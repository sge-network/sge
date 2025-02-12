package types

// keys
var (
	// MinterKey is the key to use for the keeper store.
	MinterKey = []byte{0x00}
	// ParamsKey is the key to use for the keeper store.
	ParamsKey = []byte{0x01}
)

// module constants
const (
	// ModuleName defines the module name
	ModuleName = "mint"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)
