package types

import "github.com/sge-network/sge/utils"

const (
	// ModuleName defines the module name
	ModuleName = "rewards"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_rewards"
)

var (
	RewardsKeyPrefix = []byte{0x00} // prefix for keys that store rewards
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func GetRewardKey(id string) []byte {
	return utils.StrBytes(id)
}
