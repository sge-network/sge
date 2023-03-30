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
