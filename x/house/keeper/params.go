package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/house/types"
)

// Get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.GetMinDepositAllowedAmount(ctx),
		k.GetHouseParticipationFee(ctx),
	)
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetMinDepositAllowedAmount - Minum amount of deposit acceptable.
func (k Keeper) GetMinDepositAllowedAmount(ctx sdk.Context) (res sdk.Int) {
	k.paramstore.Get(ctx, types.KeyMinDeposit, &res)
	return
}

// GetHouseParticipationFee - % of deposit to be paid for house participation by the user
func (k Keeper) GetHouseParticipationFee(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyHouseParticipationFee, &res)
	return
}
