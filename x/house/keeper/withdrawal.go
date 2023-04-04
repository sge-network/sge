package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/house/types"
)

// SetWithdrawal sets a withdrawal.
func (k Keeper) SetWithdrawal(ctx sdk.Context, withdrawal types.Withdrawal) {
	withdrawalKey := types.GetWithdrawalKey(withdrawal.Address, withdrawal.MarketUID, withdrawal.ParticipationIndex, withdrawal.ID)

	store := k.getWithdrawalStore(ctx)
	b := k.cdc.MustMarshal(&withdrawal)
	store.Set(withdrawalKey, b)
}

// GetAllWithdrawals returns all withdrawals used during genesis dump.
func (k Keeper) GetAllWithdrawals(ctx sdk.Context) (list []types.Withdrawal, err error) {
	store := k.getWithdrawalStore(ctx)
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
func (k Keeper) Withdraw(ctx sdk.Context, creator string, marketUID string, participationIndex uint64, mode types.WithdrawalMode, witAmt sdk.Int) (uint64, error) {
	// Get the deposit object
	deposit, found := k.GetDeposit(ctx, creator, marketUID, participationIndex)
	if !found {
		return 0, sdkerrors.Wrapf(types.ErrDepositNotFound, ": %s, %d", marketUID, participationIndex)
	}

	if deposit.Creator != creator {
		return 0, sdkerrors.Wrapf(types.ErrWrongWithdrawCreator, ": %s", creator)
	}

	if mode == types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL {
		if deposit.Liquidity.Sub(deposit.TotalWithdrawalAmount).LT(witAmt) {
			return 0, sdkerrors.Wrapf(types.ErrWithdrawalTooLarge, "%d", witAmt.Int64())
		}
	}

	withdrawalID := deposit.WithdrawalCount + 1

	// Create the withdrawal object
	withdrawal := types.NewWithdrawal(withdrawalID, creator, marketUID, participationIndex, witAmt, mode)

	withdrawalAmt, err := k.srKeeper.WithdrawOrderBookParticipation(ctx, creator, marketUID, participationIndex, mode, witAmt)
	if err != nil {
		return participationIndex, sdkerrors.Wrapf(types.ErrSRLiquidateProcessing, "%s", err)
	}

	withdrawal.Amount = withdrawalAmt
	k.SetWithdrawal(ctx, withdrawal)

	deposit.WithdrawalCount++
	deposit.TotalWithdrawalAmount = deposit.TotalWithdrawalAmount.Add(withdrawalAmt)
	k.SetDeposit(ctx, deposit)

	return withdrawalID, nil
}
