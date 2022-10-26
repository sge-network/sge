package types

// ReserverKey is the key to use for the keeper store.
var ReserverKey = []byte("sr")

// module constants
const (
	// ModuleName defines the module name
	ModuleName = "strategicreserve"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_strategicreserve"
)

// module accounts constants
const (
	// SRPoolName defines the account name for SR Pool
	SRPoolName = "sr_pool"

	// BetReserveName defines the account name for storing bet amount
	BetReserveName = "bet_reserve"

	// WinningsCollectorName defines the account for the storing the
	// winnings of the house from bet events
	WinningsCollectorName = "winnings_collector"
)

// prefixes
var (
	// PayoutLockPrefix defines the prefix for the KV-Store partition
	// which stores the locks for paying out the funds to the user
	PayoutLockPrefix = []byte{0x01}

	// DailyCapPrefix defines the prefix for the KV-Store partition
	// which stores the current_filled_capacity for the corresponding
	// day_id
	DailyCapPrefix = []byte{0x02}
)
