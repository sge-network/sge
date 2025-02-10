package types

import (
	yaml "gopkg.in/yaml.v2"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
)

// NewWithdrawal creates a new withdrawal object
//
//nolint:interface
func NewWithdrawal(
	id uint64,
	creator, depositorAddr, marketUID string,
	participationIndex uint64,
	amount sdkmath.Int,
	mode WithdrawalMode,
) Withdrawal {
	return Withdrawal{
		Creator:            creator,
		ID:                 id,
		Address:            depositorAddr,
		MarketUID:          marketUID,
		ParticipationIndex: participationIndex,
		Mode:               mode,
		Amount:             amount,
	}
}

// MustMarshalWithdrawal returns the withdrawal bytes. Panics if fails
func MustMarshalWithdrawal(cdc codec.BinaryCodec, withdrawal Withdrawal) []byte {
	return cdc.MustMarshal(&withdrawal)
}

// UnmarshalWithdrawal return the withdrawal
func UnmarshalWithdrawal(cdc codec.BinaryCodec, value []byte) (withdrawal Withdrawal, err error) {
	err = cdc.Unmarshal(value, &withdrawal)
	return withdrawal, err
}

// String returns a human-readable string representation of a Withdrawal.
func (w Withdrawal) String() string {
	out, err := yaml.Marshal(w)
	if err != nil {
		panic(err)
	}
	return string(out)
}
