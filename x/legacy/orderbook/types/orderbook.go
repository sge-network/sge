package types

import (
	yaml "gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/codec"
)

// NewOrderBook creates a new orderbook object
//
//nolint:interface
func NewOrderBook(
	bookUID string, oddsCount uint64,
	status OrderBookStatus,
) OrderBook {
	return OrderBook{
		UID:                bookUID,
		ParticipationCount: 0,
		Status:             status,
		OddsCount:          oddsCount,
	}
}

// UnmarshalOrderBook return the orderbook
func UnmarshalOrderBook(cdc codec.BinaryCodec, value []byte) (book OrderBook, err error) {
	err = cdc.Unmarshal(value, &book)
	return book, err
}

// String returns a human-readable string representation of a OrderBook.
func (ob OrderBook) String() string {
	out, err := yaml.Marshal(ob)
	if err != nil {
		panic(err)
	}
	return string(out)
}
