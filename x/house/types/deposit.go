package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewDeposit creates a new deposit object
func NewDeposit(
	creator, depositorAddress, marketUID string,
	amount, totalAmount sdkmath.Int,
	withdrawalCount uint64,
) Deposit {
	return Deposit{
		Creator:               creator,
		DepositorAddress:      depositorAddress,
		MarketUID:             marketUID,
		Amount:                amount,
		WithdrawalCount:       withdrawalCount,
		TotalWithdrawalAmount: totalAmount,
	}
}

// String returns a human-readable string representation of a Deposit.
func (d *Deposit) String() string {
	out, err := yaml.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// CalcHouseParticipationFeeAmount sets participation fee amount for house
func (d *Deposit) CalcHouseParticipationFeeAmount(feePercentage sdk.Dec) sdkmath.Int {
	return feePercentage.MulInt(d.Amount).RoundInt()
}
