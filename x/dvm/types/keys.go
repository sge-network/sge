package types

// module constants
const (
	// ModuleName defines the module name
	ModuleName = "dvm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_dvm"
)

// store prefixes and keys
var (
	// PubKeysListKey is the key of the list of public keys in KV-Store which points to []string
	PubKeysListKey = []byte{0x00}
)
