package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/types"
)

// GetParams Get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.GetMaxOrderBookParticipationsAllowed(ctx),
		k.GetBatchSettlementCount(ctx),
		k.GetRequeueThreshold(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetMaxOrderBookParticipationsAllowed - Max number of orderbook participations allowed.
func (k Keeper) GetMaxOrderBookParticipationsAllowed(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxOrderBookParticipations, &res)
	return
}

// GetBatchSettlementCount - number of deposits to be settled in end blocker.
func (k Keeper) GetBatchSettlementCount(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyBatchSettlementCount, &res)
	return
}

// GetRequeueThreshold - threshold below which a participation is requeued in fullfillment queue.
func (k Keeper) GetRequeueThreshold(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyRequeueThreshold, &res)
	return
}
