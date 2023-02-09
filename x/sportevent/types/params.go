package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// param keys
var (
	// KeyEventMinBetAmount is the min bet amount param key
	KeyEventMinBetAmount = []byte("EventMinBetAmount")

	// KeyMinBetFee is the minimum bet fee param key
	KeyMinBetFee = []byte("EventMinBetFee")
)

// default params
const (
	// DefaultMinBetAmount is the default minimum bet amount allowed
	DefaultMinBetAmount = 1000000
)

// DefaultMinBetFee is the default minimum bet fee amount allowed
var DefaultMinBetFee = sdk.NewInt(100000)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		EventMinBetAmount: sdk.NewInt(DefaultMinBetAmount),
		EventMinBetFee:    DefaultMinBetFee,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEventMinBetAmount, &p.EventMinBetAmount, validateMinBetAmount),
		paramtypes.NewParamSetPair(KeyMinBetFee, &p.EventMinBetFee, validateMinBetFeePercentage),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateMinBetAmount(i interface{}) error {
	return nil
}

func validateMinBetFeePercentage(i interface{}) error {
	return nil
}
