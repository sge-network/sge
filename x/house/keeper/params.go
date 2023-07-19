package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/house/types"
)

// GetParams return parameters of the module
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params for the module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetMinAllowedDepositAmount returns minimum acceptable deposit amount.
func (k Keeper) GetMinAllowedDepositAmount(ctx sdk.Context) (res sdk.Int) {
	return k.GetParams(ctx).MinDeposit
}

// GetHouseParticipationFee returns % of deposit to be paid for house participation by an account
func (k Keeper) GetHouseParticipationFee(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).HouseParticipationFee
}
