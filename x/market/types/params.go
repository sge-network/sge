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
	// KeyMinBetAmount is the min bet amount param key
	KeyMinBetAmount = []byte("MinBetAmount")

	// KeyMinBetFee is the minimum bet fee param key
	KeyMinBetFee = []byte("MinBetFee")

	// KeyMaxBetFee is the maximum bet fee param key
	KeyMaxBetFee = []byte("MaxBetFee")

	// KeyMaxSRContribution is the min bet amount param key
	KeyMaxSRContribution = []byte("MaxSRContribution")
)

// default params
var (
	// DefaultMinBetAmount is the default minimum bet amount allowed
	DefaultMinBetAmount = sdk.NewInt(1000000)

	// DefaultMaxSRContribution is the default maximum sr contribution allowed
	DefaultMaxSRContribution = sdk.NewInt(10000000)

	// DefaultMinBetFee is the default minimum bet fee amount allowed
	DefaultMinBetFee = sdk.NewInt(0)

	// DefaultMaxBetFee is the default maximum bet fee amount allowed
	DefaultMaxBetFee = sdk.NewInt(100)
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		MinBetAmount:      DefaultMinBetAmount,
		MinBetFee:         DefaultMinBetFee,
		MaxBetFee:         DefaultMaxBetFee,
		MaxSrContribution: DefaultMaxSRContribution,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinBetAmount, &p.MinBetAmount, validateMinBetAmount),
		paramtypes.NewParamSetPair(KeyMinBetFee, &p.MinBetFee, validateMinBetFeePercentage),
		paramtypes.NewParamSetPair(KeyMaxBetFee, &p.MaxBetFee, validateMaxBetFeePercentage),
		paramtypes.NewParamSetPair(KeyMaxSRContribution, &p.MaxSrContribution, validateMaxSRContribution),
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
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LTE(sdk.OneInt()) {
		return fmt.Errorf("minimum bet amount must be positive: %d", v.Int64())
	}

	return nil
}

func validateMinBetFeePercentage(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroInt()) {
		return fmt.Errorf("minimum bet fee must be positive: %d", v.Int64())
	}

	return nil
}

func validateMaxBetFeePercentage(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroInt()) {
		return fmt.Errorf("maximum bet fee must be positive: %d", v.Int64())
	}

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

// NewMarketBetConstraints creates new bet constraint pointer
func (p *Params) NewMarketBetConstraints(minAmount, betFee sdk.Int) *MarketBetConstraints {
	if minAmount.IsNil() {
		minAmount = p.MinBetAmount
	}

	if betFee.IsNil() {
		betFee = p.MinBetFee
	}

	return &MarketBetConstraints{
		MinAmount: minAmount,
		BetFee:    betFee,
	}
}
