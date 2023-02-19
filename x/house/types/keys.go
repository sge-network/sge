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
	// HouseParticipationFeeName defines the account name for house participation fee
	HouseParticipationFeeName = "house_participation_fee_pool"
)

var (
	DepositKeyPrefix    = []byte{0x00} // prefix for keys that store deposits
	WithdrawalKeyPrefix = []byte{0x01} // prefix for keys that store withdrawals
)

// GetDepositKey creates the key for deposit bond with sport event and participant
func GetDepositKey(depositorAddr string, sportEventUID string, participantID uint64) []byte {
	return append(GetDepositListPrefix(depositorAddr), append(utils.StrBytes(sportEventUID), utils.Uint64ToBytes(participantID)...)...)
}

// GetDepositListPrefix creates the key for deposit bond with sport event
func GetDepositListPrefix(depositorAddr string) []byte {
	return utils.StrBytes(depositorAddr)
}

// GetWithdrawalKey creates the key for withdrawal bond with sport event and deposit
func GetWithdrawalKey(depositorAddr string, sportEventUID string, participantID uint64, id uint64) []byte {
	return append(GetDepositKey(depositorAddr, sportEventUID, participantID), utils.Uint64ToBytes(id)...)
}

// GetWithdrawalListPrefix creates the key for withdrawals bond with sport event
func GetWithdrawalListPrefix(depositorAddr string) []byte {
	return utils.StrBytes(depositorAddr)
}
