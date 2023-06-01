package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

// singlePageNum used to return single page result in pagination.
const singlePageNum = 1

// SettleBet settles a single bet and updates it in KVStore
func (k Keeper) SettleBet(ctx sdk.Context, bettorAddressStr, betUID string) error {
	if !types.IsValidUID(betUID) {
		return types.ErrInvalidBetUID
	}

	uid2ID, found := k.GetBetID(ctx, betUID)
	if !found {
		return types.ErrNoMatchingBet
	}

	bet, found := k.GetBet(ctx, bettorAddressStr, uid2ID.ID)
	if !found {
		return types.ErrNoMatchingBet
	}

	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "%s", err)
	}

	if bet.Creator != bettorAddressStr {
		return types.ErrBettorAddressNotEqualToCreator
	}

	if err := checkBetStatus(bet.Status); err != nil {
		// bet cancellation logic will reside here if this feature is requested
		return err
	}

	// get the respective market for the bet
	market, found := k.marketKeeper.GetMarket(ctx, bet.MarketUID)
	if !found {
		return types.ErrNoMatchingMarket
	}

	if market.Status == markettypes.MarketStatus_MARKET_STATUS_ABORTED ||
		market.Status == markettypes.MarketStatus_MARKET_STATUS_CANCELED {
		payoutProfit, err := types.CalculatePayoutProfit(bet.OddsType, bet.OddsValue, bet.Amount)
		if err != nil {
			return err
		}

		if err := k.srKeeper.RefundBettor(ctx, bettorAddress, bet.Amount, bet.BetFee, payoutProfit.TruncateInt(), bet.UID); err != nil {
			return sdkerrors.Wrapf(types.ErrInSRRefund, "%s", err)
		}

		bet.Status = types.Bet_STATUS_SETTLED
		bet.Result = types.Bet_RESULT_REFUNDED

		k.updateSettlementState(ctx, bet, uid2ID.ID)

		return nil
	}

	// check if the bet odds is a winner odds or not and set the bet pointer states
	if err := processBetResultAndStatus(&bet, market); err != nil {
		return err
	}

	if err := k.settleResolvedBet(ctx, &bet); err != nil {
		return err
	}

	if err := k.srKeeper.WithdrawBetFee(ctx, sdk.MustAccAddressFromBech32(market.Creator), bet.BetFee); err != nil {
		return err
	}

	k.updateSettlementState(ctx, bet, uid2ID.ID)

	return nil
}

// updateSettlementState settles bet in the store
func (k Keeper) updateSettlementState(ctx sdk.Context, bet types.Bet, betID uint64) {
	// set current height as settlement height
	bet.SettlementHeight = ctx.BlockHeight()

	// store bet in the module state
	k.SetBet(ctx, bet, betID)

	// remove pending bet
	k.RemovePendingBet(ctx, bet.MarketUID, betID)

	// store settled bet in the module state
	k.SetSettledBet(ctx, types.NewSettledBet(bet.UID, bet.Creator), betID, ctx.BlockHeight())
}

// settleResolvedBet settles a bet by calling strategicReserve functions to unlock fund and payout
// based on bet's result, and updates status of bet to settled
func (k Keeper) settleResolvedBet(ctx sdk.Context, bet *types.Bet) error {
	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "%s", err)
	}

	payout, err := types.CalculatePayoutProfit(bet.OddsType, bet.OddsValue, bet.Amount)
	if err != nil {
		return err
	}

	if bet.Result == types.Bet_RESULT_LOST {
		if err := k.srKeeper.BettorLoses(ctx, bettorAddress, bet.Amount, payout.TruncateInt(), bet.UID, bet.BetFulfillment, bet.MarketUID); err != nil {
			return sdkerrors.Wrapf(types.ErrInSRBettorLoses, "%s", err)
		}
		bet.Status = types.Bet_STATUS_SETTLED
	} else if bet.Result == types.Bet_RESULT_WON {
		if err := k.srKeeper.BettorWins(ctx, bettorAddress, bet.Amount, payout.TruncateInt(), bet.UID, bet.BetFulfillment, bet.MarketUID); err != nil {
			return sdkerrors.Wrapf(types.ErrInSRBettorWins, "%s", err)
		}
		bet.Status = types.Bet_STATUS_SETTLED
	}
	return nil
}

