package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
	"github.com/spf13/cast"
)

type fulfillmentInfo struct {
	betUID                  string
	betID                   uint64
	oddsUID                 string
	oddsType                bettypes.OddsType
	oddsVal                 string
	maxLossMultiplier       sdk.Dec
	betAmount               sdk.Int
	payoutProfit            sdk.Dec
	fulfiledBetAmount       sdk.Int
	totalAvailableLiquidity sdk.Int

	fulfillmentQueue []uint64
	fulfillmentMap   fulfillmentMap
	inProcessItem    fulfillmentItem
	fulfillments     []*bettypes.BetFulfillment
}

func (fInfo *fulfillmentInfo) setItemFulfilledAndRemove() {
	fInfo.inProcessItem.setFulfilled()
	fInfo.removeQueueItem()
}

func (fInfo *fulfillmentInfo) removeQueueItem() {
	fInfo.fulfillmentQueue = fInfo.fulfillmentQueue[1:]
}

func (fInfo *fulfillmentInfo) hasUnfulfilledQueueItem() bool {
	return len(fInfo.fulfillmentQueue) > 0
}

func (fInfo *fulfillmentInfo) IsFulfilled() bool {
	// if the remaining payout is less than 1.00, means that the decimal part will be ignored
	return fInfo.payoutProfit.LT(sdk.OneDec()) || len(fInfo.fulfillmentQueue) == 0
}

func (fInfo *fulfillmentInfo) NoMoreLiquidityAvailable() bool {
	// if the remaining payout is less than 1.00, means that the decimal part will be ignored
	return fInfo.payoutProfit.GTE(sdk.OneDec())
}

func (fInfo *fulfillmentInfo) notEnoughLiquidityAvailable() bool {
	return fInfo.inProcessItem.availableLiquidity.ToDec().LTE(fInfo.payoutProfit)
}

func (fInfo *fulfillmentInfo) isLiquidityLessThanThreshold(threshold sdk.Int) bool {
	diff := fInfo.inProcessItem.availableLiquidity.Sub(fInfo.payoutProfit.TruncateInt())
	return diff.LTE(threshold)
}

type fulfillmentItem struct {
	availableLiquidity    sdk.Int
	participation         types.OrderBookParticipation
	participationExposure types.ParticipationExposure
}

func (item *fulfillmentItem) noLiquidityAvailable() bool {
	return item.availableLiquidity.LTE(sdk.ZeroInt())
}

func (item *fulfillmentItem) setFulfilled() {
	item.participationExposure.IsFulfilled = true
	item.participation.ExposuresNotFilled--
}

func (item *fulfillmentItem) setAvailableLiquidity(maxLossMultiplier sdk.Dec) {
	item.availableLiquidity = item.calcAvailableLiquidity(maxLossMultiplier)
}

func (item *fulfillmentItem) calcAvailableLiquidity(maxLossMultiplier sdk.Dec) sdk.Int {
	return maxLossMultiplier.
		MulInt(item.participation.CurrentRoundLiquidity).
		Sub(sdk.NewDecFromInt(item.participationExposure.Exposure)).TruncateInt()
}

func (item *fulfillmentItem) allExposureFulfilled() bool {
	return item.participation.ExposuresNotFilled == 0
}

type fulfillmentMap map[uint64]fulfillmentItem

func (fMap fulfillmentMap) setParticipation(participation types.OrderBookParticipation) {
	fItem := fMap[participation.Index]
	fItem.participation = participation
	fMap[participation.Index] = fItem
}

func (fMap fulfillmentMap) setExposure(participationIndex uint64, exposure types.ParticipationExposure) {
	fItem := fMap[participationIndex]
	fItem.participationExposure = exposure
	fMap[participationIndex] = fItem
}

