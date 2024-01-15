package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/subaccount/types"
)

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetParams return parameters of the module
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// GetWagerEnabled returns enable/disable status of wager
func (k Keeper) GetWagerEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).WagerEnabled
}

// GetDepositEnabled returns enable/disable status of deposit
func (k Keeper) GetDepositEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).DepositEnabled
}
