package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewOrderBookParticipation creates a new book participation object
//
//nolint:interfacer
func NewOrderBookParticipation(
	index uint64, orderBookUID string, participantAddress string,
	exposuresNotFilled uint64,
	liquidity, currentRoundLiquidity, totalBetAmount, currentRoundTotalBetAmount, maxLoss, currentRoundMaxLoss sdk.Int,
	currentRoundMaxLossOddsUID string, actualProfit sdk.Int,
) OrderBookParticipation {
	return OrderBookParticipation{
		Index:                      index,
		OrderBookUID:               orderBookUID,
		ParticipantAddress:         participantAddress,
		Liquidity:                  liquidity,
		CurrentRoundLiquidity:      currentRoundLiquidity,
		ExposuresNotFilled:         exposuresNotFilled,
		TotalBetAmount:             totalBetAmount,
		CurrentRoundTotalBetAmount: currentRoundTotalBetAmount,
		MaxLoss:                    maxLoss,
		CurrentRoundMaxLoss:        currentRoundMaxLoss,
		CurrentRoundMaxLossOddsUID: currentRoundMaxLossOddsUID,
		ActualProfit:               actualProfit,
	}
}

// String returns a human readable string representation of a BookParticipation.
func (p OrderBookParticipation) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// CalculateMaxLoss calculates the maxixmum amount of the tokens expected to be the
// loss of the participation according to the bet amount
func (p OrderBookParticipation) calculateMaxLoss(betAmount sdk.Int) sdk.Int {
	return p.CurrentRoundMaxLoss.Sub(betAmount)
}

// IsEligibleForNextRound determines if the participation has enough
// liquidity to be used in the next round or not
func (p *OrderBookParticipation) IsEligibleForNextRound() bool {
	return p.CurrentRoundLiquidity.GT(sdk.ZeroInt())
}

// TrimCurrentRoundLiquidity subtracts the max loss from the current round liquidity.
func (p *OrderBookParticipation) TrimCurrentRoundLiquidity() {
	maxLoss := sdk.MaxInt(sdk.ZeroInt(), p.CurrentRoundMaxLoss)
	p.CurrentRoundLiquidity = p.CurrentRoundLiquidity.Sub(maxLoss)
}

// ResetForNextRound resets the exposures, max loss and current round amount
// and make the participation ready for the next round
func (p *OrderBookParticipation) ResetForNextRound(notFilledExposures uint64) {
	// prepare participation for the next round
	p.ExposuresNotFilled = notFilledExposures
	p.MaxLoss = p.MaxLoss.Add(p.CurrentRoundMaxLoss)
	p.CurrentRoundTotalBetAmount = sdk.ZeroInt()
	p.CurrentRoundMaxLoss = sdk.ZeroInt()
}

// SetCurrentRound sets the current round total bet amount and max loss.
func (p *OrderBookParticipation) SetCurrentRound(pe *ParticipationExposure, oddsUID string, betAmount sdk.Int) {
	p.TotalBetAmount = p.TotalBetAmount.Add(betAmount)
	p.CurrentRoundTotalBetAmount = p.CurrentRoundTotalBetAmount.Add(betAmount)
	p.setMaxLoss(pe, oddsUID, betAmount)
}

// setMaxLoss sets the max loss of the the cirrent round.
func (p *OrderBookParticipation) setMaxLoss(pe *ParticipationExposure, oddsUID string, betAmount sdk.Int) {
	// max loss is the maximum amount that an exposure may lose.
	maxLoss := pe.calculateMaxLoss(p.CurrentRoundTotalBetAmount)
	switch {
	case p.CurrentRoundMaxLoss.IsNil():
		p.CurrentRoundMaxLoss = maxLoss
		p.CurrentRoundMaxLossOddsUID = oddsUID
	case p.CurrentRoundMaxLossOddsUID == oddsUID:
		p.CurrentRoundMaxLoss = maxLoss
	default:
		originalMaxLoss := p.calculateMaxLoss(betAmount)
		if maxLoss.GT(originalMaxLoss) {
			p.CurrentRoundMaxLoss = maxLoss
			p.CurrentRoundMaxLossOddsUID = oddsUID
		} else {
			p.CurrentRoundMaxLoss = originalMaxLoss
		}
	}
}
