package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	yaml "gopkg.in/yaml.v2"
)

// NewOrderBook creates a new orderbook object
//
//nolint:interfacer
func NewOrderBook(bookUID string, participationCount, oddsCount uint64, status OrderBookStatus) OrderBook {
	return OrderBook{
		UID:                bookUID,
		ParticipationCount: participationCount,
		Status:             status,
		OddsCount:          oddsCount,
	}
}

// return the orderbook
func UnmarshalOrderBook(cdc codec.BinaryCodec, value []byte) (book OrderBook, err error) {
	err = cdc.Unmarshal(value, &book)
	return book, err
}

// String returns a human readable string representation of a OrderBook.
func (ob OrderBook) String() string {
	out, err := yaml.Marshal(ob)
	if err != nil {
		panic(err)
	}
	return string(out)
}
