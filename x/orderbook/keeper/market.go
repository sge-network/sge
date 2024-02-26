package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bettypes "github.com/sge-network/sge/x/bet/types"
)

func (k Keeper) WithdrawBetFee(ctx sdk.Context, marketCreator sdk.AccAddress, betFee, betPriceLockFee sdkmath.Int) error {
	// refund market creator's account from bet fee collector.
	if err := k.refund(bettypes.BetFeeCollectorFunder{}, ctx, marketCreator, betFee); err != nil {
		return err
	}
	return k.refund(bettypes.PriceLockFunder{}, ctx, marketCreator, betPriceLockFee)
}
