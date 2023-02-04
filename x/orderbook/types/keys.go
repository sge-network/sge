package types

import (
	"github.com/sge-network/sge/utils"
)

const (
	// ModuleName is the name of the orderbook module
	ModuleName = "orderbook"

	// StoreKey is the string store representation
	StoreKey = ModuleName
)

// module accounts constants
const (
	// BookLiquidityName defines the account name for book liquidity for users
	BookLiquidityName = "book_liquidity_pool"

	// SRBookLiquidityName defines the account name for book liquidity for sr
	SRBookLiquidityName = "sr_book_liquidity_pool"
)

var (
	BookKeyPrefix            = []byte{0x00} // prefix for keys that store books
	BookParticiapntKeyPrefix = []byte{0x01} // prefix for keys that store book particiapnts
)

// GetBookKey returns the bytes of an book key
func GetBookKey(bookID string) []byte {
	return utils.StrBytes(bookID)
}

// GetBookParticipantKey creates the key for book participant bond with book
func GetBookParticipantKey(bookID string, bookParticipantNumber uint64) []byte {
	return append(utils.StrBytes(bookID), utils.Uint64ToBytes(bookParticipantNumber)...)
}
