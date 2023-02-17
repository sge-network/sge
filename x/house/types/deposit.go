package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewDeposit creates a new deposit object
//
//nolint:interfacer
func NewDeposit(depAddr sdk.AccAddress, seUID string, depAmt, withAmt sdk.Int, withdrawals uint64) Deposit {
	return Deposit{
		DepositorAddress:      depAddr.String(),
		SportEventUID:         seUID,
		Amount:                depAmt,
		Withdrawals:           withdrawals,
		TotalWithdrawalAmount: withAmt,
	}
}

// MustMarshalDeposit returns the deposit bytes. Panics if fails
func MustMarshalDeposit(cdc codec.BinaryCodec, deposit Deposit) []byte {
	return cdc.MustMarshal(&deposit)
}

// MustUnmarshalDeposit return the unmarshaled deposit from bytes.
// Panics if fails.
func MustUnmarshalDeposit(cdc codec.BinaryCodec, value []byte) Deposit {
	deposit, err := UnmarshalDeposit(cdc, value)
	if err != nil {
		panic(err)
	}

	return deposit
}

// return the deposit
func UnmarshalDeposit(cdc codec.BinaryCodec, value []byte) (deposit Deposit, err error) {
	err = cdc.Unmarshal(value, &deposit)
	return deposit, err
}

// String returns a human readable string representation of a Deposit.
func (d Deposit) String() string {
	out, _ := yaml.Marshal(d)
	return string(out)
}

func (d *Deposit) SetHouseParticipationFee(feePercentage sdk.Dec) {
	d.Fee = feePercentage.MulInt(d.Amount).RoundInt()
	d.Liquidity = d.Amount.Sub(d.Fee)
}
