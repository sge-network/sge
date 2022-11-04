package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalculateVig vig according to the active odds
// vig is extra amount of odds value more than 100
func CalculateVig(odds []*BetOdds) sdk.Dec {
	total := sdk.NewDec(0)
	for _, o := range odds {
		oVal := sdk.MustNewDecFromStr(o.Value)
		oValPerc := sdk.NewDec(vigCalculationCoefficient).Quo(oVal)
		total = total.Add(oValPerc)
	}

	return total.Sub(sdk.NewDec(maxOddsValueSum))
}
