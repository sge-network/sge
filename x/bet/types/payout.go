package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalculatePayoutProfit calculates the amount of payout profit portion according to bet odds value and amount
func CalculatePayoutProfit(oddsType OddsType, oddsVal string, amount sdkmath.Int) (sdk.Dec, error) {
	payout, err := calculatePayout(oddsType, oddsVal, amount)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	// bettor profit is the subtracted amount of payout from bet amount
	profit := payout.Sub(sdk.NewDecFromInt(amount))

	return profit, nil
}

// calculatePayout calculates the amount of payout according to bet odds value and amount
func calculatePayout(oddsType OddsType, oddsVal string, amount sdkmath.Int) (sdk.Dec, error) {
	var oType OddsTypeI

	// assign corresponding type to the interface instance
	switch oddsType {
	case OddsType_ODDS_TYPE_DECIMAL:
		oType = new(decimalOdds)

	case OddsType_ODDS_TYPE_FRACTIONAL:
		oType = new(fractionalOdds)

	case OddsType_ODDS_TYPE_MONEYLINE:
		oType = new(moneylineOdds)

	default:
		return sdk.ZeroDec(), ErrInvalidOddsType
	}

	// total payout should be paid to bettor
	payout, err := oType.CalculatePayout(oddsVal, amount)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return payout, nil
}

// CalculateBetAmount calculates the amount of bet according to bet odds value and payout profit
func CalculateBetAmount(oddsType OddsType, oddsVal string, payoutProfit sdk.Dec) (sdk.Dec, error) {
	betAmount, err := calculateBetAmount(oddsType, oddsVal, payoutProfit)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return betAmount, nil
}

// CalculateBetAmountInt calculates the amount of bet according to bet odds value and payout profit
// and returns the int and the truncated decimal part.
func CalculateBetAmountInt(
	oddsType OddsType,
	oddsVal string,
	payoutProfit sdk.Dec,
	truncatedBetAmount sdk.Dec,
) (sdkmath.Int, sdk.Dec, error) {
	expectedBetAmountDec, err := CalculateBetAmount(oddsType, oddsVal, payoutProfit)
	if err != nil {
		return sdkmath.Int{}, sdk.Dec{}, err
	}
	// add previous loop truncated value to the calculated bet amount
	expectedBetAmountDec = expectedBetAmountDec.Add(truncatedBetAmount)

	// we need for the bet amount to be of type sdkmath.Int
	// so the truncation in inevitable
	betAmount := expectedBetAmountDec.TruncateInt()

	// save the truncated amount in the calculations for the next loop
	truncatedBetAmount = truncatedBetAmount.Add(expectedBetAmountDec.Sub(sdk.NewDecFromInt(betAmount)))

	return betAmount, truncatedBetAmount, nil
}

// calculateBetAmount calculates the amount of bet according to bet odds value and payoutProfit
func calculateBetAmount(oddsType OddsType, oddsVal string, payoutProfit sdk.Dec) (sdk.Dec, error) {
	var oType OddsTypeI

	// assign corresponding type to the interface instance
	switch oddsType {
	case OddsType_ODDS_TYPE_DECIMAL:
		oType = new(decimalOdds)

	case OddsType_ODDS_TYPE_FRACTIONAL:
		oType = new(fractionalOdds)

	case OddsType_ODDS_TYPE_MONEYLINE:
		oType = new(moneylineOdds)

	default:
		return sdk.ZeroDec(), ErrInvalidOddsType
	}

	// total payout should be paid to bettor
	betAmount, err := oType.CalculateBetAmount(oddsVal, payoutProfit)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	return betAmount, nil
}
