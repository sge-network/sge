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
) (betFulfillments []*bettypes.BetFulfillment, err error) {
	// If lock exists, return error
	// Lock already exists means the bet is already placed for the given bet-uid
	if k.payoutLockExists(ctx, betUID) {
		err = sdkerrors.Wrapf(types.ErrLockAlreadyExists, "%s", betUID)
		return
	}

	// get book data by its id
	book, found := k.GetBook(ctx, bookUID)
	if !found {
		err = sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", bookUID)
		return
	}
	// get exposures of the odds user is placing the bet
	bookExposure, found := k.GetBookOddsExposure(ctx, bookUID, oddsUID)
	if !found {
		err = sdkerrors.Wrapf(types.ErrOrderBookExposureNotFound, "%s , %s", bookUID, oddsUID)
		return
	}

	// make maps for participation and its exposures
	participationMap := make(map[uint64]types.BookParticipation)
	participationExposureMap := make(map[uint64]types.ParticipationExposure)

	bps, err := k.GetParticipationsOfBook(ctx, book.ID)
	if err != nil {
		return
	}
	if book.ParticipationCount != cast.ToUint64(len(bps)) {
		err = sdkerrors.Wrapf(types.ErrBookParticipationsNotFound, "%s", bookUID)
		return
	}
	for _, bp := range bps {
		participationMap[bp.Index] = bp
	}

	pes, err := k.GetExposureByBookAndOdds(ctx, bookUID, oddsUID)
	if err != nil {
		return
	}
	if book.ParticipationCount != cast.ToUint64(len(pes)) {
		err = sdkerrors.Wrapf(types.ErrParticipationExposuresNotFound, "%s, %s", bookUID, oddsUID)
		return
	}
	for _, pe := range pes {
		participationExposureMap[pe.ParticipationIndex] = pe
	}

	remainingPayoutProfit, updatedFulfillmentQueue := payoutProfit, bookExposure.FulfillmentQueue
	fulfiledBetAmount := sdk.ZeroInt()
	truncatedBetAmount := sdk.NewDec(0)

	// continue until updatedFulfillmentQueue gets empty
	for len(updatedFulfillmentQueue) > 0 {
		var participation types.BookParticipation
		var participationExposure types.ParticipationExposure
		participationIndex := updatedFulfillmentQueue[0]

		getQueueItemParticipation := func() error {
			participation, found = participationMap[participationIndex]
			if !found {
				return sdkerrors.Wrapf(types.ErrBookParticipationNotFound, "%s, %d", bookUID, participationIndex)
			}
			participationExposure, found = participationExposureMap[participationIndex]
			if !found {
				return sdkerrors.Wrapf(types.ErrParticipationExposureNotFound, "%s, %d", bookUID, participationIndex)
			}
			if participationExposure.IsFulfilled {
				return sdkerrors.Wrapf(types.ErrParticipationExposureAlreadyFilled, "%s, %d", bookUID, participationIndex)
			}
			return nil
		}

		removeQueueItem := func() {
			updatedFulfillmentQueue = updatedFulfillmentQueue[1:]
		}

		setFulfilled := func() {
			participationExposure.IsFulfilled = true
			participation.ExposuresNotFilled--
		}

		processBetFulfillment := func(payoutToFulfill sdk.Int, isLastFulfillment bool) error {
			var betAmountToFulfill sdk.Int
			if isLastFulfillment {
				// the bet amount
				betAmountToFulfill = betAmount
			} else {
				expectedBetAmountDec, err := bettypes.CalculateBetAmount(oddsType, oddsVal, payoutToFulfill.ToDec())
				if err != nil {
					return err
				}
				expectedBetAmountDec = expectedBetAmountDec.Add(truncatedBetAmount)
				betAmountToFulfill = expectedBetAmountDec.TruncateInt()
				truncatedBetAmount = truncatedBetAmount.Add(expectedBetAmountDec.Sub(betAmountToFulfill.ToDec()))
			}

			if !isLastFulfillment {
				setFulfilled()
			}
			participationExposure.Exposure = participationExposure.Exposure.Add(payoutToFulfill)
			participationExposure.BetAmount = participationExposure.BetAmount.Add(betAmountToFulfill)
			participation.TotalBetAmount = participation.TotalBetAmount.Add(betAmountToFulfill)
			participation.CurrentRoundTotalBetAmount = participation.CurrentRoundTotalBetAmount.Add(betAmountToFulfill)

			maxLoss := participationExposure.Exposure.Sub(participation.CurrentRoundTotalBetAmount).Add(participationExposure.BetAmount)
			switch {
			case participation.CurrentRoundMaxLoss.IsNil():
				participation.CurrentRoundMaxLoss = maxLoss
				participation.CurrentRoundMaxLossOddsUID = oddsUID
			case participation.CurrentRoundMaxLossOddsUID == oddsUID:
				participation.CurrentRoundMaxLoss = maxLoss
			default:
				originalMaxLoss := participation.CurrentRoundMaxLoss.Sub(betAmountToFulfill)
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
				PayoutAmount:       payoutToFulfill,
			})
			betAmount = betAmount.Sub(betAmountToFulfill)

			fulfiledBetAmount = fulfiledBetAmount.Add(betAmountToFulfill)
			remainingPayoutProfit = remainingPayoutProfit.Sub(payoutToFulfill.ToDec())

			participationBetPair := types.NewParticipationBetPair(participation.BookUID, betUID, participation.Index)
			k.SetParticipationBetPair(ctx, participationBetPair, betID)
			return nil
		}

		refreshQueueAndState := func() (err error) {
			maxLoss := sdk.MaxInt(sdk.ZeroInt(), participation.CurrentRoundMaxLoss)
			participation.CurrentRoundLiquidity = participation.CurrentRoundLiquidity.Sub(maxLoss)

			// check if there is more liquidity amount
			eligibleForNextRound := participation.CurrentRoundLiquidity.GT(sdk.ZeroInt())

			participationExposures, err := k.GetExposureByBookAndParticipationIndex(ctx, bookUID, participationIndex)
			if err != nil {
				return
			}
			for _, pe := range participationExposures {
				k.MoveToHistoricalParticipationExposure(ctx, pe)
				if eligibleForNextRound {
					newPe := types.NewParticipationExposure(book.ID, pe.OddsUID, sdk.ZeroInt(), sdk.ZeroInt(), pe.ParticipationIndex, pe.Round+1, false)
					k.SetParticipationExposure(ctx, newPe)
					if pe.OddsUID == participationExposure.OddsUID {
						participationExposure = newPe
						participationExposureMap[pe.ParticipationIndex] = participationExposure
					}
				}
			}

			participation.ExposuresNotFilled = book.OddsCount
			participation.CurrentRoundTotalBetAmount = sdk.ZeroInt()
			participation.MaxLoss = participation.MaxLoss.Add(participation.CurrentRoundMaxLoss)
			participation.CurrentRoundMaxLoss = sdk.ZeroInt()
			participationMap[participation.Index] = participation
			k.SetBookParticipation(ctx, participation)

			if eligibleForNextRound {
				var boes []types.BookOddsExposure
				boes, err = k.GetOddsExposuresByBook(ctx, bookUID)
				if err != nil {
					return
				}
				for _, boe := range boes {
					boe.FulfillmentQueue = append(boe.FulfillmentQueue, participationIndex)
					if boe.OddsUID == participationExposure.OddsUID {
						bookExposure = boe
					}

					k.SetBookOddsExposure(ctx, boe)
				}
				updatedFulfillmentQueue = append(updatedFulfillmentQueue, participationIndex)
				bookExposure.FulfillmentQueue = updatedFulfillmentQueue
			}

			return nil
		}

		// fill participation and exposure values
		err = getQueueItemParticipation()
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
		case availableLiquidty.ToDec().LTE(remainingPayoutProfit):
			// if the available liquidity is less than remaining payout profit that
			// need to be paid, we should use all of available liquidity pull for the calculations.
			err := processBetFulfillment(availableLiquidty, false)
			if err != nil {
				return betFulfillments, err
			}
			removeQueueItem()
		default:
			// availableLiquidty is positive and more than remaining payout profit that
			// need to be paid, so we can cover all of payout profits with available liquidity.
			// this case appends the last fulfillment
			err := processBetFulfillment(remainingPayoutProfit.TruncateInt(), true)
			if err != nil {
				return betFulfillments, err
			}
		}

		k.SetParticipationExposure(ctx, participationExposure)
		k.SetBookParticipation(ctx, participation)

		// if there are no more exposures to be filled
		if participation.ExposuresNotFilled == 0 {
			err = refreshQueueAndState()
			if err != nil {
				return
			}
		}

		// if the remaining payout is less than 1.00, means that the decimal part will be ignored
		if remainingPayoutProfit.LT(sdk.OneDec()) || len(updatedFulfillmentQueue) == 0 {
			break
		}
	}

	if remainingPayoutProfit.GTE(sdk.OneDec()) {
		return betFulfillments, sdkerrors.Wrapf(types.ErrInternalProcessingBet, "insufficient liquidity in order book")
	}

	bookExposure.FulfillmentQueue = updatedFulfillmentQueue
	k.SetBookOddsExposure(ctx, bookExposure)

	// Transfer bet fee from bettor to the `bet` module account
	err = k.transferFundsFromUserToModule(ctx, bettorAddress, bettypes.ModuleName, betFee)
	if err != nil {
		return betFulfillments, err
	}

	// Transfer bet amount from bettor to `book_liquidity_pool` Account
	err = k.transferFundsFromUserToModule(ctx, bettorAddress, types.BookLiquidityName, fulfiledBetAmount)
	if err != nil {
		return betFulfillments, err
	}

	// Create a unique lock in the Payout Store for the bet
	k.SetPayoutLock(ctx, betUID)

	return betFulfillments, nil
}
