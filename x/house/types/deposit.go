package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewDeposit creates a new deposit object
//
//nolint:interfacer
func NewDeposit(creator, sportEventUID string, amount, totalAmount sdk.Int, withdrawalCount uint64) Deposit {
	return Deposit{
		Creator:               creator,
		SportEventUID:         sportEventUID,
		Amount:                amount,
		WithdrawalCount:       withdrawalCount,
		TotalWithdrawalAmount: totalAmount,
	}
}

// String returns a human readable string representation of a Deposit.
func (d Deposit) String() string {
	out, err := yaml.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func (d *Deposit) SetHouseParticipationFee(feePercentage sdk.Dec) {
	d.Fee = feePercentage.MulInt(d.Amount).RoundInt()
	d.Liquidity = d.Amount.Sub(d.Fee)
}
