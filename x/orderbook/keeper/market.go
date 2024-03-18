package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bettypes "github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

func (k Keeper) WithdrawBetFee(ctx sdk.Context, marketCreator sdk.AccAddress, betFee sdkmath.Int) error {
	// refund market creator's account from bet fee collector.
	return k.refund(bettypes.BetFeeCollectorFunder{}, ctx, marketCreator, betFee)
}

func (k Keeper) WithdrawPriceLockFee(ctx sdk.Context, marketCreator sdk.AccAddress, betFee sdkmath.Int) error {
	// refund market creator's account from price lock fee collector.
	return k.refund(markettypes.PriceLockFeeCollector{}, ctx, marketCreator, betFee)
}
