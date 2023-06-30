package types

// module constants
const (
	// ModuleName defines the module name
	ModuleName = "subaccount"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

var (
	// SubaccountIDPrefix is the key used to store the subaccount ID in the keeper KVStore
	SubaccountIDPrefix = []byte{0x00}
)
