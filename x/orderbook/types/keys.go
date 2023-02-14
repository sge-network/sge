package types

import (
	"github.com/sge-network/sge/utils"
)

const (
	// ModuleName is the name of the orderbook module
	ModuleName = "orderbook"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// module accounts constants
const (
	// BookLiquidityName defines the account name for book liquidity for users
	BookLiquidityName = "book_liquidity_pool"

	// SRProfitName defines the account name for profit for sr
	SRProfitName = "sr_profit_pool"
)

var (
	BookKeyPrefix                          = []byte{0x00} // prefix for keys that store books
	BookParticipantKeyPrefix               = []byte{0x01} // prefix for keys that store book participants
	BookOddExposureKeyPrefix               = []byte{0x02} // prefix for keys that store book odd exposures
	ParticipantExposureKeyPrefix           = []byte{0x03} // prefix for keys that store participant exposures
	PayoutLockKeyPrefix                    = []byte{0x04} // prefix for keys that store payout locks
	ParticipantExposureByPNKeyPrefix       = []byte{0x05} // prefix for keys that store participant exposures
	HistoricalParticipantExposureKeyPrefix = []byte{0x06} // prefix for keys that store historical participant exposures
	BookStatsKeyPrefix                     = []byte{0x07} // prefix for keys that store book stats
	ParticipantBetPairKeyPrefix            = []byte{0x08} // prefix for keys that store book participant and bet pairs
)

// GetBookKey returns the bytes of an book key
func GetBookKey(bookId string) []byte {
	return utils.StrBytes(bookId)
}

// GetBookParticipantKey creates the key for book participant bond with book
func GetBookParticipantKey(bookId string, bookParticipantNumber uint64) []byte {
	return append(GetBookParticipantsKey(bookId), utils.Uint64ToBytes(bookParticipantNumber)...)
}

// GetBookParticipantsKey creates the key for book participants for an book
func GetBookParticipantsKey(bookId string) []byte {
	return utils.StrBytes(bookId)
}

// GetBookOddExposureKey creates the key for book exposure for an odd
func GetBookOddExposureKey(bookId, oddId string) []byte {
	return append(GetBookOddExposuresKey(bookId), utils.StrBytes(oddId)...)
}

// GetBookOddExposuresKey creates the key for book exposure for an book
func GetBookOddExposuresKey(bookId string) []byte {
	return utils.StrBytes(bookId)
}

// GetParticipantExposureKey creates the key for participant exposure for an odd
func GetParticipantExposureKey(bookId, oddId string, pn uint64) []byte {
	return append(GetParticipantExposuresKey(bookId, oddId), utils.Uint64ToBytes(pn)...)
}

// GetParticipantExposuresKey creates the key for exposures for a book id and odd id
func GetParticipantExposuresKey(bookId, oddId string) []byte {
	return append(GetParticipantExposuresByBookKey(bookId), utils.StrBytes(oddId)...)
}

// GetParticipantExposuresByBookKey creates the key for exposures for a book id
func GetParticipantExposuresByBookKey(bookId string) []byte {
	return utils.StrBytes(bookId)
}

// GetParticipantExposureByPNKey creates the key for participant exposure for an odd by participant number
func GetParticipantExposureByPNKey(bookId, oddId string, pn uint64) []byte {
	return append(GetParticipantByPNKey(bookId, pn), utils.StrBytes(oddId)...)
}

// GetParticipantExposuresByPNKey creates the key for exposures for a book id and participant number
func GetParticipantByPNKey(bookId string, pn uint64) []byte {
	return append(utils.StrBytes(bookId), utils.Uint64ToBytes(pn)...)
}

// GetHistoricalParticipantExposureKey creates the key for participant exposure for an odd
func GetHistoricalParticipantExposureKey(bookId, oddId string, pn, round uint64) []byte {
	return append(GetParticipantExposureKey(bookId, oddId, pn), utils.Uint64ToBytes(round)...)
}

// GetParticipantBetPairKey creates the bond between participant and bet
func GetParticipantBetPairKey(bookId string, bookParticipantNumber uint64, betId uint64) []byte {
	return append(GetParticipantByPNKey(bookId, bookParticipantNumber), utils.Uint64ToBytes(betId)...)
}
