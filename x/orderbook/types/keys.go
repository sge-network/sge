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
	BookOddsExposureKeyPrefix              = []byte{0x02} // prefix for keys that store book odds exposures
	ParticipantExposureKeyPrefix           = []byte{0x03} // prefix for keys that store participant exposures
	PayoutLockKeyPrefix                    = []byte{0x04} // prefix for keys that store payout locks
	ParticipantExposureByPNKeyPrefix       = []byte{0x05} // prefix for keys that store participant exposures
	HistoricalParticipantExposureKeyPrefix = []byte{0x06} // prefix for keys that store historical participant exposures
	BookStatsKeyPrefix                     = []byte{0x07} // prefix for keys that store book stats
	ParticipantBetPairKeyPrefix            = []byte{0x08} // prefix for keys that store book participant and bet pairs
)

// GetBookKey returns the bytes of an book key
func GetBookKey(bookID string) []byte {
	return utils.StrBytes(bookID)
}

// GetBookParticipantKey creates the key for book participant bond with book
func GetBookParticipantKey(bookID string, bookParticipantNumber uint64) []byte {
	return append(GetBookParticipantsKey(bookID), utils.Uint64ToBytes(bookParticipantNumber)...)
}

// GetBookParticipantsKey creates the key for book participants for an book
func GetBookParticipantsKey(bookID string) []byte {
	return utils.StrBytes(bookID)
}

// GetBookOddsExposureKey creates the key for book exposure for an odd
func GetBookOddsExposureKey(bookID, oddsID string) []byte {
	return append(GetBookOddsExposuresKey(bookID), utils.StrBytes(oddsID)...)
}

// GetBookOddsExposuresKey creates the key for book exposure for an book
func GetBookOddsExposuresKey(bookID string) []byte {
	return utils.StrBytes(bookID)
}

// GetParticipantExposureKey creates the key for participant exposure for an odd
func GetParticipantExposureKey(bookID, oddsID string, pn uint64) []byte {
	return append(GetParticipantExposuresKey(bookID, oddsID), utils.Uint64ToBytes(pn)...)
}

// GetParticipantExposuresKey creates the key for exposures for a book id and odds id
func GetParticipantExposuresKey(bookID, oddsID string) []byte {
	return append(GetParticipantExposuresByBookKey(bookID), utils.StrBytes(oddsID)...)
}

// GetParticipantExposuresByBookKey creates the key for exposures for a book id
func GetParticipantExposuresByBookKey(bookID string) []byte {
	return utils.StrBytes(bookID)
}

// GetParticipantExposureByPNKey creates the key for participant exposure for an odds by participant number
func GetParticipantExposureByPNKey(bookID, oddsID string, pn uint64) []byte {
	return append(GetParticipantByPNKey(bookID, pn), utils.StrBytes(oddsID)...)
}

// GetParticipantExposuresByPNKey creates the key for exposures for a book id and participant number
func GetParticipantByPNKey(bookID string, pn uint64) []byte {
	return append(utils.StrBytes(bookID), utils.Uint64ToBytes(pn)...)
}

// GetHistoricalParticipantExposureKey creates the key for participant exposure for an odd
func GetHistoricalParticipantExposureKey(bookID, oddsID string, pn, round uint64) []byte {
	return append(GetParticipantExposureKey(bookID, oddsID, pn), utils.Uint64ToBytes(round)...)
}

// GetParticipantBetPairKey creates the bond between participant and bet
func GetParticipantBetPairKey(bookID string, bookParticipantNumber uint64, betID uint64) []byte {
	return append(GetParticipantByPNKey(bookID, bookParticipantNumber), utils.Uint64ToBytes(betID)...)
}
