package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// RefundBettor process bets in case market gets cancelled or aborted.
func (k Keeper) RefundBettor(ctx sdk.Context, bettorAddress sdk.AccAddress, betAmount, payout sdk.Int, uniqueLock string) error {
	// Idempotency check: If lock does not exist, return error
	if !k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	// Transfer bet amount from `bet_reserve` to bettor's account
	err := k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, bettorAddress, betAmount)
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
	betFulfillments []*bettypes.BetFulfillment,
	bookUID string,
) error {
	// Idempotency check: If lock does not exist, return error
	if !k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	for _, betFulfillment := range betFulfillments {
		bookParticipation, found := k.GetBookParticipation(ctx, bookUID, betFulfillment.ParticipationIndex)
		if !found {
			return sdkerrors.Wrapf(types.ErrBookParticipationNotFound, "%s, %d", bookUID, betFulfillment.ParticipationIndex)
		}

		// Transfer payout from the `book_liquidity_pool` account to bettor
		err := k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, bettorAddress, betFulfillment.PayoutAmount)
		if err != nil {
			return err
		}

		// Transfer bet amount from the `book_liquidity_pool` account to bettor
		err = k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, bettorAddress, betFulfillment.BetAmount)
		if err != nil {
			return err
		}

		bookParticipation.ActualProfit = bookParticipation.ActualProfit.Sub(betFulfillment.PayoutAmount)
		k.SetBookParticipation(ctx, bookParticipation)
	}

	// Delete lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)

	return nil
}

// BettorLoses process bets in case bettor looses
func (k Keeper) BettorLoses(ctx sdk.Context, address sdk.AccAddress,
	betAmount sdk.Int, payoutProfit sdk.Int, uniqueLock string, betFulfillments []*bettypes.BetFulfillment, bookUID string,
) error {
	// Idempotency check: If lock does not exist, return error
	if !k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	for _, betFulfillment := range betFulfillments {
		// Update amount to be transferred to house
		bookParticipation, found := k.GetBookParticipation(ctx, bookUID, betFulfillment.ParticipationIndex)
		if !found {
			return sdkerrors.Wrapf(types.ErrBookParticipationNotFound, "%s, %d", bookUID, betFulfillment.ParticipationIndex)
		}
		bookParticipation.ActualProfit = bookParticipation.ActualProfit.Add(betFulfillment.BetAmount)
		k.SetBookParticipation(ctx, bookParticipation)
	}

	// Delete lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)
	return nil
}
