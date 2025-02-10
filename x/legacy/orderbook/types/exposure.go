package types

import (
	yaml "gopkg.in/yaml.v2"

	sdkmath "cosmossdk.io/math"
)

// NewOrderBookOddsExposure creates a new book odds exposure object
//
//nolint:interface
func NewOrderBookOddsExposure(
	orderBookUID, oddsUID string,
	fulfillmentQueue []uint64,
) OrderBookOddsExposure {
	return OrderBookOddsExposure{
		OrderBookUID:     orderBookUID,
		OddsUID:          oddsUID,
		FulfillmentQueue: fulfillmentQueue,
	}
}

// String returns a human-readable string representation of a BookOddsExposure.
func (boe OrderBookOddsExposure) String() string {
	out, err := yaml.Marshal(boe)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// NewParticipationExposure creates a new participation exposure object
//
//nolint:interface
func NewParticipationExposure(
	orderBookUID, oddsUID string,
	exposure, betAmount sdkmath.Int,
	participationIndex, round uint64,
	isFulfilled bool,
) ParticipationExposure {
	return ParticipationExposure{
		OrderBookUID:       orderBookUID,
		OddsUID:            oddsUID,
		ParticipationIndex: participationIndex,
		Exposure:           exposure,
		BetAmount:          betAmount,
		IsFulfilled:        isFulfilled,
		Round:              round,
	}
}

// NextRound returns the next round participation object extracted from the current round properties.
func (pe ParticipationExposure) NextRound() ParticipationExposure {
	return NewParticipationExposure(
		pe.OrderBookUID,
		pe.OddsUID,
		sdkmath.ZeroInt(),
		sdkmath.ZeroInt(),
		pe.ParticipationIndex,
		pe.Round+1,
		false,
	)
}

// String returns a human-readable string representation of a participationExposure.
func (pe ParticipationExposure) String() string {
	out, err := yaml.Marshal(pe)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// CalculateMaxLoss calculates the maximum amount of loss for an exposure
// according to the bet amount.
func (pe ParticipationExposure) CalculateMaxLoss(totalBetAmount sdkmath.Int) sdkmath.Int {
	return pe.Exposure.Add(pe.BetAmount).Sub(totalBetAmount)
}

// SetCurrentRound sets the current round bet amount and payout profit.
func (pe *ParticipationExposure) SetCurrentRound(betAmount, payoutProfit sdkmath.Int) {
	// add the payout profit to the
	pe.Exposure = pe.Exposure.Add(payoutProfit)

	// add the bet amount that is being fulfilled to the exposure and participation
	pe.BetAmount = pe.BetAmount.Add(betAmount)
}
