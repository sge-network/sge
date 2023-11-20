package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	yaml "gopkg.in/yaml.v2"
)

var (
	defaultWagerEnabled   = true
	defaultDepositEnabled = false
)

// parameter store keys
var (
	// keyWagerEnabled is the enable/disable status of subaccount wager tx endpoint.
	keyWagerEnabled = []byte("WagerEnabled")

	// keyDepositEnabled is the enable/disable status of subaccount deposit tx endpoint.
	keyDepositEnabled = []byte("DepositEnabled")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		// TODO: in the next upgrade handler
		WagerEnabled: defaultWagerEnabled,
		// TODO: in the next upgrade handler
		DepositEnabled: defaultDepositEnabled,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams()
}

// ParamKeyTable for house module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			keyWagerEnabled,
			&p.WagerEnabled,
			validateWagerEnabled,
		),
		paramtypes.NewParamSetPair(
			keyDepositEnabled,
			&p.DepositEnabled,
			validateDepositEnabled,
		),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateWagerEnabled(p.WagerEnabled); err != nil {
		return err
	}

	return validateDepositEnabled(p.DepositEnabled)
}

// String returns a human-readable string representation of the parameters.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func validateWagerEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("%s: %T", "invalid parameter type for wager enabled", i)
	}

	return nil
}

func validateDepositEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("%s: %T", "invalid parameter type for deposit enabled", i)
	}

	return nil
}
