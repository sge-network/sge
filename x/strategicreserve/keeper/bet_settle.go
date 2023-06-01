package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// RefundBettor process bets in case market gets cancelled or aborted,
// this method transfers th bet amount from order book liquidity module account balance to the bettor account balance.
func (k Keeper) RefundBettor(
	ctx sdk.Context,
	bettorAddress sdk.AccAddress,
	betAmount, betFee, payout sdk.Int,
	uniqueLock string,
) (err error) {
	// if no lock exist means that there is nothing to be processed.
	if !k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	// transfer bet amount from `book_liquidity_pool` to bettor's account.
	err = k.transferFundsFromModuleToAccount(ctx, types.OrderBookLiquidityPool, bettorAddress, betAmount)
	if err != nil {
		return
	}

	// transfer bet fee from `bet_fee_collector` to bettor's account.
	err = k.transferFundsFromModuleToAccount(ctx, bettypes.BetFeeCollector, bettorAddress, betFee)
	if err != nil {
		return
	}

	// delete the lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)

	return nil
}

// BettorWins process bets in case bettor is the winner,
// transfers the bet amount and the payout profit to the bettor's account and,
// updates actual profit of the participation to the subtracted value from the payout profit.
func (k Keeper) BettorWins(
	ctx sdk.Context,
	bettorAddress sdk.AccAddress,
	betAmount sdk.Int,
	payoutProfit sdk.Int,
	uniqueLock string,
	betFulfillments []*bettypes.BetFulfillment,
	orderBookUID string,
) (err error) {
	// if no lock exist means that there is nothing to be processed.
	if !k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	for _, betFulfillment := range betFulfillments {
		orderBookParticipation, found := k.GetOrderBookParticipation(ctx, orderBookUID, betFulfillment.ParticipationIndex)
		if !found {
			return sdkerrors.Wrapf(types.ErrOrderBookParticipationNotFound, "%s, %d", orderBookUID, betFulfillment.ParticipationIndex)
		}

		total := betFulfillment.PayoutProfit.Add(betFulfillment.BetAmount)
		// transfer bet amount and expected payout from the `bet_collector` account to bettor
		err = k.transferFundsFromModuleToAccount(ctx, types.OrderBookLiquidityPool, bettorAddress, total)
		if err != nil {
			return
		}

		// update actual profit of the participation, the bettor is the winner, so we need to
		// payout from the participant profit.
		orderBookParticipation.ActualProfit = orderBookParticipation.ActualProfit.Sub(betFulfillment.PayoutProfit)
		k.SetOrderBookParticipation(ctx, orderBookParticipation)
	}

	// Delete lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)

	return nil
}

// BettorLoses process bets in case bettor loses,
// adds the bet amount to the actual profit of the participation
// for each of the bet fulfillemnt records and,
// removes the payout lock.
func (k Keeper) BettorLoses(ctx sdk.Context, address sdk.AccAddress,
	betAmount sdk.Int,
	payoutProfit sdk.Int,
	uniqueLock string,
	betFulfillments []*bettypes.BetFulfillment,
	orderBookUID string,
) error {
	// if no lock exist means that there is nothing to be processed.
	if !k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	for _, betFulfillment := range betFulfillments {
		// update amount to be transferred to house
		orderBookParticipation, found := k.GetOrderBookParticipation(ctx, orderBookUID, betFulfillment.ParticipationIndex)
		if !found {
			return sdkerrors.Wrapf(types.ErrOrderBookParticipationNotFound, "%s, %d", orderBookUID, betFulfillment.ParticipationIndex)
		}

		// update actual profit of the participation, the bettor is the loser, so we need to
		// add the lost bet amount to the participant profit.
		orderBookParticipation.ActualProfit = orderBookParticipation.ActualProfit.Add(betFulfillment.BetAmount)
		k.SetOrderBookParticipation(ctx, orderBookParticipation)
	}

	// Delete lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)

	return nil
}
