package types

const (
	// ModuleName defines the module name
	ModuleName = "reward"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_reward"
)

// module accounts constants
const (
	// rewardPool defines the account name for reward pool.
	rewardPool = "reward_pool"
)
