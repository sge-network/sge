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
	}
}

func NewFinishedPublicKeysChangeProposal(
	proposal PublicKeysChangeProposal,
	result ProposalResult,
	resultMetadata string,
	finishTS int64,
) PublicKeysChangeFinishedProposal {
	return PublicKeysChangeFinishedProposal{
		Proposal:   proposal,
		Result:     result,
		ResultMeta: resultMetadata,
		FinishTS:   finishTS,
	}
}

func (proposal *PublicKeysChangeProposal) IsExpired(blockTime int64) bool {
	diff := blockTime - proposal.StartTS
	return diff > MaxValidProposalSeconds
}
