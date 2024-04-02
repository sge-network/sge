package types

import sdkmath "cosmossdk.io/math"

var (
	// minCampaignFunds is the minimum campaign funds allowed grant.
	minCampaignFunds = sdkmath.NewInt(100)

	// maxWithdrawGrant is the maximum withdraw allowed grant.
	maxWithdrawGrant = sdkmath.NewInt(100)
)
