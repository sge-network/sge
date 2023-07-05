package types

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sge-network/sge/app/params"
	"gopkg.in/yaml.v2"
)

const DefaultBondDenom = params.DefaultBondDenom

var (
	KeyLockedBalanceDenom = []byte("LockedBalanceDenom")
)

// NewParams creates a new Params instance
func NewParams(lockedBalanceDenom string) Params {
	return Params{LockedBalanceDenom: lockedBalanceDenom}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultBondDenom,
	)
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			KeyLockedBalanceDenom,
			&p.LockedBalanceDenom,
			validateLockedBalanceDenom,
		),
	}
}

func validateLockedBalanceDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == "" {
		return fmt.Errorf("locked balance denom should not be empty")
	}

	return nil
}

// String returns a human-readable string representation of the parameters.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// Validate a set of params
func (p Params) Validate() error {
	if err := validateLockedBalanceDenom(p.LockedBalanceDenom); err != nil {
		return err
	}

	return nil
}

// ParamKeyTable for house module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}
