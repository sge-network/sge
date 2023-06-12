package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
)

func (k Keeper) WithdrawBetFee(ctx sdk.Context, marketCreator sdk.AccAddress, betFee sdk.Int) error {
	// refund market creator's account from bet fee collector.
	if err := k.refund(bettypes.BetFeeCollectorFunder{}, ctx, marketCreator, betFee); err != nil {
		return err
	}
	return nil
}
