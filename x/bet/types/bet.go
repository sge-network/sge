package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

func NewBetFulfillment(
	participantAddress string,
	participationIndex uint64,
	betAmount, payoutProfit sdk.Int,
) *BetFulfillment {
	return &BetFulfillment{
		ParticipantAddress: participantAddress,
		ParticipationIndex: participationIndex,
		BetAmount:          betAmount,
		PayoutProfit:       payoutProfit,
	}
}
