package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sge-network/sge/x/dvm/types"
)

type mutationModifications struct {
	Additions []string
	Deletions []string
}

// Mutation is the main transaction of DVM to add or delete the keys to the chain.
func (k msgServer) Mutation(goCtx context.Context, msg *types.MsgMutation) (*types.MsgMutationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	keys, found := k.GetPublicKeys(ctx)

	if !found {
		return nil, types.ErrNoPublicKeysFound
	}

	ticket, err := types.NewTicket(msg.Txs)
	if err != nil {
		return nil, err
	}
	// check the expiration of ticket
	err = ticket.IsValid(ctx)
	if err != nil {
		return nil, err
	}

	err = ticket.Consensus(keys.List...)
	if err != nil {
		return nil, err
	}

	var request = mutationModifications{}
	err = ticket.Unmarshal(&request)
	if err != nil {
		return nil, err
	}

	err = mutateList(&keys, request)
	if err != nil {
		return nil, err
	}

	k.Keeper.SetPublicKeys(ctx, keys)

	return &types.MsgMutationResponse{Success: true}, nil
}

// mutateList verifies and adds the Addition Keys and remove the deletion keys from the list of public Keys.
func mutateList(ks *types.PublicKeys, modifs mutationModifications) error {

	// populate map of keys
	mKeys := make(map[string]string)
	for _, v := range ks.List {
		mKeys[v] = ""
	}

	for _, v := range modifs.Additions {

		// check if pem content is a valid ED25516 key
		P, err := jwt.ParseEdPublicKeyFromPEM([]byte(v))
		if err != nil {
			return err
		}
		_ = P
		//add the key to the list
		mKeys[v] = ""
	}

	for _, v := range modifs.Deletions {
		delete(mKeys, v)
	}

	var res = []string{}
	for key := range mKeys {
		res = append(res, key)
	}
	ks.List = res

	return nil
}
