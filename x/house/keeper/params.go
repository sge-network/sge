package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/house/types"
)

// MinDeposit - Minum amount of deposit acceptable.
func (k Keeper) MinDeposit(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyMinDeposit, &res)
	return
}

// HouseParticipationFee - % of deposit to be paid for house participation by the user
func (k Keeper) HouseParticipationFee(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyHouseParticipationFee, &res)
	return
}

// Get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.MinDeposit(ctx),
		k.HouseParticipationFee(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
