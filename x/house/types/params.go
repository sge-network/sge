package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	yaml "gopkg.in/yaml.v2"
)

// House params default values
const (
	// Default minimum deposit acceptable.
	DefaultMinDeposit int64 = 100

	// Default house participation fee.
	DefaultHouseParticipationFee string = "0.1"
)

var (
	KeyMinDeposit            = []byte("MinDeposit")
	KeyHouseParticipationFee = []byte("HouseParticipationFee")
)

// ParamTable for house module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(minDeposit sdk.Int, houseParticipationFee sdk.Dec) Params {
	return Params{
		MinDeposit:            minDeposit,
		HouseParticipationFee: houseParticipationFee,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinDeposit, &p.MinDeposit, validateMinimumDeposit),
		paramtypes.NewParamSetPair(KeyHouseParticipationFee, &p.HouseParticipationFee, validateHouseParticipationFee),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		sdk.NewInt(DefaultMinDeposit),
		sdk.MustNewDecFromStr(DefaultHouseParticipationFee),
	)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// validate a set of params
func (p Params) Validate() error {
	if err := validateMinimumDeposit(p.MinDeposit); err != nil {
		return err
	}

	if err := validateHouseParticipationFee(p.HouseParticipationFee); err != nil {
		return err
	}

	return nil
}

func validateMinimumDeposit(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LTE(sdk.OneInt()) {
		return fmt.Errorf("minimum deposit must be positive and more than one: %d", v)
	}

	return nil
}

func validateHouseParticipationFee(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("house participation fee cannot be lower than 0: %d", v)
	}

	return nil
}
