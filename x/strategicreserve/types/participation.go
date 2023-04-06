package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewOrderBookParticipation creates a new book participation object
//
//nolint:interfacer
func NewOrderBookParticipation(
	index uint64, orderBookUID string, participantAddress string,
	exposuresNotFilled uint64, isModuleAccount bool,
	liquidity, currentRoundLiquidity, totalBetAmount, currentRoundTotalBetAmount, maxLoss, currentRoundMaxLoss sdk.Int,
	currentRoundMaxLossOddsUID string, actualProfit sdk.Int,
) OrderBookParticipation {
	return OrderBookParticipation{
		Index:                      index,
		OrderBookUID:               orderBookUID,
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
func (bp OrderBookParticipation) String() string {
	out, err := yaml.Marshal(bp)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// CalculateMaxLoss calculates the maxixmum amount of the tokens expected to be the
// loss of the participation according to the bet amount
func (bp OrderBookParticipation) CalculateMaxLoss(betAmount sdk.Int) sdk.Int {
	return bp.CurrentRoundMaxLoss.Sub(betAmount)
}
