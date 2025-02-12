package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// GetParams Get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetMaxOrderBookParticipationsAllowed - Max number of orderbook participations allowed.
func (k Keeper) GetMaxOrderBookParticipationsAllowed(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).MaxOrderBookParticipations
}

// GetBatchSettlementCount - number of deposits to be settled in end blocker.
func (k Keeper) GetBatchSettlementCount(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).BatchSettlementCount
}

// GetRequeueThreshold - threshold below which a participation is requeued in fullfillment queue.
func (k Keeper) GetRequeueThreshold(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).RequeueThreshold
}
