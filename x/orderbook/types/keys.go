package types

import (
	"github.com/sge-network/sge/utils"
)

const (
	// ModuleName is the name of the orderbook module
	ModuleName = "orderbook"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the message route for orderbook
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// module accounts constants
const (
	// orderBookLiquidityPool defines the account name for book liquidity for the participants.
	orderBookLiquidityPool = "orderbook_liquidity_pool"
)

var (
	OrderBookKeyPrefix              = []byte{0x00} // prefix for keys that store books
	OrderBookParticipationKeyPrefix = []byte{
		0x01,
	} // prefix for keys that store book participations
	OrderBookOddsExposureKeyPrefix = []byte{
		0x02,
	} // prefix for keys that store book odds exposures
	ParticipationExposureKeyPrefix = []byte{
		0x03,
	} // prefix for keys that store participation exposures
	ParticipationExposureByIndexKeyPrefix = []byte{
		0x04,
	} // prefix for keys that store participation exposures
	HistoricalParticipationExposureKeyPrefix = []byte{
		0x05,
	} // prefix for keys that store historical participation exposures
	OrderBookStatsKeyPrefix       = []byte{0x06} // prefix for keys that store book stats
	ParticipationBetPairKeyPrefix = []byte{
		0x07,
	} // prefix for keys that store book participation and bet pairs
	FeeGrantPrefix = []byte{0x08} // prefix for keys that store fee grants
)

// GetOrderBookKey returns the bytes of an book key
func GetOrderBookKey(bookUID string) []byte {
	return utils.StrBytes(bookUID)
}

// GetOrderBookParticipationKey creates the key for book participation bond with book
func GetOrderBookParticipationKey(bookUID string, participationIndex uint64) []byte {
	return append(GetOrderBookParticipationsKey(bookUID), utils.Uint64ToBytes(participationIndex)...)
}

// GetOrderBookParticipationsKey creates the key for book participations for an book
func GetOrderBookParticipationsKey(bookUID string) []byte {
	return utils.StrBytes(bookUID)
}

// GetOrderBookOddsExposureKey creates the key for book exposure for an odd
func GetOrderBookOddsExposureKey(bookUID, oddsUID string) []byte {
	return append(GetOrderBookOddsExposuresKey(bookUID), utils.StrBytes(oddsUID)...)
}

// GetOrderBookOddsExposuresKey creates the key for book exposure for an book
func GetOrderBookOddsExposuresKey(bookUID string) []byte {
	return utils.StrBytes(bookUID)
}

// GetParticipationExposureKey creates the key for participation exposure for an odd
func GetParticipationExposureKey(bookUID, oddsUID string, index uint64) []byte {
	return append(GetParticipationExposuresKey(bookUID, oddsUID), utils.Uint64ToBytes(index)...)
}

// GetParticipationExposuresKey creates the key for exposures for a book id and odds id
func GetParticipationExposuresKey(bookUID, oddsUID string) []byte {
	return append(GetParticipationExposuresByOrderBookKey(bookUID), utils.StrBytes(oddsUID)...)
}

// GetParticipationExposuresByOrderBookKey creates the key for exposures for a book id
func GetParticipationExposuresByOrderBookKey(bookUID string) []byte {
	return utils.StrBytes(bookUID)
}

// GetParticipationExposureByIndexKey creates the key for participation exposure for an odds by participation index
func GetParticipationExposureByIndexKey(bookUID, oddsUID string, index uint64) []byte {
	return append(GetParticipationByIndexKey(bookUID, index), utils.StrBytes(oddsUID)...)
}

// GetParticipationByIndexKey creates the key for exposures for a book id and participation index
func GetParticipationByIndexKey(bookUID string, index uint64) []byte {
	return append(utils.StrBytes(bookUID), utils.Uint64ToBytes(index)...)
}

// GetHistoricalParticipationExposureKey creates the key for participation exposure for an odd
func GetHistoricalParticipationExposureKey(bookUID, oddsUID string, index, round uint64) []byte {
	return append(GetParticipationExposureKey(bookUID, oddsUID, index), utils.Uint64ToBytes(round)...)
}

// GetParticipationBetPairKey creates the bond between participation and bet
func GetParticipationBetPairKey(bookUID string, bookParticipationIndex uint64, betID uint64) []byte {
	return append(
		GetParticipationByIndexKey(bookUID, bookParticipationIndex),
		utils.Uint64ToBytes(betID)...)
}
