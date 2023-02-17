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
	uniqueLock, bookId, oddId string,
	maxLossMultiplier sdk.Dec,
	payoutProfit sdk.Int,
	bettorAddress sdk.AccAddress,
	betFee sdk.Int,
	betAmount sdk.Int,
	oddsType bettypes.OddsType,
	oddsVal string, betId uint64,
) (error, []*bettypes.BetFullfillment) {
	betFullfillment := []*bettypes.BetFullfillment{}

	// If lock exists, return error
	// Lock already exists means the bet is already placed for the given bet-uid
	if k.payoutLockExists(ctx, uniqueLock) {
		return sdkerrors.Wrapf(types.ErrLockAlreadyExists, "%s", uniqueLock), betFullfillment
	}

	// Check bet fullfillment
	book, found := k.GetBook(ctx, bookId)
	if !found {
		return sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", bookId), betFullfillment
	}
	bookExposure, found := k.GetBookOddExposure(ctx, bookId, oddId)
	if !found {
		return sdkerrors.Wrapf(types.ErrOrderBookExposureNotFound, "%s , %s", bookId, oddId), betFullfillment
	}

	// Get participants for a book
	participantMap := make(map[uint64]types.BookParticipant)
	participantExposureMap := make(map[uint64]types.ParticipantExposure)

	bps := k.GetParticipantsByBook(ctx, book.Id)
	if int(book.Participants) != len(bps) {
		return sdkerrors.Wrapf(types.ErrBookParticipantsNotFound, "%s", bookId), betFullfillment
	}
	for _, bp := range bps {
		participantMap[bp.ParticipantNumber] = bp
	}

	pes := k.GetExposureByBookAndOdd(ctx, bookId, oddId)
	if int(book.Participants) != len(pes) {
		return sdkerrors.Wrapf(types.ErrParticipantExposuresNotFound, "%s, %s", bookId, oddId), betFullfillment
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
			return sdkerrors.Wrapf(types.ErrBookParticipantNotFound, "%s, %d", bookId, pn), betFullfillment
		}
		participantExposure, found := participantExposureMap[pn]
		if !found {
			return sdkerrors.Wrapf(types.ErrParticipantExposureNotFound, "%s, %d", bookId, pn), betFullfillment
		}
		if participantExposure.IsFullfilled {
			return sdkerrors.Wrapf(types.ErrParticipantExposureAlreadyFilled, "%s, %d", bookId, pn), betFullfillment
		}

		availableLiquidty := maxLossMultiplier.MulInt(participant.CurrentRoundLiquidity).TruncateInt().Sub(participantExposure.Exposure)
		if availableLiquidty.LTE(sdk.ZeroInt()) {
			participantExposure.IsFullfilled = true
			participant.ExposuresNotFilled -= 1
			updatedFullfillmentQueue = updatedFullfillmentQueue[1:]
		} else if availableLiquidty.LTE(remainingPayoutProfit) {
			participantExposure.Exposure = participantExposure.Exposure.Add(availableLiquidty)
			betAmount, err := bettypes.CalculateBetAmount(oddsType, oddsVal, availableLiquidty)
			if err != nil {
				return err, betFullfillment
			}
			calculatedBetAmount = calculatedBetAmount.Add(betAmount)
			participantExposure.BetAmount = participantExposure.BetAmount.Add(betAmount)
			participantExposure.IsFullfilled = true
			participant.ExposuresNotFilled -= 1
			participant.TotalBetAmount = participant.TotalBetAmount.Add(betAmount)
			participant.CurrentRoundTotalBetAmount = participant.CurrentRoundTotalBetAmount.Add(betAmount)
			updatedFullfillmentQueue = updatedFullfillmentQueue[1:]

			remainingPayoutProfit = remainingPayoutProfit.Sub(availableLiquidty)
			maxLoss := participantExposure.Exposure.Sub(participant.CurrentRoundTotalBetAmount).Add(participantExposure.BetAmount)
			if participant.CurrentRoundMaxLoss.IsNil() {
				participant.CurrentRoundMaxLoss = maxLoss
				participant.CurrentRoundMaxLossOdd = oddId
			} else if participant.CurrentRoundMaxLossOdd == oddId {
				participant.CurrentRoundMaxLoss = maxLoss
			} else {
				originalMaxLoss := participant.CurrentRoundMaxLoss.Sub(betAmount)
				if maxLoss.GT(originalMaxLoss) {
					participant.CurrentRoundMaxLoss = maxLoss
					participant.CurrentRoundMaxLossOdd = oddId
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
			participantBetPair := types.NewParticipantBetPair(participant.BookId, uniqueLock, participant.ParticipantNumber, betId)
			k.SetParticipantBetPair(ctx, participantBetPair)
		} else {
			participantExposure.Exposure = participantExposure.Exposure.Add(remainingPayoutProfit)
			betAmount, err := bettypes.CalculateBetAmount(oddsType, oddsVal, remainingPayoutProfit)
			if err != nil {
				return err, betFullfillment
			}
			calculatedBetAmount = calculatedBetAmount.Add(betAmount)
			participantExposure.BetAmount = participantExposure.BetAmount.Add(betAmount)
			participant.TotalBetAmount = participant.TotalBetAmount.Add(betAmount)
			participant.CurrentRoundTotalBetAmount = participant.CurrentRoundTotalBetAmount.Add(betAmount)
			maxLoss := participantExposure.Exposure.Sub(participant.CurrentRoundTotalBetAmount).Add(participantExposure.BetAmount)
			if participant.CurrentRoundMaxLoss.IsNil() {
				participant.CurrentRoundMaxLoss = maxLoss
				participant.CurrentRoundMaxLossOdd = oddId
			} else if participant.CurrentRoundMaxLossOdd == oddId {
				participant.CurrentRoundMaxLoss = maxLoss
			} else {
				originalMaxLoss := participant.CurrentRoundMaxLoss.Sub(betAmount)
				if maxLoss.GT(originalMaxLoss) {
					participant.CurrentRoundMaxLoss = maxLoss
					participant.CurrentRoundMaxLossOdd = oddId
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
			participantBetPair := types.NewParticipantBetPair(participant.BookId, uniqueLock, participant.ParticipantNumber, betId)
			k.SetParticipantBetPair(ctx, participantBetPair)
		}
		k.SetParticipantExposure(ctx, participantExposure)
		k.SetBookParticipant(ctx, participant)

		if participant.ExposuresNotFilled == 0 {
			// add back to queue
			maxLoss := sdk.MaxInt(sdk.ZeroInt(), participant.CurrentRoundMaxLoss)
			participant.CurrentRoundLiquidity = participant.CurrentRoundLiquidity.Sub(maxLoss)
			eligibleForNextRound := participant.CurrentRoundLiquidity.GT(sdk.ZeroInt())

			participantExposures := k.GetExposureByBookAndParticipantNumber(ctx, bookId, pn)
			for _, pe := range participantExposures {
				k.MoveToHistoricalParticipantExposure(ctx, pe)
				if eligibleForNextRound {
					newPe := types.NewParticipantExposure(book.Id, pe.OddId, sdk.ZeroInt(), sdk.ZeroInt(), pe.ParticipantNumber, pe.Round+1, false)
					k.SetParticipantExposure(ctx, newPe)
					if pe.OddId == participantExposure.OddId {
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
				boes := k.GetOddExposuresByBook(ctx, bookId)
				for _, boe := range boes {
					boe.FullfillmentQueue = append(boe.FullfillmentQueue, pn)
					if boe.OddId == participantExposure.OddId {
						bookExposure = boe
					}

					k.SetBookOddExposure(ctx, boe)
				}
				updatedFullfillmentQueue = append(updatedFullfillmentQueue, pn)
				bookExposure.FullfillmentQueue = updatedFullfillmentQueue
			}
		}

		if remainingPayoutProfit.LTE(sdk.ZeroInt()) || len(updatedFullfillmentQueue) <= 0 {
			break
		}
	}

	if remainingPayoutProfit.GT(sdk.ZeroInt()) {
		return sdkerrors.Wrapf(types.ErrInternalProcessingBet, "insufficient liquidity in order book"), betFullfillment
	}

	bookExposure.FullfillmentQueue = updatedFullfillmentQueue
	k.SetBookOddExposure(ctx, bookExposure)

	// Transfer bet fee from bettor to the `bet` module account
	err := k.transferFundsFromUserToModule(ctx, bettorAddress, bettypes.ModuleName, betFee)
	if err != nil {
		return err, betFullfillment
	}

	// Transfer bet amount from bettor to `book_liquidity_pool` Account
	err = k.transferFundsFromUserToModule(ctx, bettorAddress, types.BookLiquidityName, calculatedBetAmount)
	if err != nil {
		return err, betFullfillment
	}

	// Create a unique lock in the Payout Store for the bet
	k.setPayoutLock(ctx, uniqueLock)

	return nil, betFullfillment
}
