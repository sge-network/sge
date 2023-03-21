package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewBookParticipation creates a new book participation object
//
//nolint:interfacer
func NewBookParticipation(
	index uint64, bookUID string, participantAddress string,
	exposuresNotFilled uint64, isModuleAccount bool,
	liquidity, currentRoundLiquidity, totalBetAmount, currentRoundTotalBetAmount, maxLoss, currentRoundMaxLoss sdk.Int,
	currentRoundMaxLossOddsUID string, actualProfit sdk.Int,
) BookParticipation {
	return BookParticipation{
		Index:                      index,
		BookUID:                    bookUID,
		ParticipantAddress:         participantAddress,
		IsModuleAccount:            isModuleAccount,
		Liquidity:                  liquidity,
		CurrentRoundLiquidity:      currentRoundLiquidity,
		ExposuresNotFilled:         exposuresNotFilled,
		TotalBetAmount:             totalBetAmount,
		CurrentRoundTotalBetAmount: currentRoundTotalBetAmount,
		MaxLoss:                    maxLoss,
		CurrentRoundMaxLoss:        currentRoundMaxLoss,
		CurrentRoundMaxLossOddsUID: currentRoundMaxLossOddsUID,
		ActualProfit:               actualProfit,
	}
}

// String returns a human readable string representation of a BookParticipation.
func (bp BookParticipation) String() string {
	out, err := yaml.Marshal(bp)
	if err != nil {
		panic(err)
	}
	return string(out)
}
