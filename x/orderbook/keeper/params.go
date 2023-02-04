package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/types"
)

// MaxBookParticipants - Max number of book participants allowed.
func (k Keeper) MaxBookParticipants(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMaxBookParticipants, &res)
	return
}

// Get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.MaxBookParticipants(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
