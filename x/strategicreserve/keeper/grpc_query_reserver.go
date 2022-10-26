package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// Reserver returns current reserver.
func (k Keeper) Reserver(c context.Context, _ *types.QueryReserverRequest) (*types.QueryReserverResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	reserver := k.GetReserver(ctx)

	return &types.QueryReserverResponse{Reserver: &reserver}, nil
}
