package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	// minDepositGrant is the minimum deposit allowed grant.
	minDepositGrant = sdk.NewInt(100)
	// maxWithdrawGrant is the maximum withdraw allowed grant.
	maxWithdrawGrant = sdk.NewInt(100)
)
