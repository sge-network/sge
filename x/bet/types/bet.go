package types

func NewActiveBet(id uint64, creator string) *ActiveBet {
	return &ActiveBet{
		ID:      id,
		Creator: creator,
	}
}

func NewSettledBet(id uint64, bettorAddress string) *SettledBet {
	return &SettledBet{
		ID:            id,
		BettorAddress: bettorAddress,
	}
}
