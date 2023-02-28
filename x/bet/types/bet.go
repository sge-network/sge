package types

func NewActiveBet(uid, creator string) *ActiveBet {
	return &ActiveBet{
		UID:     uid,
		Creator: creator,
	}
}

func NewSettledBet(uid, bettorAddress string) *SettledBet {
	return &SettledBet{
		UID:           uid,
		BettorAddress: bettorAddress,
	}
}