// ProcessBetPlacement processes bet placement
func (k Keeper) ProcessBetPlacement(
	ctx sdk.Context,
	betUID, bookUID, oddsUID string,
	maxLossMultiplier sdk.Dec,
	betAmount sdk.Int,
	payoutProfit sdk.Dec,
	bettorAddress sdk.AccAddress,
	betFee sdk.Int,
	oddsType bettypes.OddsType,
	oddsVal string, betID uint64,
) ([]*bettypes.BetFulfillment, error) {
	// If lock exists, return error
	// Lock already exists means the bet is already placed for the given bet-uid
	if k.payoutLockExists(ctx, betUID) {
		return nil, sdkerrors.Wrapf(types.ErrLockAlreadyExists, "%s", betUID)
	}

	// get book data by its id
	book, found := k.GetOrderBook(ctx, bookUID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", bookUID)
	}
	// get exposures of the odds bettor is placing the bet
	bookExposure, found := k.GetOrderBookOddsExposure(ctx, bookUID, oddsUID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrOrderBookExposureNotFound, "%s , %s", bookUID, oddsUID)
	}

	fInfo, err := k.initFulfillmentInfo(ctx, betAmount, payoutProfit, betUID, betID, oddsUID, oddsType, oddsVal, maxLossMultiplier, &book)
	if err != nil {
		return nil, err
	}

	err = k.fulfillBetByParticipationQueue(ctx, &fInfo, &book, bookExposure.FulfillmentQueue)
	if err != nil {
		return nil, err
	}

	bookExposure.FulfillmentQueue = fInfo.fulfillmentQueue
	k.SetOrderBookOddsExposure(ctx, bookExposure)

	// Transfer bet fee from bettor to the `bet` module account
	if err = k.transferFundsFromAccountToModule(ctx, bettorAddress, bettypes.BetFeeCollector, betFee); err != nil {
		return nil, err
	}

	// Transfer bet amount from bettor to `bet_collector` Account
	if err = k.transferFundsFromAccountToModule(ctx, bettorAddress, types.OrderBookLiquidityPool, fInfo.fulfiledBetAmount); err != nil {
		return nil, err
	}

	// Create a unique lock in the Payout Store for the bet
	k.SetPayoutLock(ctx, betUID)

	return fInfo.fulfillments, nil
}

// fulfillBetByParticipationQueue fulfills the bet placement payout using the participations
// that is stored in the state.
func (k Keeper) fulfillBetByParticipationQueue(
	ctx sdk.Context,
	fInfo *fulfillmentInfo,
	book *types.OrderBook,
	fulfillmentQueue []uint64,
) error {
	fInfo.fulfillmentQueue = fulfillmentQueue

	// the decimal amount that is being lost in the bet amount calculation from payout profit
	truncatedBetAmount := sdk.NewDec(0)

	// continue until updatedFulfillmentQueue gets empty
	for fInfo.hasUnfulfilledQueueItem() {
		var err error
		// fill participation and exposure values
		fInfo.inProcessItem = fInfo.fulfillmentMap[fInfo.fulfillmentQueue[0]]

		// availableLiquidty is the available amount of tokens to be used from the participation exposure
		fInfo.inProcessItem.setAvailableLiquidity(fInfo.maxLossMultiplier)

		switch {
		case fInfo.inProcessItem.noLiquidityAvailable():
			fInfo.setItemFulfilledAndRemove()
		case fInfo.notEnoughLiquidityAvailable():
			var betAmountToFulfill sdk.Int
			betAmountToFulfill, truncatedBetAmount, err = bettypes.CalculateBetAmountInt(fInfo.oddsType, fInfo.oddsVal, fInfo.inProcessItem.availableLiquidity.ToDec(), truncatedBetAmount)
			if err != nil {
				return err
			}
			// if the available liquidity is less than remaining payout profit that
			// need to be paid, we should use all of available liquidity pull for the calculations.
			if err = k.fulfill(ctx, fInfo, betAmountToFulfill, fInfo.inProcessItem.availableLiquidity); err != nil {
				return err
			}
			fInfo.setItemFulfilledAndRemove()
		default:
			// availableLiquidty is positive and more than remaining payout profit that
			// need to be paid, so we can cover all of payout profits with available liquidity.
			// this case appends the last fulfillment
			if fInfo.isLiquidityLessThanThreshold(sdk.NewIntFromUint64(k.GetRequeueThreshold(ctx))) {
				fInfo.setItemFulfilledAndRemove()
			}
			if err := k.fulfill(ctx, fInfo, fInfo.betAmount, fInfo.payoutProfit.TruncateInt()); err != nil {
				return err
			}
		}

		k.SetParticipationExposure(ctx, fInfo.inProcessItem.participationExposure)
		k.SetOrderBookParticipation(ctx, fInfo.inProcessItem.participation)

		// if there are no more exposures to be filled
		if fInfo.inProcessItem.allExposureFulfilled() {
			err := k.refreshQueueAndState(ctx, fInfo, book)
			if err != nil {
				return err
			}
		}

		// if the remaining payout is less than 1.00, means that the decimal part will be ignored
		if fInfo.IsFulfilled() {
			break
		}
	}

	if fInfo.NoMoreLiquidityAvailable() {
		return sdkerrors.Wrapf(types.ErrInternalProcessingBet, "insufficient liquidity in order book")
	}

	return nil
}

