package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// priceLockPool defines the account name for price locking feature pool.
	priceLockPool = "price_lock_pool"
	// priceLockFeeCollector is the module account name for the price lock fee collection module account.
	priceLockFeeCollector = "price_lock_fee_collector"
)

const (
	// ModuleName defines the module name
	ModuleName = "market"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_market"
)

var (
	// MarketKeyPrefix is the prefix to retrieve all Market
	MarketKeyPrefix = []byte{0x00}

	// MarketStatsKey is the key for the market statistics
	MarketStatsKey = []byte{0x01}
)
