package utils

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
)

// used for skiping the "0." part of string
const decimalPointSkip = 2

// DecreaseDecPrecision decrease the precision of a sdk.Dec
// ex. 1.318791351381317535 => 1.320000000000000000 when the precision is 2 digits
func DecreaseDecPrecision(d sdk.Dec, precision int) sdk.Dec {
	// extract integer part of dec
	intPart := d.TruncateDec()

	// convert dec part to string to be able to replace digits with zero
	decPartStr := d.Sub(intPart).String()

	// pick the decimal part plus next digit to be able to round correctly
	charsToPick := decimalPointSkip + precision + 1

	// extract dec part as an integer
	// ex. 318 when the precision is 2
	decPartInt := cast.ToInt(decPartStr[decimalPointSkip:charsToPick])

	// round to nearest multiplication of 10
	// ex. 320
	dPartRounded := RoundToNearest(decPartInt, 10)

	// build the rounded dec part of the number
	roundedDec := sdk.MustNewDecFromStr("0." + cast.ToString(dPartRounded))

	// add int and dec part and return the result
	return intPart.Add(roundedDec)
}
