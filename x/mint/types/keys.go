package types

// module constants
const (
	// ModuleName defines the module name
	ModuleName = "mint"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// keys
var (
	// MinterKey is the key to use for the keeper store.
	MinterKey = []byte{0x00}
)
