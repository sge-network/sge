package types

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Orderbook params default values
const (
	// DefaultMaxOrderBookParticipations Default maximum book participations.
	DefaultMaxOrderBookParticipations uint64 = 100

	// DefaultBatchSettlementCount Default batch settlement count.
	DefaultBatchSettlementCount uint64 = 100

	// DefaultRequeueThreshold Default requeue threshold.
	DefaultRequeueThreshold uint64 = 1000
)

var (
	keyMaxOrderBookParticipations = []byte("MaxOrderBookParticipations")
	keyBatchSettlementCount       = []byte("BatchSettlementCount")
	keyRequeueThreshold           = []byte("RequeueThreshold")
)

// ParamKeyTable ParamTable for orderbook module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(maxOrderBookParticipations, batchSettlementCount, requeueThreshold uint64) Params {
	return Params{
		MaxOrderBookParticipations: maxOrderBookParticipations,
		BatchSettlementCount:       batchSettlementCount,
		RequeueThreshold:           requeueThreshold,
	}
}

// ParamSetPairs Implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			keyMaxOrderBookParticipations,
			&p.MaxOrderBookParticipations,
			validateMaxOrderBookParticipations,
		),
		paramtypes.NewParamSetPair(
			keyBatchSettlementCount,
			&p.BatchSettlementCount,
			validateBatchSettlementCount,
		),
		paramtypes.NewParamSetPair(
			keyRequeueThreshold,
			&p.RequeueThreshold,
			validateRequeueThreshold,
		),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMaxOrderBookParticipations,
		DefaultBatchSettlementCount,
		DefaultRequeueThreshold,
	)
}

// String returns a human-readable string representation of the parameters.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// Validate validates a set of params
func (p Params) Validate() error {
	if err := validateMaxOrderBookParticipations(p.MaxOrderBookParticipations); err != nil {
		return err
	}

	if err := validateBatchSettlementCount(p.BatchSettlementCount); err != nil {
		return err
	}

	return validateRequeueThreshold(p.RequeueThreshold)
}

func validateMaxOrderBookParticipations(i interface{}) error {
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

func validateRequeueThreshold(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