// BatchMarketSettlements settles bets of resolved markets
// in batch. The markets get into account in FIFO manner.
func (k Keeper) BatchMarketSettlements(ctx sdk.Context) error {
	toFetch := k.GetParams(ctx).BatchSettlementCount

	// continue looping until reach batch settlement count parameter
	for toFetch > 0 {
		// get the first resolved market to process corresponding pending bets.
		marketUID, found := k.marketKeeper.GetFirstUnsettledResolvedMarket(ctx)
		// exit loop if there is no resolved bet.
		if !found {
			return nil
		}

		// settle market pending bets.
		settledCount, err := k.batchSettlementOfMarket(ctx, marketUID, toFetch)
		if err != nil {
			return fmt.Errorf("could not settle market %s %s", marketUID, err)
		}

		// check if still there is any pending bet for the market.
		pendingBetExists, err := k.IsAnyPendingBetForMarket(ctx, marketUID)
		if err != nil {
			return fmt.Errorf("could not check the pending bets %s %s", marketUID, err)
		}

		// if there is not any pending bet for the market
		// we need to remove its uid from the list of unsettled resolved bets.
		if !pendingBetExists {
			k.marketKeeper.RemoveUnsettledResolvedMarket(ctx, marketUID)
			err = k.srKeeper.SetOrderBookAsSettled(ctx, marketUID)
			if err != nil {
				return fmt.Errorf("could not resolve strategicreserve %s %s", marketUID, err)
			}
		}

		// update counter of bets to be processed in the next iteration.
		toFetch -= settledCount
	}

	return nil
}

// batchSettlementOfMarket settles pending bets of a markets
func (k Keeper) batchSettlementOfMarket(ctx sdk.Context, marketUID string, countToBeSettled uint32) (settledCount uint32, err error) {
	// initialize iterator for the certain number of pending bets
	// equal to countToBeSettled
	iterator := sdk.KVStorePrefixIteratorPaginated(
		ctx.KVStore(k.storeKey),
		types.PendingBetListOfMarketPrefix(marketUID),
		singlePageNum,
		uint(countToBeSettled))
	defer func() {
		iterErr := iterator.Close()
		if iterErr != nil {
			err = iterErr
		}
	}()

	// settle bets for the filtered pending bets
	for ; iterator.Valid(); iterator.Next() {
		var val types.PendingBet
		k.cdc.MustUnmarshal(iterator.Value(), &val)

		err = k.SettleBet(ctx, val.Creator, val.UID)
		if err != nil {
			return
		}

		// update total settled count
		settledCount++
	}

	return
}

// checkBetStatus checks status of bet. It returns an error if
// bet is canceled or settled already
func checkBetStatus(betstatus types.Bet_Status) error {
	switch betstatus {
	case types.Bet_STATUS_SETTLED:
		return types.ErrBetIsSettled
	case types.Bet_STATUS_CANCELED:
		return types.ErrBetIsCanceled
	}

	return nil
}

// processBetResultAndStatus determines the result and status of the given bet, it can be lost or won.
func processBetResultAndStatus(bet *types.Bet, market markettypes.Market) error {
	// check if market result is declared or not
	if market.Status != markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED {
		return types.ErrResultNotDeclared
	}

	var exist bool
	for _, wid := range market.WinnerOddsUIDs {
		if wid == bet.OddsUID {
			exist = true
			break
		}
	}

	if exist {
		// bettor is winner
		bet.Result = types.Bet_RESULT_WON
		bet.Status = types.Bet_STATUS_RESULT_DECLARED
		return nil
	}

	// bettor is loser
	bet.Result = types.Bet_RESULT_LOST
	bet.Status = types.Bet_STATUS_RESULT_DECLARED
	return nil
}
