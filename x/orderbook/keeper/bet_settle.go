package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/orderbook/types"
)

// RefundBettor process bets in case market gets cancelled or aborted,
// this method transfers the bet amount from order book liquidity module account balance to the bettor account balance.
func (k Keeper) RefundBettor(
	ctx sdk.Context,
	bet bettypes.Bet,
) error {
	bettorAddress := sdk.MustAccAddressFromBech32(bet.Creator)

	// refund bettor's account from orderbook liquidity pool.
	if err := k.refund(types.OrderBookLiquidityFunder{}, ctx, bettorAddress, bet.Amount); err != nil {
		return err
	}

	// refund bettor's account from bet fee collector.
	if err := k.refund(bettypes.BetFeeCollectorFunder{}, ctx, bettorAddress, bet.Fee); err != nil {
		return err
	}

	return nil
}

// BettorWins process bets in case bettor is the winner,
// transfers the bet amount and the payout profit to the bettor's account and,
// updates actual profit of the participation to the subtracted value from the payout profit.
func (k Keeper) BettorWins(
	ctx sdk.Context,
	bet bettypes.Bet,
	orderBookUID string,
) error {
	for _, betFulfillment := range bet.BetFulfillment {
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
		if err := k.refund(types.OrderBookLiquidityFunder{}, ctx, sdk.MustAccAddressFromBech32(bet.Creator), betAmountAndPayout); err != nil {
			return err
		}

		// update actual profit of the participation, the bettor is the winner, so we need to
		// payout from the participant profit.
		orderBookParticipation.ActualProfit = orderBookParticipation.ActualProfit.Sub(
			betFulfillment.PayoutProfit,
		)
		k.SetOrderBookParticipation(ctx, orderBookParticipation)
	}

	return nil
}

// BettorLoses process bets in case bettor loses,
// adds the bet amount to the actual profit of the participation
// for each of the bet fulfillment records and,
// removes the payout lock.
func (k Keeper) BettorLoses(
	ctx sdk.Context,
	bet bettypes.Bet,
	orderBookUID string,
) error {
	for _, betFulfillment := range bet.BetFulfillment {
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

	return nil
}
