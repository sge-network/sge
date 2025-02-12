package types

import (
	sdkmath "cosmossdk.io/math"
)

// CalculatePayoutProfit calculates the amount of payout profit portion according to bet odds value and amount
func CalculatePayoutProfit(oddsVal string, amount sdkmath.Int) (sdkmath.LegacyDec, error) {
	payout, err := calculatePayout(oddsVal, amount)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	// bettor profit is the subtracted amount of payout from bet amount
	profit := payout.Sub(sdkmath.LegacyNewDecFromInt(amount))

	return profit, nil
}

// calculatePayout calculates the amount of payout according to bet odds value and amount
func calculatePayout(oddsVal string, amount sdkmath.Int) (sdkmath.LegacyDec, error) {
	// total payout should be paid to bettor
	payout, err := CalculateDecimalPayout(oddsVal, amount)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	return payout, nil
}

// CalculateBetAmount calculates the amount of bet according to bet odds value and payout profit
func CalculateBetAmount(oddsVal string, payoutProfit sdkmath.LegacyDec) (sdkmath.LegacyDec, error) {
	betAmount, err := calculateBetAmount(oddsVal, payoutProfit)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	return betAmount, nil
}

// CalculateBetAmountInt calculates the amount of bet according to bet odds value and payout profit
// and returns the int and the truncated decimal part.
func CalculateBetAmountInt(
	oddsVal string,
	payoutProfit sdkmath.LegacyDec,
	truncatedBetAmount sdkmath.LegacyDec,
) (sdkmath.Int, sdkmath.LegacyDec, error) {
	expectedBetAmountDec, err := CalculateBetAmount(oddsVal, payoutProfit)
	if err != nil {
		return sdkmath.Int{}, sdkmath.LegacyDec{}, err
	}
	// add previous loop truncated value to the calculated bet amount
	expectedBetAmountDec = expectedBetAmountDec.Add(truncatedBetAmount)

	// we need for the bet amount to be of type sdkmath.Int
	// so the truncation in inevitable
	betAmount := expectedBetAmountDec.RoundInt()

	// save the truncated amount in the calculations for the next loop
	truncatedBetAmount = truncatedBetAmount.Add(expectedBetAmountDec.Sub(sdkmath.LegacyNewDecFromInt(betAmount)))

	return betAmount, truncatedBetAmount, nil
}

// calculateBetAmount calculates the amount of bet according to bet odds value and payoutProfit
func calculateBetAmount(oddsVal string, payoutProfit sdkmath.LegacyDec) (sdkmath.LegacyDec, error) {
	// total payout should be paid to bettor
	betAmount, err := CalculateDecimalBetAmount(oddsVal, payoutProfit)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	return betAmount, nil
}
