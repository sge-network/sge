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

	withdrawalKey := types.GetWithdrawalKey(depAddress, withdrawal.SportEventUid, withdrawal.ParticipantId, withdrawal.WithdrawalNumber)
	store.Set(withdrawalKey, b)
	return nil
}

// Withdraw performs a withdrawal of coins of unused amount corresponding to a deposit.
func (k Keeper) Withdraw(ctx sdk.Context, depAddr sdk.AccAddress, sportEventUid string, pId uint64, mode types.WithdrawalMode, witAmt sdk.Int) (error, uint64) {
	var withdrawalNumber uint64
	// Get the deposit object
	depoistKey := types.GetDepositKey(depAddr, sportEventUid, pId)
	deposit, found := k.GetDeposit(ctx, depoistKey)
	if !found {
		return sdkerrors.Wrapf(types.ErrDepositNotFound, ": %s, %d", sportEventUid, pId), withdrawalNumber
	}

	if mode == types.WithdrawalMode_MODE_PARTIAL {
		if deposit.Liquidity.Sub(deposit.TotalWithdrawalAmount).LT(witAmt) {
			return sdkerrors.Wrapf(types.ErrWithdrawalTooLarge, "%d", witAmt.Int64()), withdrawalNumber
		}
	}

	// Create the withdrawal object
	withdrawal := types.NewWithdrawal(depAddr, sportEventUid, pId, deposit.Withdrawals+1, witAmt, mode)

	withdrawalAmt, err := k.orderBookKeeper.LiquidateBookParticipant(ctx, depAddr.String(), sportEventUid, pId, mode, witAmt)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrOrderBookLiquidateProcessing, "%s", err), pId
	}

	withdrawal.Amount = withdrawalAmt
	if err = k.SetWithdrawal(ctx, withdrawal); err != nil {
		return err, withdrawalNumber
	}

	deposit.Withdrawals += 1
	deposit.TotalWithdrawalAmount = deposit.TotalWithdrawalAmount.Add(withdrawalAmt)
	if err = k.SetDeposit(ctx, deposit); err != nil {
		return err, withdrawalNumber
	}

	return nil, withdrawalNumber
}
