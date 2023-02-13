package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sge-network/sge/x/dvm/types"
)

// SubmitPubkeysChangeProposal is the main transaction of DVM to add or delete the keys to the chain.
func (k msgServer) SubmitPubkeysChangeProposal(goCtx context.Context, msg *types.MsgSubmitPubkeysChangeProposalRequest) (*types.MsgSubmitPubkeysChangeProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	keys, found := k.GetKeyVault(ctx)

	if !found {
		return nil, types.ErrNoPublicKeysFound
	}

	payload := types.PubkeysChangeProposalPayload{}
	err := k.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload)
	if err != nil {
		return nil, err
	}

	err = payload.Validate(keys.PublicKeys)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ticket payload is not valid", err)
	}

	stats := k.GetProposalStats(ctx)
	stats.PubkeysChangeCount += 1

	// set active proposal
	k.Keeper.SetActivePubkeysChangeProposal(ctx,
		types.NewPublicKeysChangeProposal(
			stats.PubkeysChangeCount,
			msg.Creator,
			payload,
			ctx.BlockTime().Unix(),
		),
	)

	// set proposal statistics
	k.Keeper.SetProposalStats(ctx, stats)

	return &types.MsgSubmitPubkeysChangeProposalResponse{Success: true}, nil
}

// mutateList verifies and adds the Addition Keys and remove the deletion keys from the list of public Keys.
func mutateList(ks *types.KeyVault, modifs types.PubkeysChangeProposalPayload) error {
	// populate map of keys
	mKeys := make(map[string]string)
	for _, v := range ks.PublicKeys {
		mKeys[v] = ""
	}

	for _, v := range modifs.Additions {
		// check if pem content is a valid ED25516 key
		P, err := jwt.ParseEdPublicKeyFromPEM([]byte(v))
		if err != nil {
			return err
		}
		_ = P

		// add the key to the list
		mKeys[v] = ""
	}

	for _, v := range modifs.Deletions {
		delete(mKeys, v)
	}

	res := []string{}
	for key := range mKeys {
		res = append(res, key)
	}
	ks.PublicKeys = res

	return nil
}
