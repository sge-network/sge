package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

// CalculateDecimalPayout calculates total payout of a certain bet amount by decimal odds calculations
func CalculateDecimalPayout(oddsVal string, amount sdkmath.Int) (sdkmath.LegacyDec, error) {
	// decimal odds value should be sdkmath.LegacyDec, so convert it directly
	oddsDecVal, err := sdkmath.LegacyNewDecFromStr(oddsVal)
	if err != nil {
		return sdkmath.LegacyZeroDec(),
			sdkerrors.Wrapf(ErrDecimalOddsIncorrectFormat, "%s", err)
	}

	// odds value should not be negative or zero
	if !oddsDecVal.IsPositive() {
		return sdkmath.LegacyZeroDec(),
			sdkerrors.Wrapf(ErrDecimalOddsShouldBePositive, "%s", oddsVal)
	}

	// odds value should not be less than 1
	if oddsDecVal.LTE(sdkmath.LegacyOneDec()) {
		return sdkmath.LegacyZeroDec(),
			sdkerrors.Wrapf(ErrDecimalOddsCanNotBeLessThanOne, "%s", oddsVal)
	}

	// calculate payout
	payout := oddsDecVal.MulInt(amount)

	// get the integer part of the payout
	return payout, nil
}

// CalculateBetAmount calculates bet amount
func CalculateDecimalBetAmount(oddsVal string, payoutProfit sdkmath.LegacyDec) (sdkmath.LegacyDec, error) {
	// decimal odds value should be sdkmath.LegacyDec, so convert it directly
	oddsDecVal, err := sdkmath.LegacyNewDecFromStr(oddsVal)
	if err != nil {
		return sdkmath.LegacyZeroDec(),
			sdkerrors.Wrapf(ErrDecimalOddsIncorrectFormat, "%s", err)
	}

	// odds value should not be negative or zero
	if !oddsDecVal.IsPositive() {
		return sdkmath.LegacyZeroDec(),
			sdkerrors.Wrapf(ErrDecimalOddsShouldBePositive, "%s", oddsVal)
	}

	// odds value should not be less than 1
	if oddsDecVal.LTE(sdkmath.LegacyOneDec()) {
		return sdkmath.LegacyZeroDec(),
			sdkerrors.Wrapf(ErrDecimalOddsCanNotBeLessThanOne, "%s", oddsVal)
	}

	// calculate bet amount
	betAmount := payoutProfit.Quo(oddsDecVal.Sub(sdkmath.LegacyOneDec()))

	// get the integer part of the bet amount
	return betAmount, nil
}
