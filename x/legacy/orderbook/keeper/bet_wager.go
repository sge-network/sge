package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"

	bettypes "github.com/sge-network/sge/x/legacy/bet/types"
	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// ProcessWager processes bet placement
func (k Keeper) ProcessWager(
	ctx sdk.Context,
	betUID, bookUID, oddsUID string,
	maxLossMultiplier sdkmath.LegacyDec,
	betAmount sdkmath.Int,
	payoutProfit sdkmath.LegacyDec,
	bettorAddress sdk.AccAddress,
	betFee sdkmath.Int,
	oddsVal string, betID uint64,
	odds map[string]*bettypes.BetOddsCompact,
	oddUIDS []string,
) ([]*bettypes.BetFulfillment, error) {
	// get book data by its id
	book, found := k.GetOrderBook(ctx, bookUID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", bookUID)
	}
	// get exposures of the odds for which bet is placed
	bookExposure, found := k.GetOrderBookOddsExposure(ctx, bookUID, oddsUID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrOrderBookExposureNotFound, "%s , %s", bookUID, oddsUID)
	}

	fInfo, err := k.initFulfillmentInfo(
		ctx,
		betAmount,
		payoutProfit,
		betUID,
		betID,
		oddsUID,
		oddsVal,
		maxLossMultiplier,
		&book,
		odds,
		oddUIDS,
	)
	if err != nil {
		return nil, err
	}

	if err = k.fulfillBetByParticipationQueue(ctx, &fInfo, &book, bookExposure.FulfillmentQueue); err != nil {
		return nil, err
	}

	bookExposure.FulfillmentQueue = fInfo.updatedfulfillmentQueue
	k.SetOrderBookOddsExposure(ctx, bookExposure)

	// fund bet fee collector from bettor's account.
	if err := k.fund(bettypes.BetFeeCollectorFunder{}, ctx, bettorAddress, betFee); err != nil {
		return nil, err
	}

	// fund order book liquidity pool from bettor's account.
	if err := k.fund(types.OrderBookLiquidityFunder{}, ctx, bettorAddress, fInfo.fulfilledBetAmount); err != nil {
		return nil, err
	}

	return fInfo.fulfillments, nil
}

