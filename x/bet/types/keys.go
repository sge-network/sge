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

var (
	// BetListPrefix is the prefix to retrieve all Bet
	BetListPrefix = []byte{0x00}

	// BetIDListPrefix is the prefix to retrieve all Bet IDs
	BetIDListPrefix = []byte{0x01}

	// BetStatsKey is the key for the bet statistics
	BetStatsKey = []byte{0x03}
)

func BetListByCreatorKey(creator string) []byte {
	return append(BetListPrefix, utils.StrBytes(creator)...)
}

func BetListByIDKey(creator string, id uint64) []byte {
	return append(utils.StrBytes(creator), utils.Uint64ToBytes(id)...)
}
