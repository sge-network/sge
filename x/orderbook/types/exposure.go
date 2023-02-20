package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewBookOddsExposure creates a new book odds exposure object
//
//nolint:interfacer
func NewBookOddsExposure(bookID, oddsID string, fulfillmentQueue []uint64) BookOddsExposure {
	return BookOddsExposure{
		BookID:           bookID,
		OddsID:           oddsID,
		FulfillmentQueue: fulfillmentQueue,
	}
}

// String returns a human readable string representation of a BookOddsExposure.
func (boe BookOddsExposure) String() string {
	out, _ := yaml.Marshal(boe)
	return string(out)
}

// NewParticipationExposure creates a new participation exposure object
//
//nolint:interfacer
func NewParticipationExposure(bookID, oddsID string, exposure, betAmount sdk.Int, participationIndex, round uint64, isFulfilled bool) ParticipationExposure {
	return ParticipationExposure{
		BookID:             bookID,
		OddsID:             oddsID,
		ParticipationIndex: participationIndex,
		Exposure:           exposure,
		BetAmount:          betAmount,
		IsFulfilled:        isFulfilled,
		Round:              round,
	}
}

// String returns a human readable string representation of a participationExposure.
func (pe ParticipationExposure) String() string {
	out, _ := yaml.Marshal(pe)
	return string(out)
}
