package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cast"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

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

	betFulfillments, updatedFulfillmentQueue, fulfiledBetAmount, err := k.fulFillQueueBets(ctx,
		betUID,
		betID,
		oddsUID,
		oddsVal,
		oddsType,
		betAmount,
		payoutProfit,
		maxLossMultiplier,
		&book,
		&bookExposure,
	)
	if err != nil {
		return nil, err
	}

	bookExposure.FulfillmentQueue = updatedFulfillmentQueue
	k.SetOrderBookOddsExposure(ctx, bookExposure)

	// Transfer bet fee from bettor to the `bet` module account
	if err = k.transferFundsFromAccountToModule(ctx, bettorAddress, bettypes.ModuleName, betFee); err != nil {
		return nil, err
	}

	// Transfer bet amount from bettor to `book_liquidity_pool` Account
	if err = k.transferFundsFromAccountToModule(ctx, bettorAddress, types.OrderBookLiquidityName, fulfiledBetAmount); err != nil {
		return nil, err
	}

	// Create a unique lock in the Payout Store for the bet
	k.SetPayoutLock(ctx, betUID)

	return betFulfillments, nil
}

func (k Keeper) fulFillQueueBets(
	ctx sdk.Context,
	betUID string, betID uint64,
	oddsUID, oddsVal string, oddsType bettypes.OddsType,
	betAmount sdk.Int,
	payoutProfit sdk.Dec,
	maxLossMultiplier sdk.Dec,
	book *types.OrderBook,
	bookExposure *types.OrderBookOddsExposure,
) (
	betFulfillments []*bettypes.BetFulfillment,
	updatedQueue []uint64,
	fulfiledBetAmount sdk.Int,
	err error,
) {
	participationMap, participationExposureMap, err := k.getExposuresMap(ctx, oddsUID, book)
	if err != nil {
		return
	}

	updatedQueue = bookExposure.FulfillmentQueue

	// the decimal amount that is being lost in the bet amount calculation from payout profit
	truncatedBetAmount := sdk.NewDec(0)

	// continue until updatedFulfillmentQueue gets empty
	for len(updatedQueue) > 0 {
		var participation types.OrderBookParticipation
		var participationExposure types.ParticipationExposure
		participationIndex := updatedQueue[0]

		getQueueItemInfo := func() error {
			var found bool
			participation, found = participationMap[participationIndex]
			if !found {
				return sdkerrors.Wrapf(types.ErrOrderBookParticipationNotFound, "%s, %d", book.UID, participationIndex)
			}
			participationExposure, found = participationExposureMap[participationIndex]
			if !found {
				return sdkerrors.Wrapf(types.ErrParticipationExposureNotFound, "%s, %d", book.UID, participationIndex)
			}
			if participationExposure.IsFulfilled {
				return sdkerrors.Wrapf(types.ErrParticipationExposureAlreadyFilled, "%s, %d", book.UID, participationIndex)
			}
			return nil
		}

		removeQueueItem := func() {
			updatedQueue = updatedQueue[1:]
		}

		setFulfilled := func() {
			participationExposure.IsFulfilled = true
			participation.ExposuresNotFilled--
		}

		fulfill := func(payoutProfitToFulfill sdk.Int, isLastFulfillment bool) error {
			var betAmountToFulfill sdk.Int

			// if the fulfillment is in the last state,
			// the bet amount should be fulfilled totally
			// because there is no decimal truncation.
			if isLastFulfillment {
				// the bet amount
				betAmountToFulfill = betAmount
			} else {
				// findout what would be the bet amount according to the pyout profit
				expectedBetAmountDec, err := bettypes.CalculateBetAmount(oddsType, oddsVal, payoutProfitToFulfill.ToDec())
				if err != nil {
					return err
				}
				// add previous loop truncated value to the calculated bet amount
				expectedBetAmountDec = expectedBetAmountDec.Add(truncatedBetAmount)

				// we need for the bet amount to be of type sdk.Int
				// so the truncation in inevitable
				betAmountToFulfill = expectedBetAmountDec.TruncateInt()

				// save the truncated amount in the calculations for the next loop
				truncatedBetAmount = truncatedBetAmount.Add(expectedBetAmountDec.Sub(betAmountToFulfill.ToDec()))
			}

			if !isLastFulfillment {
				setFulfilled()
			}
			// add the payout profit to the
			participationExposure.Exposure = participationExposure.Exposure.Add(payoutProfitToFulfill)

			// add the bet amount that is being fulfilled to the exposure and participation
			participationExposure.BetAmount = participationExposure.BetAmount.Add(betAmountToFulfill)
			participation.TotalBetAmount = participation.TotalBetAmount.Add(betAmountToFulfill)
			participation.CurrentRoundTotalBetAmount = participation.CurrentRoundTotalBetAmount.Add(betAmountToFulfill)

			// max loss is the maximum amount that an exposure may lose.
			maxLoss := participationExposure.CalculateMaxLoss(participation.CurrentRoundTotalBetAmount)
			switch {
			case participation.CurrentRoundMaxLoss.IsNil():
				participation.CurrentRoundMaxLoss = maxLoss
				participation.CurrentRoundMaxLossOddsUID = oddsUID
			case participation.CurrentRoundMaxLossOddsUID == oddsUID:
				participation.CurrentRoundMaxLoss = maxLoss
			default:
				originalMaxLoss := participation.CalculateMaxLoss(betAmountToFulfill)
				if maxLoss.GT(originalMaxLoss) {
					participation.CurrentRoundMaxLoss = maxLoss
					participation.CurrentRoundMaxLossOddsUID = oddsUID
				} else {
					participation.CurrentRoundMaxLoss = originalMaxLoss
				}
			}

			betFulfillments = append(betFulfillments, &bettypes.BetFulfillment{
				ParticipantAddress: participation.ParticipantAddress,
				ParticipationIndex: participation.Index,
				BetAmount:          betAmountToFulfill,
				PayoutProfit:       payoutProfitToFulfill,
			})

			// the amount has been fulfulled, so it should be subtracted from the bet amount of the
			betAmount = betAmount.Sub(betAmountToFulfill)

			// add the flfilled bet amount to the fulfillment amount tracker variable
			fulfiledBetAmount = fulfiledBetAmount.Add(betAmountToFulfill)

			// subtract the payout profit that is fulfilled from the initial payout profit
			// to prevent being calculated multiple times
			payoutProfit = payoutProfit.Sub(payoutProfitToFulfill.ToDec())

			// store the bet pair in the state
			participationBetPair := types.NewParticipationBetPair(participation.OrderBookUID, betUID, participation.Index)
			k.SetParticipationBetPair(ctx, participationBetPair, betID)

			return nil
		}

		refreshQueueAndState := func() (err error) {
			maxLoss := sdk.MaxInt(sdk.ZeroInt(), participation.CurrentRoundMaxLoss)
			participation.CurrentRoundLiquidity = participation.CurrentRoundLiquidity.Sub(maxLoss)

			// check if there is more liquidity amount
			eligibleForNextRound := participation.CurrentRoundLiquidity.GT(sdk.ZeroInt())

			participationExposures, err := k.GetExposureByOrderBookAndParticipationIndex(ctx, book.UID, participationIndex)
			if err != nil {
				return
			}

			// prepare the participation exposure map for the next round of calculations.
			for _, pe := range participationExposures {
				k.MoveToHistoricalParticipationExposure(ctx, pe)
				if eligibleForNextRound {
					newPe := types.NewParticipationExposure(book.UID, pe.OddsUID, sdk.ZeroInt(), sdk.ZeroInt(), pe.ParticipationIndex, pe.Round+1, false)
					k.SetParticipationExposure(ctx, newPe)
					if pe.OddsUID == participationExposure.OddsUID {
						participationExposure = newPe
						participationExposureMap[pe.ParticipationIndex] = participationExposure
					}
				}
			}

			// prepare participation for the next round
			participation.ExposuresNotFilled = book.OddsCount
			participation.CurrentRoundTotalBetAmount = sdk.ZeroInt()
			participation.MaxLoss = participation.MaxLoss.Add(participation.CurrentRoundMaxLoss)
			participation.CurrentRoundMaxLoss = sdk.ZeroInt()
			participationMap[participation.Index] = participation
			k.SetOrderBookParticipation(ctx, participation)

			if eligibleForNextRound {
				var boes []types.OrderBookOddsExposure
				boes, err = k.GetOddsExposuresByOrderBook(ctx, book.UID)
				if err != nil {
					return
				}
				for i, boe := range boes {
					boe.FulfillmentQueue = append(boe.FulfillmentQueue, participationIndex)
					if boe.OddsUID == participationExposure.OddsUID {
						// use the index to prevent implicit memory aliasing.
						bookExposure = &boes[i]
					}

					k.SetOrderBookOddsExposure(ctx, boe)
				}
				updatedQueue = append(updatedQueue, participationIndex)
				bookExposure.FulfillmentQueue = updatedQueue
			}

			return nil
		}

		// fill participation and exposure values
		err = getQueueItemInfo()
		if err != nil {
			return
		}

		// availableLiquidty is the available amount of tokens to be used from the participation exposure
		availableLiquidty := maxLossMultiplier.
			MulInt(participation.CurrentRoundLiquidity).
			Sub(sdk.NewDecFromInt(participationExposure.Exposure)).TruncateInt()

		switch {
		case availableLiquidty.LTE(sdk.ZeroInt()):
			setFulfilled()
			removeQueueItem()
		case availableLiquidty.ToDec().LTE(payoutProfit):
			// if the available liquidity is less than remaining payout profit that
			// need to be paid, we should use all of available liquidity pull for the calculations.
			err = fulfill(availableLiquidty, false)
			if err != nil {
				return
			}
			removeQueueItem()
		default:
			// availableLiquidty is positive and more than remaining payout profit that
			// need to be paid, so we can cover all of payout profits with available liquidity.
			// this case appends the last fulfillment
			err = fulfill(payoutProfit.TruncateInt(), true)
			if err != nil {
				return
			}
		}

		k.SetParticipationExposure(ctx, participationExposure)
		k.SetOrderBookParticipation(ctx, participation)

		// if there are no more exposures to be filled
		if participation.ExposuresNotFilled == 0 {
			err = refreshQueueAndState()
			if err != nil {
				return
			}
		}

		// if the remaining payout is less than 1.00, means that the decimal part will be ignored
		if payoutProfit.LT(sdk.OneDec()) || len(updatedQueue) == 0 {
			break
		}
	}

	if payoutProfit.GTE(sdk.OneDec()) {
		err = sdkerrors.Wrapf(types.ErrInternalProcessingBet, "insufficient liquidity in order book")
		return
	}

	return betFulfillments, []uint64{}, sdk.Int{}, nil
}

