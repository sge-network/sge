package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// JWT constants
const (
	// JWTHeaderIndex is the index of the header in the JWT ticket
	JWTHeaderIndex = 0

	// JWTPayloadIndex is the index of the payload in the JWT ticket
	JWTPayloadIndex = 1

	// JWTSeparator is the separator character between JWT ticket parts
	JWTSeparator = "."

	// DefaultTimeWeight is the default weight of the time for JWT ticket expiration
	DefaultTimeWeight = 1

	// MinPubKeysCount is the minimum allowed public keys in the key vault
	MinPubKeysCount = 3

	// MaxPubKeysCount is the maximum allowed public keys in the key vault
	MaxPubKeysCount = 5
)

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