// initFulfillmentInfo initializes the fulfillment info for the queue iteration process.
func (k Keeper) initFulfillmentInfo(
	ctx sdk.Context,
	betAmount sdk.Int,
	payoutProfit sdk.Dec,
	betUID string,
	betID uint64,
	oddsUID string,
	oddsType bettypes.OddsType,
	oddsVal string,
	maxLossMultiplier sdk.Dec,
	book *types.OrderBook,
) (
	fInfo fulfillmentInfo,
	err error,
) {
	fInfo = fulfillmentInfo{
		// bet specs
		betAmount:               betAmount,
		payoutProfit:            payoutProfit,
		betUID:                  betUID,
		betID:                   betID,
		oddsUID:                 oddsUID,
		oddsType:                oddsType,
		oddsVal:                 oddsVal,
		maxLossMultiplier:       maxLossMultiplier,
		totalAvailableLiquidity: sdk.NewInt(0),

		//  in process maps
		fulfillmentMap: make(map[uint64]fulfillmentItem),

		// initialize the fulfilled bet amount with 0
		fulfiledBetAmount: sdk.NewInt(0),
	}

	bps, err := k.GetParticipationsOfOrderBook(ctx, book.UID)
	if err != nil {
		return
	}
	if book.ParticipationCount != cast.ToUint64(len(bps)) {
		err = sdkerrors.Wrapf(types.ErrBookParticipationsNotFound, "%s", book.UID)
		return
	}

	pes, err := k.GetExposureByOrderBookAndOdds(ctx, book.UID, oddsUID)
	if err != nil {
		return
	}
	if book.ParticipationCount != cast.ToUint64(len(pes)) {
		err = sdkerrors.Wrapf(types.ErrParticipationExposuresNotFound, "%s, %s", book.UID, oddsUID)
		return
	}

	for _, bp := range bps {
		fInfo.fulfillmentMap[bp.Index] = fulfillmentItem{
			participation: bp,
		}
	}

	for _, pe := range pes {
		item, found := fInfo.fulfillmentMap[pe.ParticipationIndex]
		if found {
			item.participationExposure = pe
			fInfo.fulfillmentMap[pe.ParticipationIndex] = item

			fInfo.totalAvailableLiquidity = fInfo.totalAvailableLiquidity.
				Add(item.calcAvailableLiquidity(maxLossMultiplier))
		}
	}

	err = fInfo.validate()
	if err != nil {
		return
	}

	return fInfo, nil
}

func (fInfo *fulfillmentInfo) validate() error {
	for _, participationIndex := range fInfo.fulfillmentQueue {
		_, found := fInfo.fulfillmentMap[participationIndex]
		if !found {
			return sdkerrors.Wrapf(types.ErrOrderBookParticipationNotFound, "%d", participationIndex)
		}

		participationExposure, found := fInfo.fulfillmentMap[participationIndex]
		if !found {
			return sdkerrors.Wrapf(types.ErrParticipationExposureNotFound, "%d", participationIndex)
		}
		if participationExposure.participationExposure.IsFulfilled {
			return sdkerrors.Wrapf(types.ErrParticipationExposureAlreadyFilled, "%d", participationIndex)
		}
	}

	if fInfo.totalAvailableLiquidity.LT(fInfo.payoutProfit.TruncateInt()) {
		return sdkerrors.Wrapf(types.ErrParticipationsCanNotCoverthePayoutProfit, "total liquidity %s, payout %s", fInfo.totalAvailableLiquidity, fInfo.payoutProfit)
	}

	return nil
}

