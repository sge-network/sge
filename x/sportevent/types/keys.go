package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// ModuleName defines the module name
	ModuleName = "sportevent"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_sportevent"
)

var (
	// SportEventKeyPrefix is the prefix to retrieve all SportEvent
	SportEventKeyPrefix = []byte{0x00}

	// SportEventStatsKey is the key for the sport-event statistics
	SportEventStatsKey = []byte{0x01}
)
