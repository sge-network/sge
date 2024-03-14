package utils

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalculatePriceReimbursement(sgeAmount sdkmath.Int, initialPrice, finalPrice sdk.Dec) sdkmath.Int {
	firstValueInDollars := sdk.NewDecFromInt(sgeAmount).Mul(initialPrice)
	secondValueInDollars := sdk.NewDecFromInt(sgeAmount).Mul(finalPrice)
	// sge reimbursement tokens
	reimbursementInDollar := firstValueInDollars.Sub(secondValueInDollars)
	reimbursementInSGE := reimbursementInDollar.Quo(finalPrice)
	return reimbursementInSGE.TruncateInt()
}
