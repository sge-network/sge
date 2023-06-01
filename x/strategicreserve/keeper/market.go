package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
)

func (k Keeper) WithdrawBetFee(ctx sdk.Context, marketCreator sdk.AccAddress, betFee sdk.Int) error {
	if err := k.transferFundsFromModuleToAccount(ctx, bettypes.BetFeeCollector, marketCreator, betFee); err != nil {
		return err
	}
	return nil
}
