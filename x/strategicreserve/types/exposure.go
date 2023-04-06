package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewOrderBookOddsExposure creates a new book odds exposure object
//
//nolint:interfacer
func NewOrderBookOddsExposure(orderBookUID, oddsUID string, fulfillmentQueue []uint64) OrderBookOddsExposure {
	return OrderBookOddsExposure{
		OrderBookUID:     orderBookUID,
		OddsUID:          oddsUID,
		FulfillmentQueue: fulfillmentQueue,
	}
}

// String returns a human readable string representation of a BookOddsExposure.
func (boe OrderBookOddsExposure) String() string {
	out, err := yaml.Marshal(boe)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// NewParticipationExposure creates a new participation exposure object
//
//nolint:interfacer
func NewParticipationExposure(orderBookUID, oddsUID string, exposure, betAmount sdk.Int, participationIndex, round uint64, isFulfilled bool) ParticipationExposure {
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

// String returns a human readable string representation of a participationExposure.
func (pe ParticipationExposure) String() string {
	out, err := yaml.Marshal(pe)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// CalculateMaxLoss calculates the maximum amount of loss for an exposure
// according to the bet amount.
func (pe ParticipationExposure) CalculateMaxLoss(totalBetAmount sdk.Int) sdk.Int {
	return pe.Exposure.
		Sub(totalBetAmount).
		Add(pe.BetAmount)
}