// fulfillBetByParticipationQueue fulfills the bet wagering payout using the participations
// that is stored in the state.
func (k Keeper) fulfillBetByParticipationQueue(
	ctx sdk.Context,
	fInfo *fulfillmentInfo,
	book *types.OrderBook,
	fulfillmentQueue []uint64,
) error {
	fInfo.fulfillmentQueue = fulfillmentQueue
	fInfo.updatedfulfillmentQueue = fulfillmentQueue

	// the decimal amount that is being lost in the bet amount calculation from payout profit
	truncatedBetAmount := sdkmath.LegacyNewDec(0)
	requeThreshold := sdkmath.NewIntFromUint64(k.GetRequeueThreshold(ctx))
	// continue until updatedFulfillmentQueue gets empty
	for fInfo.hasUnfulfilledQueueItem() {
		var err error
		// fill participation and exposure values
		fInfo.inProcessItem = fInfo.fulfillmentMap[fInfo.fulfillmentQueue[0]]

		// availableLiquidity is the available amount of tokens to be used from the participation exposure
		fInfo.inProcessItem.setAvailableLiquidity(fInfo.maxLossMultiplier)

		setFulfilled := false
		switch {
		case fInfo.inProcessItem.noLiquidityAvailable():
			setFulfilled = true
		case fInfo.notEnoughLiquidityAvailable():
			var betAmountToFulfill sdkmath.Int
			betAmountToFulfill, truncatedBetAmount, err = bettypes.CalculateBetAmountInt(
				fInfo.oddsVal,
				sdkmath.LegacyNewDecFromInt(fInfo.inProcessItem.availableLiquidity),
				truncatedBetAmount,
			)
			if err != nil {
				return err
			}

			// if the available liquidity is less than remaining payout profit that
			// need to be paid, we should use all of available liquidity pull for the calculations.
			k.fulfill(ctx, fInfo, betAmountToFulfill, fInfo.inProcessItem.availableLiquidity)
			setFulfilled = true
		default:
			// availableLiquidity is positive and more than remaining payout profit that
			// need to be paid, so we can cover all payout profits with available liquidity.
			// this case appends the last fulfillment
			if fInfo.inProcessItem.isLiquidityLessThanThreshold(requeThreshold, fInfo.payoutProfit.TruncateInt()) {
				setFulfilled = true
			}

			k.fulfill(ctx, fInfo, fInfo.betAmount, fInfo.payoutProfit.TruncateInt())
		}

		if setFulfilled {
			fInfo.setItemFulfilledAndRemove()
			if fInfo.inProcessItem.participation.IsEligibleForNextRoundPreLiquidityReduction() {
				eUpdate, err := fInfo.checkFullfillmentForOtherOdds(requeThreshold)
				if err != nil {
					return err
				}
				for _, exposure := range eUpdate {
					k.SetParticipationExposure(ctx, exposure)
				}
			}
		}

		k.SetParticipationExposure(ctx, fInfo.inProcessItem.participationExposure)
		k.SetOrderBookParticipation(ctx, fInfo.inProcessItem.participation)

		// if there are no more exposures to be filled
		if fInfo.inProcessItem.allExposureFulfilled() && fInfo.inProcessItem.participation.IsEligibleForNextRoundPreLiquidityReduction() {
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
		return types.ErrInsufficientLiquidityInOrderBook
	}

	return nil
}

// initFulfillmentInfo initializes the fulfillment info for the queue iteration process.
//
//nolint:nakedret
func (k Keeper) initFulfillmentInfo(
	ctx sdk.Context,
	betAmount sdkmath.Int,
	payoutProfit sdkmath.LegacyDec,
	betUID string,
	betID uint64,
	oddsUID string,
	oddsVal string,
	maxLossMultiplier sdkmath.LegacyDec,
	book *types.OrderBook,
	odds map[string]*bettypes.BetOddsCompact,
	oddUIDS []string,
) (
	fInfo fulfillmentInfo,
	err error,
) {
	fInfo = fulfillmentInfo{
		// bet specs
		betAmount:         betAmount,
		payoutProfit:      payoutProfit,
		betUID:            betUID,
		betID:             betID,
		oddsUID:           oddsUID,
		oddsVal:           oddsVal,
		maxLossMultiplier: maxLossMultiplier,

		//  in process maps
		fulfillmentMap: make(map[uint64]fulfillmentItem),

		// initialize the fulfilled bet amount with 0
		fulfilledBetAmount: sdkmath.NewInt(0),
		betOdds:            odds,
		oddUIDS:            oddUIDS,
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

	peMap, err := k.GetExposureByOrderBook(ctx, book.UID)
	if err != nil {
		return
	}
	if book.ParticipationCount != cast.ToUint64(len(peMap)) {
		err = sdkerrors.Wrapf(types.ErrParticipationExposuresNotFound, "%s, %s", book.UID, oddsUID)
		return
	}

	for _, bp := range bps {
		exposures, found := peMap[bp.Index]
		if !found {
			err = sdkerrors.Wrapf(types.ErrParticipationExposuresNotFound, "%s, %d", book.UID, bp.Index)
			return
		}
		fInfo.fulfillmentMap[bp.Index] = fulfillmentItem{
			participation: bp,
			allExposures:  exposures,
		}
	}

	for _, pe := range pes {
		item, found := fInfo.fulfillmentMap[pe.ParticipationIndex]
		if found {
			item.participationExposure = pe
			fInfo.fulfillmentMap[pe.ParticipationIndex] = item
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
			return sdkerrors.Wrapf(
				types.ErrParticipationExposureAlreadyFilled,
				"%d",
				participationIndex,
			)
		}
	}

	return nil
}

// fulfill processes the participation and exposures in according to the expected bet amount to be fulfilled.
func (k Keeper) fulfill(
	ctx sdk.Context,
	fInfo *fulfillmentInfo,
	betAmountToFulfill,
	payoutProfitToFulfill sdkmath.Int,
) {
	fInfo.inProcessItem.participationExposure.SetCurrentRound(betAmountToFulfill, payoutProfitToFulfill)
	fInfo.inProcessItem.participation.SetCurrentRound(
		&fInfo.inProcessItem.participationExposure,
		fInfo.oddsUID,
		betAmountToFulfill,
	)

	fInfo.fulfillments = append(fInfo.fulfillments, bettypes.NewBetFulfillment(
		fInfo.inProcessItem.participation.ParticipantAddress,
		fInfo.inProcessItem.participation.Index,
		betAmountToFulfill,
		payoutProfitToFulfill,
	))

	// the amount has been fulfilled, so it should be subtracted from the bet amount of the
	fInfo.betAmount = fInfo.betAmount.Sub(betAmountToFulfill)

	// add the fulfilled bet amount to the fulfillment amount tracker variable
	fInfo.fulfilledBetAmount = fInfo.fulfilledBetAmount.Add(betAmountToFulfill)

	// subtract the payout profit that is fulfilled from the initial payout profit
	// to prevent being calculated multiple times
	fInfo.payoutProfit = fInfo.payoutProfit.Sub(sdkmath.LegacyNewDecFromInt(payoutProfitToFulfill))

	// store the bet pair in the state
	participationBetPair := types.NewParticipationBetPair(
		fInfo.inProcessItem.participation.OrderBookUID,
		fInfo.betUID,
		fInfo.inProcessItem.participation.Index,
	)
	k.SetParticipationBetPair(ctx, participationBetPair, fInfo.betID)
}

// prepareParticipationExposuresForNextRound prepares the participation exposures for the next round of queue process.
func (k Keeper) prepareParticipationExposuresForNextRound(
	ctx sdk.Context,
	fInfo *fulfillmentInfo,
	bookUID string,
) error {
	participationExposures, err := k.GetExposureByOrderBookAndParticipationIndex(
		ctx,
		bookUID,
		fInfo.inProcessItem.participation.Index,
	)
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
				fInfo.fulfillmentMap.setExposure(
					pe.ParticipationIndex,
					fInfo.inProcessItem.participationExposure,
				)
			}
		}
	}

	return nil
}

