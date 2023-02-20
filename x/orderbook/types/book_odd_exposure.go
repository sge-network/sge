package types

import (
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
