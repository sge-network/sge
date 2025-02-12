package types

import (
	"github.com/sge-network/sge/utils"
)

const (
	// ModuleName is the name of the house module
	ModuleName = "house"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the house module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the house module
	RouterKey = ModuleName
)

// module accounts constants
const (
	// houseFeeCollector defines the account name for house participation fee
	houseFeeCollector = "house_fee_collector"
)

var (
	DepositKeyPrefix    = []byte{0x00} // prefix for keys that store deposits
	WithdrawalKeyPrefix = []byte{0x01} // prefix for keys that store withdrawals
)

// GetDepositKey creates the key for deposit bond with market and participation
func GetDepositKey(depositorAddr, marketUID string, participationIndex uint64) []byte {
	return append(
		GetDepositListPrefix(depositorAddr),
		append(utils.StrBytes(marketUID), utils.Uint64ToBytes(participationIndex)...)...)
}

// GetDepositListPrefix creates the key for deposit bond with market
func GetDepositListPrefix(depositorAddr string) []byte {
	return utils.StrBytes(depositorAddr)
}

// GetWithdrawalKey creates the key for withdrawal bond with market and deposit
func GetWithdrawalKey(
	depositorAddr string,
	marketUID string,
	participationIndex uint64,
	id uint64,
) []byte {
	return append(
		GetDepositKey(depositorAddr, marketUID, participationIndex),
		utils.Uint64ToBytes(id)...)
}

// GetWithdrawalListPrefix creates the key for withdrawals bond with market
func GetWithdrawalListPrefix(depositorAddr string) []byte {
	return utils.StrBytes(depositorAddr)
}