// prepareParticipationForNextRound prepares the participation for the next round of queue process.
func (k Keeper) prepareParticipationForNextRound(
	ctx sdk.Context,
	fInfo *fulfillmentInfo,
	notFilledExposures uint64,
) {
	// prepare participation for the next round
	fInfo.inProcessItem.participation.ResetForNextRound(notFilledExposures)
	fInfo.fulfillmentMap.setParticipation(fInfo.inProcessItem.participation)

	// store modified participation in the module state
	k.SetOrderBookParticipation(ctx, fInfo.inProcessItem.participation)
}

// prepareOddsExposuresForNextRound prepares the odds exposures for the next round of queue process.
func (k Keeper) prepareOddsExposuresForNextRound(
	ctx sdk.Context,
	fInfo *fulfillmentInfo,
	bookUID string,
) error {
	if fInfo.inProcessItem.participation.IsEligibleForNextRound() {
		boes, err := k.GetOddsExposuresByOrderBook(ctx, bookUID)
		if err != nil {
			return err
		}
		for _, boe := range boes {
			if len(boe.FulfillmentQueue) != 0 && boe.FulfillmentQueue[0] == fInfo.inProcessItem.participation.Index {
				boe.FulfillmentQueue = boe.FulfillmentQueue[1:]
			}
			boe.FulfillmentQueue = append(boe.FulfillmentQueue, fInfo.inProcessItem.participation.Index)
			k.SetOrderBookOddsExposure(ctx, boe)
		}
		fInfo.updatedfulfillmentQueue = append(fInfo.updatedfulfillmentQueue, fInfo.inProcessItem.participation.Index)
	}

	return nil
}

