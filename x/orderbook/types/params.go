package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	yaml "gopkg.in/yaml.v2"
)

// OrderBook params default values
const (
	// Default maximum book participations.
	DefaultMaxBookParticipations uint64 = 100

	// Default batch settlement count.
	DefaultBatchSettlementCount uint64 = 100
)

var (
	KeyMaxBookParticipations = []byte("MaxBookParticipationss")

	KeyBatchSettlementCount = []byte("BatchSettlementCount")
)

// ParamTable for orderbook module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(maxBookParticipations, batchSettlementCount uint64) Params {
	return Params{
		MaxBookParticipations: maxBookParticipations,
		BatchSettlementCount:  batchSettlementCount,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxBookParticipations, &p.MaxBookParticipations, validateMaxBookParticipations),
		paramtypes.NewParamSetPair(KeyBatchSettlementCount, &p.BatchSettlementCount, validateBatchSettlementCount),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxBookParticipations,
		DefaultBatchSettlementCount,
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
	if err := validateMaxBookParticipations(p.MaxBookParticipations); err != nil {
		return err
	}

	if err := validateBatchSettlementCount(p.BatchSettlementCount); err != nil {
		return err
	}

	return nil
}

func validateMaxBookParticipations(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("maximum book participations must be positive: %d", v)
	}

	return nil
}

func validateBatchSettlementCount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("batch settlement count must be positive: %d", v)
	}

	return nil
}
