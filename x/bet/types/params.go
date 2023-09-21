package types

import (
	fmt "fmt"

	"gopkg.in/yaml.v2"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	batchSettlementCount  = 1000
	maxBetByUIDQueryCount = 10
)

var (
	defaultMinAmount = sdkmath.NewInt(1000000)
	defaultFee       = sdkmath.NewInt(100)
)

// parameter store keys
var (
	// keyBatchSettlementCount is the batch settlement
	// count of bets
	keyBatchSettlementCount = []byte("BatchSettlementCount")

	// keyMaxBetByUIDQueryCount is the max count of
	// the queryable bets by UID list.
	keyMaxBetByUIDQueryCount = []byte("MaxBetByUidQueryCount")

	// keyWagerConstraints is the default bet placement
	// constraints.
	keyWagerConstraints = []byte("WagerConstraints")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		BatchSettlementCount:  batchSettlementCount,
		MaxBetByUidQueryCount: maxBetByUIDQueryCount,
		Constraints: Constraints{
			MinAmount: defaultMinAmount,
			Fee:       defaultFee,
		},
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(
			keyBatchSettlementCount,
			&p.BatchSettlementCount,
			validateBatchSettlementCount,
		),
		paramtypes.NewParamSetPair(
			keyMaxBetByUIDQueryCount,
			&p.MaxBetByUidQueryCount,
			validateMaxBetByUIDQueryCount,
		),
		paramtypes.NewParamSetPair(
			keyWagerConstraints,
			&p.Constraints,
			validateConstraints,
		),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateBatchSettlementCount(p.BatchSettlementCount); err != nil {
		return err
	}

	if err := validateMaxBetByUIDQueryCount(p.MaxBetByUidQueryCount); err != nil {
		return err
	}

	return validateConstraints(p.Constraints)
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
		return fmt.Errorf("%s: %T", ErrTextInvalidParamType, i)
	}

	if v <= 0 {
		return fmt.Errorf("%s: %d", ErrTextBatchSettlementCountMustBePositive, v)
	}

	return nil
}

func validateMaxBetByUIDQueryCount(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("%s: %T", ErrTextInvalidParamType, i)
	}

	if v <= 0 {
		return fmt.Errorf("%s: %d", ErrTextMaxBetUIDQueryCountMustBePositive, v)
	}

	return nil
}

func validateConstraints(i interface{}) error {
	v, ok := i.(Constraints)
	if !ok {
		return fmt.Errorf("%s: %T", ErrTextInvalidParamType, i)
	}

	if v.MinAmount.LTE(sdk.OneInt()) {
		return fmt.Errorf("minimum bet amount must be more than one: %d", v.MinAmount.Int64())
	}

	if v.Fee.LT(sdk.ZeroInt()) {
		return fmt.Errorf("minimum bet fee must be positive: %d", v.Fee.Int64())
	}

	return nil
}
