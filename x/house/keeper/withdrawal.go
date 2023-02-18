package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/house/types"
)

// SetWithdrawal sets a withdrawal.
func (k Keeper) SetWithdrawal(ctx sdk.Context, withdrawal types.Withdrawal) error {
	depAddress := sdk.MustAccAddressFromBech32(withdrawal.DepositorAddress)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.WithdrawalKeyPrefix)
	b := types.MustMarshalWithdrawal(k.cdc, withdrawal)

	withdrawalKey := types.GetWithdrawalKey(depAddress, withdrawal.SportEventUID, withdrawal.ParticipantID, withdrawal.WithdrawalNumber)
	store.Set(withdrawalKey, b)
	return nil
}

// Withdraw performs a withdrawal of coins of unused amount corresponding to a deposit.
func (k Keeper) Withdraw(ctx sdk.Context, depAddr sdk.AccAddress, sportEventUID string, pID uint64, mode types.WithdrawalMode, witAmt sdk.Int) (uint64, error) {
	var withdrawalNumber uint64
	// Get the deposit object
	depoistKey := types.GetDepositKey(depAddr, sportEventUID, pID)
	deposit, found := k.GetDeposit(ctx, depoistKey)
	if !found {
		return withdrawalNumber, sdkerrors.Wrapf(types.ErrDepositNotFound, ": %s, %d", sportEventUID, pID)
	}

	if mode == types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL {
		if deposit.Liquidity.Sub(deposit.TotalWithdrawalAmount).LT(witAmt) {
			return withdrawalNumber, sdkerrors.Wrapf(types.ErrWithdrawalTooLarge, "%d", witAmt.Int64())
		}
	}

	// Create the withdrawal object
	withdrawal := types.NewWithdrawal(depAddr, sportEventUID, pID, deposit.Withdrawals+1, witAmt, mode)

	withdrawalAmt, err := k.orderBookKeeper.LiquidateBookParticipant(ctx, depAddr.String(), sportEventUID, pID, mode, witAmt)
	if err != nil {
		return pID, sdkerrors.Wrapf(types.ErrOrderBookLiquidateProcessing, "%s", err)
	}

	withdrawal.Amount = withdrawalAmt
	if err = k.SetWithdrawal(ctx, withdrawal); err != nil {
		return withdrawalNumber, err
	}

	deposit.Withdrawals++
	deposit.TotalWithdrawalAmount = deposit.TotalWithdrawalAmount.Add(withdrawalAmt)
	if err = k.SetDeposit(ctx, deposit); err != nil {
		return withdrawalNumber, err
	}

	return withdrawalNumber, nil
}
