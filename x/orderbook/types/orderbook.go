package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	yaml "gopkg.in/yaml.v2"
)

// NewBook creates a new orderbook object
//
//nolint:interfacer
func NewBook(bookID string, participationCount, oddsCount uint64, status OrderBookStatus) OrderBook {
	return OrderBook{
		ID:                 bookID,
		ParticipationCount: participationCount,
		Status:             status,
		OddsCount:          oddsCount,
	}
}

// MustMarshalBook returns the orderbook bytes. Panics if fails
func MustMarshalBook(cdc codec.BinaryCodec, book OrderBook) []byte {
	return cdc.MustMarshal(&book)
}

// MustUnmarshalBook return the unmarshaled orderbook from bytes.
// Panics if fails.
func MustUnmarshalBook(cdc codec.BinaryCodec, value []byte) OrderBook {
	book, err := UnmarshalBook(cdc, value)
	if err != nil {
		panic(err)
	}

	return book
}

// return the orderbook
func UnmarshalBook(cdc codec.BinaryCodec, value []byte) (book OrderBook, err error) {
	err = cdc.Unmarshal(value, &book)
	return book, err
}

// String returns a human readable string representation of a OrderBook.
func (ob OrderBook) String() string {
	out, _ := yaml.Marshal(ob)
	return string(out)
}
