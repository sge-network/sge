package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/orderbook/types"
)

// RefundBettor process bets in case market gets cancelled or aborted,
// this method transfers the bet amount from order book liquidity module account balance to the bettor account balance.
func (k Keeper) RefundBettor(
	ctx sdk.Context,
	bettorAddress sdk.AccAddress,
	betAmount, betFee, _ sdk.Int,
	_ string,
) error {
	// refund bettor's account from orderbook liquidity pool.
	if err := k.refund(types.OrderBookLiquidityFunder{}, ctx, bettorAddress, betAmount); err != nil {
		return err
	}

	// refund bettor's account from bet fee collector.
	if err := k.refund(bettypes.BetFeeCollectorFunder{}, ctx, bettorAddress, betFee); err != nil {
		return err
	}

	for _, hook := range k.hooks {
		hook.AfterBettorRefund(ctx, bettorAddress, betAmount, betFee)
	}
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
	_ string,
	betFulfillments []*bettypes.BetFulfillment,
	orderBookUID string,
) error {
	for _, betFulfillment := range betFulfillments {
		orderBookParticipation, found := k.GetOrderBookParticipation(
			ctx,
			orderBookUID,
			betFulfillment.ParticipationIndex,
		)
		if !found {
			return sdkerrors.Wrapf(
				types.ErrOrderBookParticipationNotFound,
				"%s, %d",
				orderBookUID,
				betFulfillment.ParticipationIndex,
			)
		}

		betAmountAndPayout := betFulfillment.PayoutProfit.Add(betFulfillment.BetAmount)
		// refund bettor's account from orderbook liquidity pool.
		if err := k.refund(types.OrderBookLiquidityFunder{}, ctx, bettorAddress, betAmountAndPayout); err != nil {
			return err
		}

		// update actual profit of the participation, the bettor is the winner, so we need to
		// payout from the participant profit.
		orderBookParticipation.ActualProfit = orderBookParticipation.ActualProfit.Sub(
			betFulfillment.PayoutProfit,
		)
		k.SetOrderBookParticipation(ctx, orderBookParticipation)
	}

	for _, h := range k.hooks {
		h.AfterBettorWin(ctx, bettorAddress, betAmount, payoutProfit)
	}

	return nil
}

// BettorLoses process bets in case bettor loses,
// adds the bet amount to the actual profit of the participation
// for each of the bet fulfillment records and,
// removes the payout lock.
func (k Keeper) BettorLoses(ctx sdk.Context, address sdk.AccAddress,
	betAmount sdk.Int,
	_ sdk.Int,
	_ string,
	betFulfillments []*bettypes.BetFulfillment,
	orderBookUID string,
) error {
	for _, betFulfillment := range betFulfillments {
		// update amount to be transferred to house
		orderBookParticipation, found := k.GetOrderBookParticipation(
			ctx,
			orderBookUID,
			betFulfillment.ParticipationIndex,
		)
		if !found {
			return sdkerrors.Wrapf(
				types.ErrOrderBookParticipationNotFound,
				"%s, %d",
				orderBookUID,
				betFulfillment.ParticipationIndex,
			)
		}

		// update actual profit of the participation, the bettor is the loser, so we need to
		// add the lost bet amount to the participant profit.
		orderBookParticipation.ActualProfit = orderBookParticipation.ActualProfit.Add(
			betFulfillment.BetAmount,
		)
		k.SetOrderBookParticipation(ctx, orderBookParticipation)
	}

	for _, h := range k.hooks {
		h.AfterBettorLoss(ctx, address, betAmount)
	}

	return nil
}
