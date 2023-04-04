package types

const (
	// minVoteCountForDecision is the minimum count of votes for a proposal
	// to be valid for making decision.
	minVoteCountForDecision = 3
	// MaxValidProposalMinutes is the maximum elapsed time in minutes since
	// the start time of a proposal to be acceptable.
	MaxValidProposalMinutes = 30
	// MaxValidProposalSeconds is the maximum elapsed time in seconds since
	// the start time of a proposal to be acceptable.
	MaxValidProposalSeconds = MaxValidProposalMinutes * 60 // 30 minutes
)
