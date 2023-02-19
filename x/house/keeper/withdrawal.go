package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/house/types"
)

// SetWithdrawal sets a withdrawal.
func (k Keeper) SetWithdrawal(ctx sdk.Context, withdrawal types.Withdrawal) {
	withdrawalKey := types.GetWithdrawalKey(withdrawal.Creator, withdrawal.SportEventUID, withdrawal.ParticipantID, withdrawal.ID)

	store := k.getWithdrawalsStore(ctx)
	b := k.cdc.MustMarshal(&withdrawal)
	store.Set(withdrawalKey, b)
}

// GetAllWithdrawals returns all withdrawals used during genesis dump.
func (k Keeper) GetAllWithdrawals(ctx sdk.Context) (list []types.Withdrawal, err error) {
	store := k.getWithdrawalsStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Withdrawal
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Withdraw performs a withdrawal of coins of unused amount corresponding to a deposit.
func (k Keeper) Withdraw(ctx sdk.Context, creator string, sportEventUID string, participantID uint64, mode types.WithdrawalMode, witAmt sdk.Int) (uint64, error) {
	var withdrawalNumber uint64
	// Get the deposit object
	deposit, found := k.GetDeposit(ctx, creator, sportEventUID, participantID)
	if !found {
		return withdrawalNumber, sdkerrors.Wrapf(types.ErrDepositNotFound, ": %s, %d", sportEventUID, participantID)
	}

	if mode == types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL {
		if deposit.Liquidity.Sub(deposit.TotalWithdrawalAmount).LT(witAmt) {
			return withdrawalNumber, sdkerrors.Wrapf(types.ErrWithdrawalTooLarge, "%d", witAmt.Int64())
		}
	}

	withdrawalID := deposit.WithdrawalCount + 1

	// Create the withdrawal object
	withdrawal := types.NewWithdrawal(withdrawalID, creator, sportEventUID, participantID, witAmt, mode)

	withdrawalAmt, err := k.orderBookKeeper.LiquidateBookParticipant(ctx, creator, sportEventUID, participantID, mode, witAmt)
	if err != nil {
		return participantID, sdkerrors.Wrapf(types.ErrOrderBookLiquidateProcessing, "%s", err)
	}

	withdrawal.Amount = withdrawalAmt
	k.SetWithdrawal(ctx, withdrawal)

	deposit.WithdrawalCount++
	deposit.TotalWithdrawalAmount = deposit.TotalWithdrawalAmount.Add(withdrawalAmt)
	k.SetDeposit(ctx, deposit)

	return withdrawalNumber, nil
}
