package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewDeposit creates a new deposit object
func NewDeposit(creator, marketUID string, amount, totalAmount sdk.Int, withdrawalCount uint64) Deposit {
	return Deposit{
		Creator:               creator,
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

// SetHouseParticipationFee sets participation fee for house
func (d *Deposit) SetHouseParticipationFee(feePercentage sdk.Dec) {
	d.Fee = feePercentage.MulInt(d.Amount).RoundInt()
	d.Liquidity = d.Amount.Sub(d.Fee)
}
