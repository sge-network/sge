package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PubKeys returns list of the public keys
func (k Keeper) PubKeys(goCtx context.Context, req *types.QueryPubKeysRequest) (*types.QueryPubKeysResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, types.ErrTextInvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Process the query

	res, found := k.GetKeyVault(ctx)

	if !found {
		return nil, types.ErrNoPublicKeysFound
	}

	return &types.QueryPubKeysResponse{List: res.PublicKeys}, nil
}
