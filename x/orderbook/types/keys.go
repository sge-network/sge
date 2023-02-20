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
	BookKeyPrefix                            = []byte{0x00} // prefix for keys that store books
	BookParticipationKeyPrefix               = []byte{0x01} // prefix for keys that store book participations
	BookOddsExposureKeyPrefix                = []byte{0x02} // prefix for keys that store book odds exposures
	ParticipationExposureKeyPrefix           = []byte{0x03} // prefix for keys that store participation exposures
	PayoutLockKeyPrefix                      = []byte{0x04} // prefix for keys that store payout locks
	ParticipationExposureByIndexKeyPrefix    = []byte{0x05} // prefix for keys that store participation exposures
	HistoricalParticipationExposureKeyPrefix = []byte{0x06} // prefix for keys that store historical participation exposures
	BookStatsKeyPrefix                       = []byte{0x07} // prefix for keys that store book stats
	ParticipationBetPairKeyPrefix            = []byte{0x08} // prefix for keys that store book participation and bet pairs
)

// GetBookKey returns the bytes of an book key
func GetBookKey(bookID string) []byte {
	return utils.StrBytes(bookID)
}

// GetBookParticipationKey creates the key for book participation bond with book
func GetBookParticipationKey(bookID string, participationIndex uint64) []byte {
	return append(GetBookParticipationsKey(bookID), utils.Uint64ToBytes(participationIndex)...)
}

// GetBookParticipationsKey creates the key for book participations for an book
func GetBookParticipationsKey(bookID string) []byte {
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

// GetParticipationExposureKey creates the key for participation exposure for an odd
func GetParticipationExposureKey(bookID, oddsID string, index uint64) []byte {
	return append(GetParticipationExposuresKey(bookID, oddsID), utils.Uint64ToBytes(index)...)
}

// GetParticipationExposuresKey creates the key for exposures for a book id and odds id
func GetParticipationExposuresKey(bookID, oddsID string) []byte {
	return append(GetParticipationExposuresByBookKey(bookID), utils.StrBytes(oddsID)...)
}

// GetParticipationExposuresByBookKey creates the key for exposures for a book id
func GetParticipationExposuresByBookKey(bookID string) []byte {
	return utils.StrBytes(bookID)
}

// GetParticipationExposureByIndexKey creates the key for participation exposure for an odds by participation number
func GetParticipationExposureByIndexKey(bookID, oddsID string, index uint64) []byte {
	return append(GetParticipationByIndexKey(bookID, index), utils.StrBytes(oddsID)...)
}

// GetParticipationByIndexKey creates the key for exposures for a book id and participation number
func GetParticipationByIndexKey(bookID string, index uint64) []byte {
	return append(utils.StrBytes(bookID), utils.Uint64ToBytes(index)...)
}

// GetHistoricalParticipationExposureKey creates the key for participation exposure for an odd
func GetHistoricalParticipationExposureKey(bookID, oddsID string, index, round uint64) []byte {
	return append(GetParticipationExposureKey(bookID, oddsID, index), utils.Uint64ToBytes(round)...)
}

// GetParticipationBetPairKey creates the bond between participation and bet
func GetParticipationBetPairKey(bookID string, bookParticipationNumber uint64, betID uint64) []byte {
	return append(GetParticipationByIndexKey(bookID, bookParticipationNumber), utils.Uint64ToBytes(betID)...)
}