func (k Keeper) getExposuresMap(
	ctx sdk.Context,
	oddsUID string,
	book *types.OrderBook,
) (
	participationMap map[uint64]types.OrderBookParticipation,
	participationExposureMap map[uint64]types.ParticipationExposure,
	err error,
) {
	participationMap = make(map[uint64]types.OrderBookParticipation)
	participationExposureMap = make(map[uint64]types.ParticipationExposure)

	bps, err := k.GetParticipationsOfOrderBook(ctx, book.UID)
	if err != nil {
		return
	}
	if book.ParticipationCount != cast.ToUint64(len(bps)) {
		err = sdkerrors.Wrapf(types.ErrBookParticipationsNotFound, "%s", book.UID)
		return
	}
	for _, bp := range bps {
		participationMap[bp.Index] = bp
	}

	pes, err := k.GetExposureByOrderBookAndOdds(ctx, book.UID, oddsUID)
	if err != nil {
		return
	}
	if book.ParticipationCount != cast.ToUint64(len(pes)) {
		err = sdkerrors.Wrapf(types.ErrParticipationExposuresNotFound, "%s, %s", book.UID, oddsUID)
		return
	}
	for _, pe := range pes {
		participationExposureMap[pe.ParticipationIndex] = pe
	}

	return participationMap, participationExposureMap, nil
}
