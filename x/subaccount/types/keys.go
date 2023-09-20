package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/sge-network/sge/utils"
)

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

	// LockedBalancePrefix is the key used to store the locked balance in the keeper KVStore
	LockedBalancePrefix = []byte{0x03}

	// BalancePrefix saves the balance of an account.
	BalancePrefix = []byte{0x04}
)

func SubAccountOwnerKey(address sdk.AccAddress) []byte {
	return append(SubAccountOwnerPrefix, address...)
}

func SubAccountKey(subAccountAddress sdk.AccAddress) []byte {
	return append(SubAccountOwnerReversePrefix, subAccountAddress...)
}

func LockedBalanceKey(subAccountAddress sdk.AccAddress, unlockTime uint64) []byte {
	return append(LockedBalancePrefix, append(address.MustLengthPrefix(subAccountAddress), utils.Uint64ToBytes(unlockTime)...)...)
}

func LockedBalancePrefixKey(subAccountAddress sdk.AccAddress) []byte {
	return append(LockedBalancePrefix, address.MustLengthPrefix(subAccountAddress)...)
}

func BalanceKey(address sdk.AccAddress) []byte {
	return append(BalancePrefix, address.Bytes()...)
}
