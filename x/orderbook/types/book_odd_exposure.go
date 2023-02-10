package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	yaml "gopkg.in/yaml.v2"
)

// NewBookOddExposure creates a new book odd exposure object
//nolint:interfacer
func NewBookOddExposure(bookId, oddId string, fullfillmentQueue []uint64) BookOddExposure {
	return BookOddExposure{
		BookId:            bookId,
		OddId:             oddId,
		FullfillmentQueue: fullfillmentQueue,
	}
}

// MustMarshalBookOddExposure returns the book exposure bytes. Panics if fails
func MustMarshalBookOddExposure(cdc codec.BinaryCodec, boe BookOddExposure) []byte {
	return cdc.MustMarshal(&boe)
}

// MustUnmarshalBookOddExposure return the unmarshaled book odd exposure from bytes.
// Panics if fails.
func MustUnmarshalBookOddExposure(cdc codec.BinaryCodec, value []byte) BookOddExposure {
	boe, err := UnmarshalBookOddExposure(cdc, value)
	if err != nil {
		panic(err)
	}

	return boe
}

// return the book odd exposure
func UnmarshalBookOddExposure(cdc codec.BinaryCodec, value []byte) (boe BookOddExposure, err error) {
	err = cdc.Unmarshal(value, &boe)
	return boe, err
}

// String returns a human readable string representation of a BookOddExposure.
func (boe BookOddExposure) String() string {
	out, _ := yaml.Marshal(boe)
	return string(out)
}
