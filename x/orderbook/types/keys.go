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
func GetBookKey(bookUID string) []byte {
	return utils.StrBytes(bookUID)
}

// GetBookParticipationKey creates the key for book participation bond with book
func GetBookParticipationKey(bookUID string, participationIndex uint64) []byte {
	return append(GetBookParticipationsKey(bookUID), utils.Uint64ToBytes(participationIndex)...)
}

// GetBookParticipationsKey creates the key for book participations for an book
func GetBookParticipationsKey(bookUID string) []byte {
	return utils.StrBytes(bookUID)
}

// GetBookOddsExposureKey creates the key for book exposure for an odd
func GetBookOddsExposureKey(bookUID, oddsUID string) []byte {
	return append(GetBookOddsExposuresKey(bookUID), utils.StrBytes(oddsUID)...)
}

// GetBookOddsExposuresKey creates the key for book exposure for an book
func GetBookOddsExposuresKey(bookUID string) []byte {
	return utils.StrBytes(bookUID)
}

// GetParticipationExposureKey creates the key for participation exposure for an odd
func GetParticipationExposureKey(bookUID, oddsUID string, index uint64) []byte {
	return append(GetParticipationExposuresKey(bookUID, oddsUID), utils.Uint64ToBytes(index)...)
}

// GetParticipationExposuresKey creates the key for exposures for a book id and odds id
func GetParticipationExposuresKey(bookUID, oddsUID string) []byte {
	return append(GetParticipationExposuresByBookKey(bookUID), utils.StrBytes(oddsUID)...)
}

// GetParticipationExposuresByBookKey creates the key for exposures for a book id
func GetParticipationExposuresByBookKey(bookUID string) []byte {
	return utils.StrBytes(bookUID)
}

// GetParticipationExposureByIndexKey creates the key for participation exposure for an odds by participation number
func GetParticipationExposureByIndexKey(bookUID, oddsUID string, index uint64) []byte {
	return append(GetParticipationByIndexKey(bookUID, index), utils.StrBytes(oddsUID)...)
}

// GetParticipationByIndexKey creates the key for exposures for a book id and participation number
func GetParticipationByIndexKey(bookUID string, index uint64) []byte {
	return append(utils.StrBytes(bookUID), utils.Uint64ToBytes(index)...)
}

// GetHistoricalParticipationExposureKey creates the key for participation exposure for an odd
func GetHistoricalParticipationExposureKey(bookUID, oddsUID string, index, round uint64) []byte {
	return append(GetParticipationExposureKey(bookUID, oddsUID, index), utils.Uint64ToBytes(round)...)
}

// GetParticipationBetPairKey creates the bond between participation and bet
func GetParticipationBetPairKey(bookUID string, bookParticipationIndex uint64, betID uint64) []byte {
	return append(GetParticipationByIndexKey(bookUID, bookParticipationIndex), utils.Uint64ToBytes(betID)...)
}
