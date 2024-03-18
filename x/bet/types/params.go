package types

import (
	fmt "fmt"

	"gopkg.in/yaml.v2"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	defaultBatchSettlementCount  = 1000
	defaultMaxBetByUIDQueryCount = 10
)

var (
	defaultMinAmount               = sdkmath.NewInt(1000000)
	defaultFee                     = sdkmath.NewInt(100)
	defaultPriceLockFeePercent     = sdk.NewDecWithPrec(5, 2)
	defaultMinPriceLockPoolBalance = sdkmath.NewInt(1000000000)
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

	// keyMinPriceLockPoolBalance is the minimum pool balance of the price lock.
	keyMinPriceLockPoolBalance = []byte("MinPriceLockPoolBalance")
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	batchSettlementCount uint32,
	maxBetByUIDQueryCount uint32,
	minAmount sdkmath.Int,
	fee sdkmath.Int,
	priceLockFeePercent sdk.Dec,
	minPriceLockPoolBalance sdkmath.Int,
) Params {
	return Params{
		BatchSettlementCount:  batchSettlementCount,
		MaxBetByUidQueryCount: maxBetByUIDQueryCount,
		Constraints: Constraints{
			MinAmount:           minAmount,
			Fee:                 fee,
			PriceLockFeePercent: priceLockFeePercent,
		},
		MinPriceLockPoolBalance: minPriceLockPoolBalance,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		defaultBatchSettlementCount,
		defaultMaxBetByUIDQueryCount,
		defaultMinAmount,
		defaultFee,
		defaultPriceLockFeePercent,
		defaultMinPriceLockPoolBalance,
	)
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
		paramtypes.NewParamSetPair(
			keyMinPriceLockPoolBalance,
			&p.MinPriceLockPoolBalance,
			validateMinPriceLockPoolBalance,
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

	if err := validateConstraints(p.Constraints); err != nil {
		return err
	}

	return validateMinPriceLockPoolBalance(p.MinPriceLockPoolBalance)
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

	if v.PriceLockFeePercent.LT(sdk.ZeroDec()) {
		return fmt.Errorf("minimum bet price lock fee must be positive: %s", v.PriceLockFeePercent)
	}

	return nil
}

func validateMinPriceLockPoolBalance(i interface{}) error {
	v, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LTE(sdk.OneInt()) {
		return fmt.Errorf("minimum price lock pool balance must be positive and more than one: %d", v)
	}

	return nil
}
