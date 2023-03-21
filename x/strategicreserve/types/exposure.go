package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewBookOddsExposure creates a new book odds exposure object
//
//nolint:interfacer
func NewBookOddsExposure(bookUID, oddsUID string, fulfillmentQueue []uint64) BookOddsExposure {
	return BookOddsExposure{
		BookUID:          bookUID,
		OddsUID:          oddsUID,
		FulfillmentQueue: fulfillmentQueue,
	}
}

// String returns a human readable string representation of a BookOddsExposure.
func (boe BookOddsExposure) String() string {
	out, err := yaml.Marshal(boe)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// NewParticipationExposure creates a new participation exposure object
//
//nolint:interfacer
func NewParticipationExposure(bookUID, oddsUID string, exposure, betAmount sdk.Int, participationIndex, round uint64, isFulfilled bool) ParticipationExposure {
	return ParticipationExposure{
		BookUID:            bookUID,
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
