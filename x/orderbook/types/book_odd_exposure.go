package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	yaml "gopkg.in/yaml.v2"
)

// NewBookOddsExposure creates a new book odds exposure object
//
//nolint:interfacer
func NewBookOddsExposure(bookID, oddsID string, fullfillmentQueue []uint64) BookOddsExposure {
	return BookOddsExposure{
		BookID:            bookID,
		OddsID:            oddsID,
		FullfillmentQueue: fullfillmentQueue,
	}
}

// MustMarshalBookOddsExposure returns the book exposure bytes. Panics if fails
func MustMarshalBookOddsExposure(cdc codec.BinaryCodec, boe BookOddsExposure) []byte {
	return cdc.MustMarshal(&boe)
}

// MustUnmarshalBookOddsExposure return the unmarshaled book odds exposure from bytes.
// Panics if fails.
func MustUnmarshalBookOddsExposure(cdc codec.BinaryCodec, value []byte) BookOddsExposure {
	boe, err := UnmarshalBookOddsExposure(cdc, value)
	if err != nil {
		panic(err)
	}

	return boe
}

// return the book odds exposure
func UnmarshalBookOddsExposure(cdc codec.BinaryCodec, value []byte) (boe BookOddsExposure, err error) {
	err = cdc.Unmarshal(value, &boe)
	return boe, err
}

// String returns a human readable string representation of a BookOddsExposure.
func (boe BookOddsExposure) String() string {
	out, _ := yaml.Marshal(boe)
	return string(out)
}
