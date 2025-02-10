package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/bet/types"
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

// GetConstraints get bet constraint values of the bet wager constraints
func (k Keeper) GetConstraints(ctx sdk.Context) (params types.Constraints) {
	return k.GetParams(ctx).Constraints
}
