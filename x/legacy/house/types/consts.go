package types

import sdkmath "cosmossdk.io/math"

var (
	// minDepositGrant is the minimum deposit allowed grant.
	minDepositGrant = sdkmath.NewInt(100)
	// maxWithdrawGrant is the maximum withdraw allowed grant.
	maxWithdrawGrant = sdkmath.NewInt(100)
)
