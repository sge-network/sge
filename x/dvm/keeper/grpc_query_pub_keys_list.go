package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListPubKeys returns list of the public keys
func (k Keeper) ListPubKeys(goCtx context.Context, req *types.QueryListPubKeyAllRequest) (*types.QueryListPubKeyAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, types.ErrTextInvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Process the query

	res, found := k.GetPublicKeysAll(ctx)

	if !found {
		return nil, types.ErrNoPublicKeysFound
	}

	return &types.QueryListPubKeyAllResponse{List: res.List}, nil
}
