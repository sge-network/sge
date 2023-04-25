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

func (proposal *PublicKeysChangeProposal) DecideResult(keyvault *KeyVault) ProposalResult {
	var yesCount, noCount int64
	for _, v := range proposal.Votes {
		switch v.Vote {
		case ProposalVote_PROPOSAL_VOTE_YES:
			yesCount++
		case ProposalVote_PROPOSAL_VOTE_NO:
			noCount++
		}
	}

	// minimum accepted count of yes votes (note vote) for the decision.
	majorityCount := keyvault.MajorityCount()

	// check if minimum majority vote count is met or not
	if noCount >= majorityCount {
		return ProposalResult_PROPOSAL_RESULT_REJECTED
	} else if yesCount > majorityCount {
		return ProposalResult_PROPOSAL_RESULT_APPROVED
	}

	// else if the yes and no votes counts are equal and we can not make decision for
	// result of the proposal
	return ProposalResult_PROPOSAL_RESULT_UNSPECIFIED
}
