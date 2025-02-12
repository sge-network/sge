package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bettypes "github.com/sge-network/sge/x/legacy/bet/types"
)

func (k Keeper) WithdrawBetFee(ctx sdk.Context, marketCreator sdk.AccAddress, betFee sdkmath.Int) error {
	// refund market creator's account from bet fee collector.
	return k.refund(bettypes.BetFeeCollectorFunder{}, ctx, marketCreator, betFee)
}
