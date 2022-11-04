package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/sge-network/sge/app/params"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// param keys
var (
	// KeyEventMinBetAmount is the min bet amount param key
	KeyEventMinBetAmount = []byte("EventMinBetAmount")

	// KeyMaxBetCap is the max bet cap param key
	KeyMaxBetCap = []byte("EventMaxBetCap")

	// KeyMinBetFee is the minimum bet fee param key
	KeyMinBetFee = []byte("EventMinBetFee")

	// KeyMinBetFee is the minimum bet fee param key
	KeyMinVig = []byte("KeyMinVig")

	// KeyMinBetFee is the minimum bet fee param key
	KeyMaxVig = []byte("KeyMaxVig")

	// KeyMinBetFee is the minimum bet fee param key
	KeyMaxEventLoss = []byte("KeyMaxEventLoss")
)

// default params
const (
	// DefaultMinBetAmount is the default minimum bet amount allowed
	DefaultMinBetAmount = 1000000

	// DefaultMaxBetCap is the default max bet cap  allowed
	DefaultMaxBetCap = 1000000000000

	// DefaultMaxEventLoss is the default max loss event can have
	DefaultMaxEventLoss = 1000000000000

	// DefaultMaxVig is the default max vig can have
	DefaultMaxVig = 15

	// DefaultMinVig is the default min vig can have
	DefaultMinVig = 5
)

var (
	// DefaultMinBetFee is the default minimum bet fee amount allowed
	DefaultMinBetFee = sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(100000))
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		EventMinBetAmount: sdk.NewInt(DefaultMinBetAmount),
		EventMaxBetCap:    sdk.NewInt(DefaultMaxBetCap),
		EventMinBetFee:    DefaultMinBetFee,
		EventMaxVig:       sdk.NewDec(DefaultMaxVig),
		EventMinVig:       sdk.NewDec(DefaultMinVig),
		EventMaxLoss:      sdk.NewInt(DefaultMaxEventLoss),
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
		paramtypes.NewParamSetPair(KeyMaxBetCap, &p.EventMaxBetCap, validateMaxBetCap),
		paramtypes.NewParamSetPair(KeyMinBetFee, &p.EventMinBetFee, validateMinBetFeePercentage),
		paramtypes.NewParamSetPair(KeyMinVig, &p.EventMinVig, validateMinVig),
		paramtypes.NewParamSetPair(KeyMaxVig, &p.EventMaxVig, validateMaxVig),
		paramtypes.NewParamSetPair(KeyMaxEventLoss, &p.EventMaxLoss, validateMaxEventLoss),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateMinBetAmount(i interface{}) error {
	return nil
}

func validateMaxBetCap(i interface{}) error {
	return nil
}

func validateMinBetFeePercentage(i interface{}) error {
	return nil
}

func validateMinVig(i interface{}) error {
	return nil
}

func validateMaxVig(i interface{}) error {
	return nil
}

func validateMaxEventLoss(i interface{}) error {
	return nil
}
