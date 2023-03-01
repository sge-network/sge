package types

import (
	fmt "fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

const (
	batchSettlementCount = 1000
)

// parameter store keys
var (
	// KeyBatchSettlementCount is the batch settlement
	// count of bets
	KeyBatchSettlementCount = []byte("BatchSettlementCount")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		BatchSettlementCount: batchSettlementCount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBatchSettlementCount, &p.BatchSettlementCount, validateBatchSettlementCount),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateBatchSettlementCount(p.BatchSettlementCount); err != nil {
		return err
	}
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

func validateBatchSettlementCount(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf(ErrTextInvalidParamType, i)
	}

	if v <= 0 {
		return fmt.Errorf(ErrTextBatchSettlementCountMustBePositive, v)
	}

	return nil
}
