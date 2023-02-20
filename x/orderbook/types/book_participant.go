package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewBookParticipation creates a new book participation object
//
//nolint:interfacer
func NewBookParticipation(
	index uint64, bookID string, participantAddress string,
	exposuresNotFilled uint64, isModuleAccount bool,
	liquidity, currentRoundLiquidity, totalBetAmount, currentRoundTotalBetAmount, maxLoss, currentRoundMaxLoss sdk.Int,
	currentRoundMaxLossOddsID string, actualProfit sdk.Int,
) BookParticipation {
	return BookParticipation{
		Index:                      index,
		BookID:                     bookID,
		ParticipantAddress:         participantAddress,
		IsModuleAccount:            isModuleAccount,
		Liquidity:                  liquidity,
		CurrentRoundLiquidity:      currentRoundLiquidity,
		ExposuresNotFilled:         exposuresNotFilled,
		TotalBetAmount:             totalBetAmount,
		CurrentRoundTotalBetAmount: currentRoundTotalBetAmount,
		MaxLoss:                    maxLoss,
		CurrentRoundMaxLoss:        currentRoundMaxLoss,
		CurrentRoundMaxLossOddsID:  currentRoundMaxLossOddsID,
		ActualProfit:               actualProfit,
	}
}

// String returns a human readable string representation of a BookParticipation.
func (bp BookParticipation) String() string {
	out, _ := yaml.Marshal(bp)
	return string(out)
}
