package types

func NewVote(pubKey string, vote ProposalVote) *Vote {
	return &Vote{
		PublicKey: pubKey,
		Vote:      vote,
	}
}
