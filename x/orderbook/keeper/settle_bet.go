package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/orderbook/types"
	srtypes "github.com/sge-network/sge/x/strategicreserve/types"
)

// RefundBettor process bets in case sports event gets cancelled or aborted.
func (k Keeper) RefundBettor(ctx sdk.Context, bettorAddress sdk.AccAddress, betAmount, payout sdk.Int, uniqueLock string) error {
	// Idempotency check: If lock does not exist, return error
	if !k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	// Transfer bet amount from `bet_reserve` to bettor's account
	err := k.transferFundsFromModuleToUser(ctx, srtypes.BetReserveName, bettorAddress, betAmount)
	if err != nil {
		return err
	}

	// Delete the lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)
	return nil
}

// BettorLoses process bets in case bettor looses
func (k Keeper) BettorWins(
	ctx sdk.Context,
	bettorAddress sdk.AccAddress,
	betAmount sdk.Int,
	payoutProfit sdk.Int,
	uniqueLock string,
	betFullfillments []*bettypes.BetFullfillment,
	bookId string,
) error {
	// Idempotency check: If lock does not exist, return error
	if !k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	for _, betFullfillment := range betFullfillments {
		bookParticipant, found := k.GetBookParticipant(ctx, bookId, betFullfillment.ParticipantNumber)
		if !found {
			return sdkerrors.Wrapf(types.ErrBookParticipantNotFound, "%s, %d", bookId, betFullfillment.ParticipantNumber)
		}

		// Transfer payout from the `book_liquidity_pool` account to bettor
		err := k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, bettorAddress, betFullfillment.PayoutAmount)
		if err != nil {
			return err
		}

		// Transfer bet amount from the `book_liquidity_pool` account to bettor
		err = k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, bettorAddress, betFullfillment.BetAmount)
		if err != nil {
			return err
		}

		bookParticipant.ActualProfit = bookParticipant.ActualProfit.Sub(betFullfillment.PayoutAmount)
		k.SetBookParticipant(ctx, bookParticipant)
	}

	// Delete lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)

	return nil
}

// BettorLoses process bets in case bettor looses
func (k Keeper) BettorLoses(ctx sdk.Context, address sdk.AccAddress,
	betAmount sdk.Int, payoutProfit sdk.Int, uniqueLock string, betFullfillments []*bettypes.BetFullfillment, bookId string,
) error {
	// Idempotency check: If lock does not exist, return error
	if !k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	for _, betFullfillment := range betFullfillments {
		// Update amount to be transferred to house
		bookParticipant, found := k.GetBookParticipant(ctx, bookId, betFullfillment.ParticipantNumber)
		if !found {
			return sdkerrors.Wrapf(types.ErrBookParticipantNotFound, "%s, %d", bookId, betFullfillment.ParticipantNumber)
		}
		bookParticipant.ActualProfit = bookParticipant.ActualProfit.Add(betFullfillment.BetAmount)
		k.SetBookParticipant(ctx, bookParticipant)
	}

	// Delete lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)
	return nil
}
