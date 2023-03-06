package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/bet/types"
	sporteventtypes "github.com/sge-network/sge/x/sportevent/types"
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidCreator, err)
	}

	if bet.Creator != bettorAddressStr {
		return types.ErrBettorAddressNotEqualToCreator
	}

	if err := checkBetStatus(bet.Status); err != nil {
		// bet cancelation logic will reside here if this feature is requested
		return err
	}

	// get the respective sport-event for the bet
	sportEvent, found := k.sportEventKeeper.GetSportEvent(ctx, bet.SportEventUID)
	if !found {
		return types.ErrNoMatchingSportEvent
	}

	if sportEvent.Status == sporteventtypes.SportEventStatus_SPORT_EVENT_STATUS_ABORTED ||
		sportEvent.Status == sporteventtypes.SportEventStatus_SPORT_EVENT_STATUS_CANCELED {
		bet.Result = types.Bet_RESULT_ABORTED

		payoutProfit, err := types.CalculatePayoutProfit(bet.OddsType, bet.OddsValue, bet.Amount)
		if err != nil {
			return err
		}

		if err := k.orderbookKeeper.RefundBettor(ctx, bettorAddress, bet.Amount, payoutProfit.TruncateInt(), bet.UID); err != nil {
			return sdkerrors.Wrapf(types.ErrInSRRefund, "%s", err)
		}

		bet.Status = types.Bet_STATUS_SETTLED

		k.updateSettlementState(ctx, bet, uid2ID.ID)

		return nil
	}

	// check if the bet odds is a winner odds or not and set the bet pointer states
	if err := processBetResultAndStatus(&bet, sportEvent); err != nil {
		return err
	}

	if err := k.settleResolvedBet(ctx, &bet); err != nil {
		return err
	}

	k.updateSettlementState(ctx, bet, uid2ID.ID)

	return nil
}

// updateSettlementState settles bet in the store
func (k Keeper) updateSettlementState(ctx sdk.Context, bet types.Bet, betID uint64) {
	// set current height as settlement heigth
	bet.SettlementHeight = ctx.BlockHeight()

	// store bet in the module state
	k.SetBet(ctx, bet, betID)

	// remove active bet
	k.RemoveActiveBet(ctx, bet.SportEventUID, betID)

	// store settled bet in the module state
	k.SetSettledBet(ctx, types.NewSettledBet(bet.UID, bet.Creator), betID, ctx.BlockHeight())
}

// settleResolvedBet settles a bet by calling strategicReserve functions to unlock fund and payout
// based on bet's result, and updates status of bet to settled
func (k Keeper) settleResolvedBet(ctx sdk.Context, bet *types.Bet) error {
	bettorAddress, err := sdk.AccAddressFromBech32(bet.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidCreator, err)
	}

	payout, err := types.CalculatePayoutProfit(bet.OddsType, bet.OddsValue, bet.Amount)
	if err != nil {
		return err
	}

	if bet.Result == types.Bet_RESULT_LOST {
		if err := k.orderbookKeeper.BettorLoses(ctx, bettorAddress, bet.Amount, payout.TruncateInt(), bet.UID, bet.BetFulfillment, bet.SportEventUID); err != nil {
			return sdkerrors.Wrapf(types.ErrInSRBettorLoses, "%s", err)
		}
		bet.Status = types.Bet_STATUS_SETTLED
	} else if bet.Result == types.Bet_RESULT_WON {
		if err := k.orderbookKeeper.BettorWins(ctx, bettorAddress, bet.Amount, payout.TruncateInt(), bet.UID, bet.BetFulfillment, bet.SportEventUID); err != nil {
			return sdkerrors.Wrapf(types.ErrInSRBettorWins, "%s", err)
		}
		bet.Status = types.Bet_STATUS_SETTLED
	}
	return nil
}

// BatchSportEventSettlements settles bets of resolved sport-events
// in batch. The sport-events get into account according in FIFO.
func (k Keeper) BatchSportEventSettlements(ctx sdk.Context) error {
	toFetch := k.GetParams(ctx).BatchSettlementCount

	// continue looping until reach batch settlement count parameter
	for toFetch > 0 {
		// get the first resolved sport-event to process corresponding active bets.
		sportEventUID, found := k.sportEventKeeper.GetFirstUnsettledResovedSportEvent(ctx)
		// exit loop if there is no resolved bet.
		if !found {
			return nil
		}

		// settle sport event active bets.
		settledCount, err := k.batchSettlementOfSportEvent(ctx, sportEventUID, toFetch)
		if err != nil {
			return fmt.Errorf("could not settle sport event %s %s", sportEventUID, err)
		}

		// check if still there is any active bet for the sport-event.
		isThereAnyActiveBet, err := k.IsAnyActiveBetForSportEvent(ctx, sportEventUID)
		if err != nil {
			return fmt.Errorf("could not check the active bets %s %s", sportEventUID, err)
		}

		// if there is not any active bet for the sport-event
		// we need to remove its uid from the list of unsettled resolved bets.
		if !isThereAnyActiveBet {
			k.sportEventKeeper.RemoveUnsettledResolvedSportEvent(ctx, sportEventUID)
			err = k.orderbookKeeper.AddBookSettlement(ctx, sportEventUID)
			if err != nil {
				return fmt.Errorf("could not resolve orderbook %s %s", sportEventUID, err)
			}
		}

		// update counter of bets to be processed in the next iteration.
		toFetch -= settledCount
	}

	return nil
}

// batchSettlementOfSportEvent settles active of a sport-events
func (k Keeper) batchSettlementOfSportEvent(ctx sdk.Context, sportEventUID string, countToBeSettled uint32) (settledCount uint32, err error) {
	// initialize iterator for the certain number of active bets
	// equal to countToBeSettled
	iterator := sdk.KVStorePrefixIteratorPaginated(
		ctx.KVStore(k.storeKey),
		types.ActiveBetListOfSportEventPrefix(sportEventUID),
		singlePageNum,
		uint(countToBeSettled))
	defer func() {
		iterErr := iterator.Close()
		if iterErr != nil {
			err = iterErr
		}
	}()

	// settle bets for the filtered active bets
	for ; iterator.Valid(); iterator.Next() {
		var val types.ActiveBet
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
func processBetResultAndStatus(bet *types.Bet, sportEvent sporteventtypes.SportEvent) error {
	// check if sport-event result is declared or not
	if sportEvent.Status != sporteventtypes.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED {
		return types.ErrResultNotDeclared
	}

	var exist bool
	for _, wid := range sportEvent.WinnerOddsUIDs {
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
