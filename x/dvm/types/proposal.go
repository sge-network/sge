package types

func NewPublicKeysChangeProposal(
	id uint64,
	creator string,
	modifications PubkeysChangeProposalPayload,
	startTS int64,
) PublicKeysChangeProposal {
	return PublicKeysChangeProposal{
		Id:            id,
		Creator:       creator,
		Modifications: modifications,
		StartTS:       startTS,
		Status:        ProposalStatus_PROPOSAL_STATUS_ACTIVE,
	}
}

func (proposal *PublicKeysChangeProposal) IsExpired(blockTime int64) bool {
	diff := blockTime - proposal.StartTS
	return diff > MaxValidProposalSeconds
}

func (proposal *PublicKeysChangeProposal) DecideResult() ProposalResult {
	var yesCount, noCount int
	for _, v := range proposal.Votes {
		switch v.Vote {
		case ProposalVote_PROPOSAL_VOTE_YES:
			yesCount++
		case ProposalVote_PROPOSAL_VOTE_NO:
			noCount++
		}
	}

	// check if minimum vote count is met or not
	if yesCount >= minVoteCountForDecision ||
		noCount >= minVoteCountForDecision {
		// minumum vote count is met, so if the yes votes count is more than rejected,
		// the proposal is approved,  otherwise is rejected.
		if yesCount > noCount {
			return ProposalResult_PROPOSAL_RESULT_APPROVED
		} else if yesCount < noCount {
			return ProposalResult_PROPOSAL_RESULT_REJECTED
		}
		// else if the yes and no votes counts are equal and we can not make decision for
		// result of the proposal
	}

	return ProposalResult_PROPOSAL_RESULT_UNSPECIFIED
}
