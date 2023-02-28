package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// params keys
var (
	// KeyCommitteeMembers defines the key for committee members
	KeyCommitteeMembers = []byte("CommitteeMembers")
)

// default values
var (
	// DefaultCommitteeMembers defines the default value of committee members
	DefaultCommitteeMembers = []string(nil)

	// InitialSrPool defines the value of the locked amount
	// and the unlocked amount in the SR_Pool account initially
	// when the chain is started
	InitialSrPool = SRPool{
		LockedAmount:   sdk.ZeroInt(),
		UnlockedAmount: sdk.NewInt(150000000000000),
	}
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		CommitteeMembers: DefaultCommitteeMembers,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyCommitteeMembers, &p.CommitteeMembers, validateCommitteeMembers),
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

// String implements the Stringer interface.
func (sp SRPool) String() string {
	out, err := yaml.Marshal(sp)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// String implements the Stringer interface.
func (r Reserver) String() string {
	out, err := yaml.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func validateCommitteeMembers(i interface{}) error {
	_, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
