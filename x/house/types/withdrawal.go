package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "gopkg.in/yaml.v2"
)

// NewWithdrawal creates a new withdrawal object
//
//nolint:interfacer
func NewWithdrawal(id uint64, depositorAddr, sportEventUID string, participationIndex uint64, amount sdk.Int, mode WithdrawalMode) Withdrawal {
	return Withdrawal{
		ID:                 id,
		DepositorAddress:   depositorAddr,
		SportEventUID:      sportEventUID,
		ParticipationIndex: participationIndex,
		Mode:               mode,
		Amount:             amount,
	}
}

// MustMarshalWithdrawal returns the withdrawal bytes. Panics if fails
func MustMarshalWithdrawal(cdc codec.BinaryCodec, withdrawal Withdrawal) []byte {
	return cdc.MustMarshal(&withdrawal)
}

// return the withdrawal
func UnmarshalWithdrawal(cdc codec.BinaryCodec, value []byte) (withdrawal Withdrawal, err error) {
	err = cdc.Unmarshal(value, &withdrawal)
	return withdrawal, err
}

// String returns a human readable string representation of a Withdrawal.
func (w Withdrawal) String() string {
	out, err := yaml.Marshal(w)
	if err != nil {
		panic(err)
	}
	return string(out)
}
