package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (

	// MaxValidProposalMinutes is the maximum elapsed time in minutes since
	// the start time of a proposal to be acceptable.
	MaxValidProposalMinutes = 30
	// MaxValidProposalSeconds is the maximum elapsed time in seconds since
	// the start time of a proposal to be acceptable.
	MaxValidProposalSeconds = MaxValidProposalMinutes * 60 // 30 minutes
)

// minVoteMajorityForDecisionPercentage is the minimum majority percentage of votes for
// a proposal to be valid for making decision.
var minVoteMajorityForDecisionPercentage = sdk.NewDecWithPrec(6667, 4) // 66.67%
