package types

import "github.com/cosmos/cosmos-sdk/types"

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

	// SubAccountOwnerPrefix is the key used to store the subaccount owner in the keeper KVStore
	SubAccountOwnerPrefix = []byte{0x01}

	// SubAccountOwnerReversePrefix is the key used to store the subaccount owner by ID in the keeper KVStore
	SubAccountOwnerReversePrefix = []byte{0x02}
)

func SubAccountOwnerKey(address types.AccAddress) []byte {
	return append(SubAccountOwnerPrefix, address...)
}

func SubAccountKey(id uint64) []byte {
	return append(SubAccountOwnerReversePrefix, types.Uint64ToBigEndian(id)...)
}
