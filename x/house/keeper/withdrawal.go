package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/house/types"
)

// SetWithdrawal sets a withdrawal.
func (k Keeper) SetWithdrawal(ctx sdk.Context, withdrawal types.Withdrawal) {
	withdrawalKey := types.GetWithdrawalKey(
		withdrawal.Address,
		withdrawal.MarketUID,
		withdrawal.ParticipationIndex,
		withdrawal.ID,
	)

	store := k.getWithdrawalStore(ctx)
	b := k.cdc.MustMarshal(&withdrawal)
	store.Set(withdrawalKey, b)
}

// GetWithdraw returns a specific withdrawal from the store.
func (k Keeper) GetWithdraw(ctx sdk.Context, depositorAddress,
	marketUID string, participationIndex, id uint64,
) (val types.Withdrawal, found bool) {
	marketsStore := k.getWithdrawalStore(ctx)
	withdrawKey := types.GetWithdrawalKey(depositorAddress, marketUID, participationIndex, id)
	b := marketsStore.Get(withdrawKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
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
func (k Keeper) Withdraw(
	ctx sdk.Context,
	deposit types.Deposit,
	creator, depositorAddr string,
	marketUID string,
	participationIndex uint64,
	mode types.WithdrawalMode,
	withdrawableAmount sdk.Int,
) (uint64, error) {
	// set next id
	withdrawalID := deposit.WithdrawalCount + 1

	err := k.orderbookKeeper.WithdrawOrderBookParticipation(
		ctx,
		marketUID,
		participationIndex,
		withdrawableAmount,
	)
	if err != nil {
		return 0, sdkerrors.Wrapf(types.ErrOBLiquidateProcessing, "%s", err)
	}

	// Create the withdrawal object
	withdrawal := types.NewWithdrawal(
		withdrawalID,
		creator,
		depositorAddr,
		marketUID,
		participationIndex,
		withdrawableAmount,
		mode,
	)
	k.SetWithdrawal(ctx, withdrawal)

	deposit.WithdrawalCount++
	deposit.TotalWithdrawalAmount = deposit.TotalWithdrawalAmount.Add(withdrawableAmount)
	k.SetDeposit(ctx, deposit)

	_ = k.orderbookKeeper.PublishOrderBookEvent(ctx, marketUID)

	return withdrawalID, nil
}
