package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

// CalculateDecimalPayout calculates total payout of a certain bet amount by decimal odds calculations
func CalculateDecimalPayout(oddsVal string, amount sdkmath.Int) (sdk.Dec, error) {
	// decimal odds value should be sdk.Dec, so convert it directly
	oddsDecVal, err := sdk.NewDecFromStr(oddsVal)
	if err != nil {
		return sdk.ZeroDec(),
			sdkerrtypes.Wrapf(ErrDecimalOddsIncorrectFormat, "%s", err)
	}

	// odds value should not be negative or zero
	if !oddsDecVal.IsPositive() {
		return sdk.ZeroDec(),
			sdkerrtypes.Wrapf(ErrDecimalOddsShouldBePositive, "%s", oddsVal)
	}

	// odds value should not be less than 1
	if oddsDecVal.LTE(sdk.OneDec()) {
		return sdk.ZeroDec(),
			sdkerrtypes.Wrapf(ErrDecimalOddsCanNotBeLessThanOne, "%s", oddsVal)
	}

	// calculate payout
	payout := oddsDecVal.MulInt(amount)

	// get the integer part of the payout
	return payout, nil
}

// CalculateBetAmount calculates bet amount
func CalculateDecimalBetAmount(oddsVal string, payoutProfit sdk.Dec) (sdk.Dec, error) {
	// decimal odds value should be sdk.Dec, so convert it directly
	oddsDecVal, err := sdk.NewDecFromStr(oddsVal)
	if err != nil {
		return sdk.ZeroDec(),
			sdkerrtypes.Wrapf(ErrDecimalOddsIncorrectFormat, "%s", err)
	}

	// odds value should not be negative or zero
	if !oddsDecVal.IsPositive() {
		return sdk.ZeroDec(),
			sdkerrtypes.Wrapf(ErrDecimalOddsShouldBePositive, "%s", oddsVal)
	}

	// odds value should not be less than 1
	if oddsDecVal.LTE(sdk.OneDec()) {
		return sdk.ZeroDec(),
			sdkerrtypes.Wrapf(ErrDecimalOddsCanNotBeLessThanOne, "%s", oddsVal)
	}

	// calculate bet amount
	betAmount := payoutProfit.Quo(oddsDecVal.Sub(sdk.OneDec()))

	// get the integer part of the bet amount
	return betAmount, nil
}
