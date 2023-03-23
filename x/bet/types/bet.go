package types

func NewPendingBet(uid, creator string) *PendingBet {
	return &PendingBet{
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
