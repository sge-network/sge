package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/sportevent/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetDefaultBetConstraints get bet constraint values of the bet constraints
func (k Keeper) GetDefaultBetConstraints(ctx sdk.Context) (params *types.EventBetConstraints) {
	p := k.GetParams(ctx)
	return types.NewEventBetConstraints(p.EventMinBetAmount, p.EventMinBetFee)
}
