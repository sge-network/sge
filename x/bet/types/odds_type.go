package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type OddsTypeI interface {
	// CalculatePayout calculates total payout of a certain bet amount
	CalculatePayout(oddsVal string, amount sdk.Int) (sdk.Int, error)
}

// decimalOdds is the type to define OddsTypeI interface
// for the decimal odds type
type decimalOdds struct{}

// CalculatePayout calculates total payout of a certain bet amount by decimal odds calculations
func (c *decimalOdds) CalculatePayout(oddsVal string, amount sdk.Int) (sdk.Int, error) {

	// decimal odds value should be sdk.Dec, so convert it directly
	oddsDecVal, err := sdk.NewDecFromStr(oddsVal)
	if err != nil {
		return sdk.ZeroInt(),
			sdkerrors.Wrapf(ErrInConvertingOddsToDec, "%s", err)
	}

	// odds value should not be negative or zero
	if !oddsDecVal.IsPositive() {
		return sdk.ZeroInt(),
			sdkerrors.Wrapf(ErrDecimalOddsCanNotBeNegative, "%s", oddsVal)
	}

	// odds value should not be less than 1
	if oddsDecVal.LT(sdk.NewDec(1)) {
		return sdk.ZeroInt(),
			sdkerrors.Wrapf(ErrDecimalOddsCanNotBeLessThanOne, "%s", oddsVal)
	}

	// calculate payout
	payout := oddsDecVal.MulInt(amount)

	// get the integer part of the payout
	return payout.RoundInt(), nil
}

// fractionalOdds is the type to define OddsTypeI interface
// for the fractional odds type
type fractionalOdds struct{}

// CalculatePayout calculates total payout of a certain bet amount by fractional odds calculations
func (c *fractionalOdds) CalculatePayout(oddsVal string, amount sdk.Int) (sdk.Int, error) {
	fraction := strings.Split(oddsVal, "/")

	// the fraction should contain two parts such as (firstpart)/secondpart)
	if len(fraction) != 2 {
		return sdk.ZeroInt(),
			ErrFractionalOddsIncorrectFormat
	}

	// the fraction part should be an integer, values such as 1.5/2 is not accepted
	firstPart, ok := sdk.NewIntFromString(fraction[0])
	if !ok {
		return sdk.ZeroInt(),
			sdkerrors.Wrapf(ErrInConvertingOddsToInt, "%s", oddsVal)
	}

	// the fraction part should be an integer, values such as 1/2.5 is not accepted
	secondPart, ok := sdk.NewIntFromString(fraction[1])
	if !ok {
		return sdk.ZeroInt(),
			sdkerrors.Wrapf(ErrInConvertingOddsToInt, "%s", oddsVal)
	}

	// fraction parts should be positive
	if !firstPart.IsPositive() || !secondPart.IsPositive() {
		return sdk.ZeroInt(), ErrFractionalOddsCanNotBeNegativeOrZero
	}

	// calculate the coefficient by dividing sdk.Dec values of fraction parts
	// this helps not to lost precision in the division
	coefficient := firstPart.ToDec().Quo(secondPart.ToDec())

	// TODO: decice to decrease the precision or not
	// decrease coefficient to prevent calculation error
	// coefficient = utils.DecreaseDecPrecision(coefficient, payOutCaclulationPrecision)

	// calculate payout
	payout := amount.ToDec().Add(amount.ToDec().Mul(coefficient))

	// get the integer part of the payout
	return payout.RoundInt(), nil
}

// moneylineOdds is the type to define OddsTypeI interface
// for the moneyline odds type
type moneylineOdds struct{}

// CalculatePayout calculates total payout of a certain bet amount by moneyline odds calculations
func (c *moneylineOdds) CalculatePayout(oddsVal string, amount sdk.Int) (sdk.Int, error) {
	// moneyline odds value sould be integer
	oddsValue, ok := sdk.NewIntFromString(oddsVal)
	if !ok {
		return sdk.ZeroInt(),
			sdkerrors.Wrapf(ErrInConvertingOddsToInt, "%s", oddsVal)
	}

	// moneyline values can be negative or positive, but zero is not acceptable
	if oddsValue.IsZero() {
		return sdk.ZeroInt(), ErrMoneylineOddsCanNotBeZero
	}

	// calculate payout
	var payout, coefficient sdk.Dec
	// calculate coefficient of the payout calculations by using sdk.Dec values of odds value
	// we should extract absolute number to prevent negative payout
	if oddsValue.IsPositive() {
		coefficient = oddsValue.ToDec().Quo(sdk.NewDec(100)).Abs()
	} else {
		coefficient = sdk.NewDec(100).QuoInt(oddsValue).Abs()
	}

	// TODO: decice to decrease the precision or not
	// decrease coefficient to prevent calculation error
	// coefficient = utils.DecreaseDecPrecision(coefficient, payOutCaclulationPrecision)

	// bet amount should be multiplied by the coefficient
	payout = amount.ToDec().Add(amount.ToDec().Mul(coefficient))

	// get the integer part of the payout
	return payout.RoundInt(), nil
}