// fulfill processes the participation and exposures in according to the expected bet amount to be fulfilled.
func (k Keeper) fulfill(
	ctx sdk.Context,
	fInfo *fulfillmentInfo,
	betAmountToFulfill,
	payoutProfitToFulfill sdk.Int,
) error {
	fInfo.inProcessItem.participationExposure.SetCurrentRound(betAmountToFulfill, payoutProfitToFulfill)
	fInfo.inProcessItem.participation.SetCurrentRound(&fInfo.inProcessItem.participationExposure, fInfo.oddsUID, betAmountToFulfill)

	fInfo.fulfillments = append(fInfo.fulfillments, bettypes.NewBetFulfillment(
		fInfo.inProcessItem.participation.ParticipantAddress,
		fInfo.inProcessItem.participation.Index,
		betAmountToFulfill,
		payoutProfitToFulfill,
	))

	// the amount has been fulfulled, so it should be subtracted from the bet amount of the
	fInfo.betAmount = fInfo.betAmount.Sub(betAmountToFulfill)

	// add the flfilled bet amount to the fulfillment amount tracker variable
	fInfo.fulfiledBetAmount = fInfo.fulfiledBetAmount.Add(betAmountToFulfill)

	// subtract the payout profit that is fulfilled from the initial payout profit
	// to prevent being calculated multiple times
	fInfo.payoutProfit = fInfo.payoutProfit.Sub(payoutProfitToFulfill.ToDec())

	// store the bet pair in the state
	participationBetPair := types.NewParticipationBetPair(
		fInfo.inProcessItem.participation.OrderBookUID,
		fInfo.betUID,
		fInfo.inProcessItem.participation.Index,
	)
	k.SetParticipationBetPair(ctx, participationBetPair, fInfo.betID)

	return nil
}

// prepareParticipationExposuresForNextRound prepares the participation exposures for the next round of queue process.
func (k Keeper) prepareParticipationExposuresForNextRound(ctx sdk.Context, fInfo *fulfillmentInfo, bookUID string) error {
	participationExposures, err := k.GetExposureByOrderBookAndParticipationIndex(ctx, bookUID, fInfo.inProcessItem.participation.Index)
	if err != nil {
		return err
	}

	// prepare the participation exposure map for the next round of calculations.
	for _, pe := range participationExposures {
		k.MoveToHistoricalParticipationExposure(ctx, pe)
		if fInfo.inProcessItem.participation.IsEligibleForNextRound() {
			newPe := pe.NextRound()
			k.SetParticipationExposure(ctx, newPe)
			if pe.OddsUID == fInfo.inProcessItem.participationExposure.OddsUID {
				fInfo.inProcessItem.participationExposure = newPe
				fInfo.fulfillmentMap.setExposure(pe.ParticipationIndex, fInfo.inProcessItem.participationExposure)
			}
		}
	}

	return nil
}

// prepareParticipationForNextRound prepares the participation for the next round of queue process.
func (k Keeper) prepareParticipationForNextRound(ctx sdk.Context, fInfo *fulfillmentInfo, notFilledExposures uint64) {
	// prepare participation for the next round
	fInfo.inProcessItem.participation.ResetForNextRound(notFilledExposures)
	fInfo.fulfillmentMap.setParticipation(fInfo.inProcessItem.participation)

	// store modified participation in the module state
	k.SetOrderBookParticipation(ctx, fInfo.inProcessItem.participation)
}

// prepareOddsExposuresForNextRound prepares the odds expsures for the next round of queue process.
func (k Keeper) prepareOddsExposuresForNextRound(ctx sdk.Context, fInfo *fulfillmentInfo, bookUID string) error {
	if fInfo.inProcessItem.participation.IsEligibleForNextRound() {
		boes, err := k.GetOddsExposuresByOrderBook(ctx, bookUID)
		if err != nil {
			return err
		}
		for _, boe := range boes {
			boe.FulfillmentQueue = append(boe.FulfillmentQueue, fInfo.inProcessItem.participation.Index)
			k.SetOrderBookOddsExposure(ctx, boe)
		}
		fInfo.fulfillmentQueue = append(fInfo.fulfillmentQueue, fInfo.inProcessItem.participation.Index)
	}

	return nil
}

// refreshQueueAndState refresh the fulfillment queue for the next round.
func (k Keeper) refreshQueueAndState(ctx sdk.Context, fInfo *fulfillmentInfo, book *types.OrderBook) error {
	fInfo.inProcessItem.participation.TrimCurrentRoundLiquidity()

	err := k.prepareParticipationExposuresForNextRound(ctx, fInfo, book.UID)
	if err != nil {
		return err
	}

	// prepare participation for the next round
	k.prepareParticipationForNextRound(ctx, fInfo, book.OddsCount)

	err = k.prepareOddsExposuresForNextRound(ctx, fInfo, book.UID)
	if err != nil {
		return err
	}

	return nil
}
