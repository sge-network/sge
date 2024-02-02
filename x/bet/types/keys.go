package types

import (
	"encoding/binary"

	"github.com/sge-network/sge/utils"
)

var _ binary.ByteOrder

const (
	// ModuleName defines the module name
	ModuleName = "bet"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_bet"
)

const (
	// betFeeCollector is the module account name for the bet fee collection module account.
	betFeeCollector = "bet_fee_collector"
	// priceLockPool defines the account name for price locking feature pool.
	priceLockPool = "price_lock_pool"
)

var (
	// BetListPrefix is the prefix to retrieve all Bet
	BetListPrefix = []byte{0x00}
	// BetIDListPrefix is the prefix to retrieve all Bet IDs
	BetIDListPrefix = []byte{0x01}
	// BetStatsKey is the key for the bet statistics
	BetStatsKey = []byte{0x02}
	// PendingBetListPrefix is the prefix to retrieve all pending bets
	PendingBetListPrefix = []byte{0x03}
	// SettledBetListPrefix is the prefix to retrieve all settled bets
	SettledBetListPrefix = []byte{0x04}
)

// BetListByCreatorPrefix returns prefix of the certain creator bet list.
func BetListByCreatorPrefix(creator string) []byte {
	return append(BetListPrefix, utils.StrBytes(creator)...)
}

// BetIDKey returns key of a certain bet of a creator.
func BetIDKey(creator string, id uint64) []byte {
	return append(utils.StrBytes(creator), utils.Uint64ToBytes(id)...)
}

// PendingBetListOfMarketPrefix returns prefix of
// pending bet list of a certain market .
func PendingBetListOfMarketPrefix(marketID string) []byte {
	return append(PendingBetListPrefix, utils.StrBytes(marketID)...)
}

// PendingBetOfMarketKey return the key of
// a certain pending bets of a market.
func PendingBetOfMarketKey(marketID string, id uint64) []byte {
	return append(utils.StrBytes(marketID), utils.Uint64ToBytes(id)...)
}

// SettledBetListOfBlockHeightPrefix returns prefix of
// settled bet list on a certain block height.
func SettledBetListOfBlockHeightPrefix(blockHeight int64) []byte {
	return append(SettledBetListPrefix, utils.Int64ToBytes(blockHeight)...)
}

// SettledBetOfMarketKey return the key of
// settled bet list on a certain block height.
func SettledBetOfMarketKey(blockHeight int64, id uint64) []byte {
	return append(utils.Int64ToBytes(blockHeight), utils.Uint64ToBytes(id)...)
}