// refreshQueueAndState refresh the fulfillment queue for the next round.
func (k Keeper) refreshQueueAndState(
	ctx sdk.Context,
	fInfo *fulfillmentInfo,
	book *types.OrderBook,
) error {
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

type fulfillmentInfo struct {
	betUID             string
	betID              uint64
	oddsUID            string
	oddsVal            string
	maxLossMultiplier  sdkmath.LegacyDec
	betAmount          sdkmath.Int
	payoutProfit       sdkmath.LegacyDec
	fulfilledBetAmount sdkmath.Int

	fulfillmentQueue        []uint64
	updatedfulfillmentQueue []uint64
	fulfillmentMap          fulfillmentMap
	inProcessItem           fulfillmentItem
	secondaryProcessItem    fulfillmentItem
	fulfillments            []*bettypes.BetFulfillment
	betOdds                 map[string]*bettypes.BetOddsCompact
	oddUIDS                 []string
}

//nolint:nakedret
func (fInfo *fulfillmentInfo) checkFullfillmentForOtherOdds(requeThreshold sdkmath.Int) (eUpdate []types.ParticipationExposure, err error) {
	if fInfo.inProcessItem.participation.ExposuresNotFilled == 0 {
		return
	}

	// check if other exposures are fulfilled according to new max loss multipliers
	for _, oddUID := range fInfo.oddUIDS {
		if oddUID == fInfo.oddsUID {
			continue
		}
		exposure, ok := fInfo.inProcessItem.allExposures[oddUID]
		if !ok {
			err = sdkerrors.Wrapf(types.ErrParticipationExposuresNotFound, "%s %d", oddUID, fInfo.inProcessItem.participation.Index)
			return
		}
		if exposure.IsFulfilled {
			continue
		}
		betOdd, ok := fInfo.betOdds[oddUID]
		if !ok {
			err = sdkerrors.Wrapf(bettypes.ErrOddsDataNotFound, "%s", oddUID)
			return
		}
		fInfo.secondaryProcessItem = fulfillmentItem{
			participation:         fInfo.inProcessItem.participation,
			participationExposure: *fInfo.inProcessItem.allExposures[oddUID],
		}

		// availableLiquidity is the available amount of tokens to be used from the participation exposure
		fInfo.secondaryProcessItem.setAvailableLiquidity(betOdd.MaxLossMultiplier)
		if fInfo.secondaryProcessItem.isLiquidityLessThanThreshold(requeThreshold, sdkmath.ZeroInt()) {
			fInfo.inProcessItem.setFulfilledSecondary(&fInfo.secondaryProcessItem)
			eUpdate = append(eUpdate, fInfo.secondaryProcessItem.participationExposure)
		}
	}
	return eUpdate, err
}

func (fInfo *fulfillmentInfo) setItemFulfilledAndRemove() {
	fInfo.inProcessItem.setFulfilled()
	fInfo.removeQueueItem()
}

func (fInfo *fulfillmentInfo) removeQueueItem() {
	fInfo.fulfillmentQueue = fInfo.fulfillmentQueue[1:]
	fInfo.updatedfulfillmentQueue = fInfo.updatedfulfillmentQueue[1:]
}

func (fInfo *fulfillmentInfo) hasUnfulfilledQueueItem() bool {
	return len(fInfo.fulfillmentQueue) > 0
}

func (fInfo *fulfillmentInfo) IsFulfilled() bool {
	// if the remaining payout is less than 1.00, means that the decimal part will be ignored
	return fInfo.payoutProfit.LT(sdkmath.LegacyOneDec()) || len(fInfo.fulfillmentQueue) == 0
}

func (fInfo *fulfillmentInfo) NoMoreLiquidityAvailable() bool {
	// if the remaining payout is less than 1.00, means that the decimal part will be ignored
	return fInfo.payoutProfit.GTE(sdkmath.LegacyOneDec())
}

func (fInfo *fulfillmentInfo) notEnoughLiquidityAvailable() bool {
	return fInfo.inProcessItem.availableLiquidity.LTE(fInfo.payoutProfit.TruncateInt())
}

type fulfillmentItem struct {
	availableLiquidity    sdkmath.Int
	participation         types.OrderBookParticipation
	participationExposure types.ParticipationExposure
	allExposures          map[string]*types.ParticipationExposure
}

func (fItem *fulfillmentItem) isLiquidityLessThanThreshold(threshold, adjAmount sdkmath.Int) bool {
	diff := fItem.availableLiquidity.Sub(adjAmount)
	return diff.LTE(threshold)
}

func (fItem *fulfillmentItem) noLiquidityAvailable() bool {
	return fItem.availableLiquidity.LTE(sdkmath.ZeroInt())
}

func (fItem *fulfillmentItem) setFulfilled() {
	fItem.participationExposure.IsFulfilled = true
	fItem.participation.ExposuresNotFilled--
}

func (fItem *fulfillmentItem) setFulfilledSecondary(sItem *fulfillmentItem) {
	sItem.participationExposure.IsFulfilled = true
	fItem.participation.ExposuresNotFilled--
}

func (fItem *fulfillmentItem) setAvailableLiquidity(maxLossMultiplier sdkmath.LegacyDec) {
	fItem.availableLiquidity = fItem.calcAvailableLiquidity(maxLossMultiplier)
}

func (fItem *fulfillmentItem) calcAvailableLiquidity(maxLossMultiplier sdkmath.LegacyDec) sdkmath.Int {
	return maxLossMultiplier.
		MulInt(fItem.participation.CurrentRoundLiquidity).
		Sub(sdkmath.LegacyNewDecFromInt(fItem.participationExposure.Exposure)).TruncateInt()
}

func (fItem *fulfillmentItem) allExposureFulfilled() bool {
	return fItem.participation.ExposuresNotFilled == 0
}

type fulfillmentMap map[uint64]fulfillmentItem

func (fMap fulfillmentMap) setParticipation(participation types.OrderBookParticipation) {
	fItem := fMap[participation.Index]
	fItem.participation = participation
	fMap[participation.Index] = fItem
}

func (fMap fulfillmentMap) setExposure(
	participationIndex uint64,
	exposure types.ParticipationExposure,
) {
	fItem := fMap[participationIndex]
	fItem.participationExposure = exposure
	fMap[participationIndex] = fItem
}
