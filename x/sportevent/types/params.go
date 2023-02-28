package types

import (
	"fmt"

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

	// KeyMaxSRContribution is the min bet amount param key
	KeyMaxSRContribution = []byte("EventMaxSRContribution")
)

// default params
var (
	// DefaultMinBetAmount is the default minimum bet amount allowed
	DefaultMinBetAmount = sdk.NewInt(1000000)

	// DefaultMaxSRContribution is the default maximum sr contribution allowed
	DefaultMaxSRContribution = sdk.NewInt(10000000)

	// DefaultMinBetFee is the default minimum bet fee amount allowed
	DefaultMinBetFee = sdk.NewInt(0)
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		EventMinBetAmount:      DefaultMinBetAmount,
		EventMinBetFee:         DefaultMinBetFee,
		EventMaxSrContribution: DefaultMaxSRContribution,
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
		paramtypes.NewParamSetPair(KeyMaxSRContribution, &p.EventMaxSrContribution, validateMaxSRContribution),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func validateMinBetAmount(i interface{}) error {
	return nil
}

func validateMinBetFeePercentage(i interface{}) error {
	return nil
}

func validateMaxSRContribution(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LTE(sdk.OneInt()) {
		return fmt.Errorf("max SR contribution must be positive: %d", v.Int64())
	}

	return nil
}

// NewEventBetConstraints creates new bet constraint pointer
func (p *Params) NewEventBetConstraints(minAmount, betFee sdk.Int) *EventBetConstraints {
	if minAmount.IsNil() {
		minAmount = p.EventMinBetAmount
	}

	if betFee.IsNil() {
		betFee = p.EventMinBetFee
	}

	return &EventBetConstraints{
		MinAmount: minAmount,
		BetFee:    betFee,
	}
}
