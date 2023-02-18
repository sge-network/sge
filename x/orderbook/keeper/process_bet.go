package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/orderbook/types"
)

// ProcessBetPlacement processes bet placement
func (k Keeper) ProcessBetPlacement(
	ctx sdk.Context,
	uniqueLock, bookID, oddsID string,
	maxLossMultiplier sdk.Dec,
	payoutProfit sdk.Int,
	bettorAddress sdk.AccAddress,
	betFee sdk.Int,
	betAmount sdk.Int,
	oddsType bettypes.OddsType,
	oddsVal string, betID uint64,
) ([]*bettypes.BetFullfillment, error) {
	betFullfillment := []*bettypes.BetFullfillment{}

	// If lock exists, return error
	// Lock already exists means the bet is already placed for the given bet-uid
	if k.payoutLockExists(ctx, uniqueLock) {
		return betFullfillment, sdkerrors.Wrapf(types.ErrLockAlreadyExists, "%s", uniqueLock)
	}

	// Check bet fullfillment
	book, found := k.GetBook(ctx, bookID)
	if !found {
		return betFullfillment, sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", bookID)
	}
	bookExposure, found := k.GetBookOddsExposure(ctx, bookID, oddsID)
	if !found {
		return betFullfillment, sdkerrors.Wrapf(types.ErrOrderBookExposureNotFound, "%s , %s", bookID, oddsID)
	}

	// Get participants for a book
	participantMap := make(map[uint64]types.BookParticipant)
	participantExposureMap := make(map[uint64]types.ParticipantExposure)

	bps := k.GetParticipantsByBook(ctx, book.ID)
	if int(book.Participants) != len(bps) {
		return betFullfillment, sdkerrors.Wrapf(types.ErrBookParticipantsNotFound, "%s", bookID)
	}
	for _, bp := range bps {
		participantMap[bp.ParticipantNumber] = bp
	}

	pes := k.GetExposureByBookAndOdds(ctx, bookID, oddsID)
	if int(book.Participants) != len(pes) {
		return betFullfillment, sdkerrors.Wrapf(types.ErrParticipantExposuresNotFound, "%s, %s", bookID, oddsID)
	}
	for _, pe := range pes {
		participantExposureMap[pe.ParticipantNumber] = pe
	}

	remainingPayoutProfit, updatedFullfillmentQueue := payoutProfit, bookExposure.FullfillmentQueue
	continueLoop, calculatedBetAmount := len(updatedFullfillmentQueue) > 0, sdk.ZeroInt()
	for continueLoop {
		pn := updatedFullfillmentQueue[0]
		participant, found := participantMap[pn]
		if !found {
			return betFullfillment, sdkerrors.Wrapf(types.ErrBookParticipantNotFound, "%s, %d", bookID, pn)
		}
		participantExposure, found := participantExposureMap[pn]
		if !found {
			return betFullfillment, sdkerrors.Wrapf(types.ErrParticipantExposureNotFound, "%s, %d", bookID, pn)
		}
		if participantExposure.IsFullfilled {
			return betFullfillment, sdkerrors.Wrapf(types.ErrParticipantExposureAlreadyFilled, "%s, %d", bookID, pn)
		}

		availableLiquidty := maxLossMultiplier.MulInt(participant.CurrentRoundLiquidity).TruncateInt().Sub(participantExposure.Exposure)
		switch {
		case availableLiquidty.LTE(sdk.ZeroInt()):
			participantExposure.IsFullfilled = true
			participant.ExposuresNotFilled--
			updatedFullfillmentQueue = updatedFullfillmentQueue[1:]
		case availableLiquidty.LTE(remainingPayoutProfit):
			participantExposure.Exposure = participantExposure.Exposure.Add(availableLiquidty)
			betAmount, err := bettypes.CalculateBetAmount(oddsType, oddsVal, availableLiquidty)
			if err != nil {
				return betFullfillment, err
			}
			calculatedBetAmount = calculatedBetAmount.Add(betAmount)
			participantExposure.BetAmount = participantExposure.BetAmount.Add(betAmount)
			participantExposure.IsFullfilled = true
			participant.ExposuresNotFilled--
			participant.TotalBetAmount = participant.TotalBetAmount.Add(betAmount)
			participant.CurrentRoundTotalBetAmount = participant.CurrentRoundTotalBetAmount.Add(betAmount)
			updatedFullfillmentQueue = updatedFullfillmentQueue[1:]

			remainingPayoutProfit = remainingPayoutProfit.Sub(availableLiquidty)
			maxLoss := participantExposure.Exposure.Sub(participant.CurrentRoundTotalBetAmount).Add(participantExposure.BetAmount)
			switch {
			case participant.CurrentRoundMaxLoss.IsNil():
				participant.CurrentRoundMaxLoss = maxLoss
				participant.CurrentRoundMaxLossOddsID = oddsID
			case participant.CurrentRoundMaxLossOddsID == oddsID:
				participant.CurrentRoundMaxLoss = maxLoss
			default:
				originalMaxLoss := participant.CurrentRoundMaxLoss.Sub(betAmount)
				if maxLoss.GT(originalMaxLoss) {
					participant.CurrentRoundMaxLoss = maxLoss
					participant.CurrentRoundMaxLossOddsID = oddsID
				} else {
					participant.CurrentRoundMaxLoss = originalMaxLoss
				}
			}

			betFullfillment = append(betFullfillment, &bettypes.BetFullfillment{
				ParticipantAddress: participant.ParticipantAddress,
				ParticipantNumber:  participant.ParticipantNumber,
				BetAmount:          betAmount,
				PayoutAmount:       availableLiquidty,
			})
			participantBetPair := types.NewParticipantBetPair(participant.BookID, uniqueLock, participant.ParticipantNumber, betID)
			k.SetParticipantBetPair(ctx, participantBetPair)
		default:
			participantExposure.Exposure = participantExposure.Exposure.Add(remainingPayoutProfit)
			betAmount, err := bettypes.CalculateBetAmount(oddsType, oddsVal, remainingPayoutProfit)
			if err != nil {
				return betFullfillment, err
			}
			calculatedBetAmount = calculatedBetAmount.Add(betAmount)
			participantExposure.BetAmount = participantExposure.BetAmount.Add(betAmount)
			participant.TotalBetAmount = participant.TotalBetAmount.Add(betAmount)
			participant.CurrentRoundTotalBetAmount = participant.CurrentRoundTotalBetAmount.Add(betAmount)
			maxLoss := participantExposure.Exposure.Sub(participant.CurrentRoundTotalBetAmount).Add(participantExposure.BetAmount)
			switch {
			case participant.CurrentRoundMaxLoss.IsNil():
				participant.CurrentRoundMaxLoss = maxLoss
				participant.CurrentRoundMaxLossOddsID = oddsID
			case participant.CurrentRoundMaxLossOddsID == oddsID:
				participant.CurrentRoundMaxLoss = maxLoss
			default:
				originalMaxLoss := participant.CurrentRoundMaxLoss.Sub(betAmount)
				if maxLoss.GT(originalMaxLoss) {
					participant.CurrentRoundMaxLoss = maxLoss
					participant.CurrentRoundMaxLossOddsID = oddsID
				} else {
					participant.CurrentRoundMaxLoss = originalMaxLoss
				}
			}

			betFullfillment = append(betFullfillment, &bettypes.BetFullfillment{
				ParticipantAddress: participant.ParticipantAddress,
				ParticipantNumber:  participant.ParticipantNumber,
				BetAmount:          betAmount,
				PayoutAmount:       remainingPayoutProfit,
			})
			remainingPayoutProfit = remainingPayoutProfit.Sub(remainingPayoutProfit)
			participantBetPair := types.NewParticipantBetPair(participant.BookID, uniqueLock, participant.ParticipantNumber, betID)
			k.SetParticipantBetPair(ctx, participantBetPair)
		}

		k.SetParticipantExposure(ctx, participantExposure)
		k.SetBookParticipant(ctx, participant)

		if participant.ExposuresNotFilled == 0 {
			// add back to queue
			maxLoss := sdk.MaxInt(sdk.ZeroInt(), participant.CurrentRoundMaxLoss)
			participant.CurrentRoundLiquidity = participant.CurrentRoundLiquidity.Sub(maxLoss)
			eligibleForNextRound := participant.CurrentRoundLiquidity.GT(sdk.ZeroInt())

			participantExposures := k.GetExposureByBookAndParticipantNumber(ctx, bookID, pn)
			for _, pe := range participantExposures {
				k.MoveToHistoricalParticipantExposure(ctx, pe)
				if eligibleForNextRound {
					newPe := types.NewParticipantExposure(book.ID, pe.OddsID, sdk.ZeroInt(), sdk.ZeroInt(), pe.ParticipantNumber, pe.Round+1, false)
					k.SetParticipantExposure(ctx, newPe)
					if pe.OddsID == participantExposure.OddsID {
						participantExposure = newPe
						participantExposureMap[pe.ParticipantNumber] = participantExposure
					}
				}
			}

			participant.ExposuresNotFilled = book.NumberOfOdds
			participant.CurrentRoundTotalBetAmount = sdk.ZeroInt()
			participant.MaxLoss = participant.MaxLoss.Add(participant.CurrentRoundMaxLoss)
			participant.CurrentRoundMaxLoss = sdk.ZeroInt()
			participantMap[participant.ParticipantNumber] = participant
			k.SetBookParticipant(ctx, participant)

			if eligibleForNextRound {
				boes := k.GetOddsExposuresByBook(ctx, bookID)
				for _, boe := range boes {
					boe.FullfillmentQueue = append(boe.FullfillmentQueue, pn)
					if boe.OddsID == participantExposure.OddsID {
						bookExposure = boe
					}

					k.SetBookOddsExposure(ctx, boe)
				}
				updatedFullfillmentQueue = append(updatedFullfillmentQueue, pn)
				bookExposure.FullfillmentQueue = updatedFullfillmentQueue
			}
		}

		if remainingPayoutProfit.LTE(sdk.ZeroInt()) || len(updatedFullfillmentQueue) == 0 {
			break
		}
	}

	if remainingPayoutProfit.GT(sdk.ZeroInt()) {
		return betFullfillment, sdkerrors.Wrapf(types.ErrInternalProcessingBet, "insufficient liquidity in order book")
	}

	bookExposure.FullfillmentQueue = updatedFullfillmentQueue
	k.SetBookOddsExposure(ctx, bookExposure)

	// Transfer bet fee from bettor to the `bet` module account
	err := k.transferFundsFromUserToModule(ctx, bettorAddress, bettypes.ModuleName, betFee)
	if err != nil {
		return betFullfillment, err
	}

	// Transfer bet amount from bettor to `book_liquidity_pool` Account
	err = k.transferFundsFromUserToModule(ctx, bettorAddress, types.BookLiquidityName, calculatedBetAmount)
	if err != nil {
		return betFullfillment, err
	}

	// Create a unique lock in the Payout Store for the bet
	k.setPayoutLock(ctx, uniqueLock)

	return betFullfillment, nil
}
