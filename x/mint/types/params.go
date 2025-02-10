package types

import (
	"fmt"
	"math"
	"math/big"
	"strings"

	yaml "gopkg.in/yaml.v2"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/app/params"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// inflation constants
const (
	initialInflation       = 0
	initialPhaseProvision  = 0
	initialPhaseStep       = 0
	initialTruncatedTokens = 0
)

// phase step constants
const (
	// EndPhaseAlias is an alias for built-in end phase if
	// no valid phase available according to the current block
	EndPhaseAlias = -1
)

// block timing
const (
	yearHours         = 8766 // 8766 is coming from 365.25×24h
	YearSeconds       = 60 * 60 * yearHours
	expectedBlockTime = 5 // in seconds
	BlocksPerYear     = uint64(YearSeconds / expectedBlockTime)
)

var (
	// DefaultExcludeAmount is the default value for exclude amount
	DefaultExcludeAmount = sdkmath.NewInt(int64(0))

	// DefaultPhases is the default value for inflation phases
	DefaultPhases = []Phase{
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.229787234042553191"),
		},
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.286259541984732824"),
		},
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.150250417362270451"),
		},
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.116459627329192547"),
		},
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.088041085840058694"),
		},
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.063246661981728742"),
		},
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.040871934604904632"),
		},
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.032042723631508678"),
		},
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.019710906701708279"),
		},
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.003903708523096942"),
		},
	}
)

// NewParams creates a new Params instance
func NewParams(mintDenom string, blocksPerYear uint64, excludeAmount sdkmath.Int, phases []Phase) Params {
	return Params{
		MintDenom:     mintDenom,
		BlocksPerYear: blocksPerYear,
		ExcludeAmount: excludeAmount,
		Phases:        phases,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		MintDenom:     params.DefaultBondDenom,
		BlocksPerYear: BlocksPerYear,
		Phases:        DefaultPhases,
		ExcludeAmount: DefaultExcludeAmount,
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}

	if err := validateBlocksPerYear(p.BlocksPerYear); err != nil {
		return err
	}

	if err := validatePhases(p.Phases); err != nil {
		return err
	}

	return validateExcludeAmount(p.ExcludeAmount)
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// GetPhaseAtStep returns the phase object at certain step
func (p Params) GetPhaseAtStep(phaseStep int) Phase {
	if phaseStep == EndPhaseAlias {
		return EndPhase()
	}
	if phaseStep == 0 {
		return NonePhase()
	}
	phasesCount := len(p.Phases)
	phaseIndex := phaseStep - 1
	if phasesCount > phaseIndex {
		return p.Phases[phaseIndex]
	}
	return EndPhase()
}

// IsEndPhaseByStep checks if the phase is end phase by step
func (p Params) IsEndPhaseByStep(phaseStep int) bool {
	phase := p.GetPhaseAtStep(phaseStep)
	return IsEndPhase(phase)
}

// NonePhase returns none phase object
// none phase is the initial phase of inflation with height 1
func NonePhase() Phase {
	return Phase{Inflation: sdkmath.LegacyMustNewDecFromStr("0"), YearCoefficient: sdkmath.LegacyZeroDec()}
}

// EndPhase returns end phase which there is no phase item with remaining blocks
func EndPhase() Phase {
	maxUInt64 := new(big.Int).SetUint64(math.MaxUint64)
	return Phase{
		Inflation:       sdkmath.LegacyMustNewDecFromStr("0"),
		YearCoefficient: sdkmath.LegacyNewDecFromBigInt(maxUInt64),
	}
}

// IsEndPhase returns true if the phase is equal to end phase props
func IsEndPhase(phase Phase) bool {
	endPhase := EndPhase()
	if phase.Inflation.Equal(endPhase.Inflation) &&
		phase.YearCoefficient.Equal(endPhase.YearCoefficient) {
		return true
	}
	return false
}

// getPhaseBlocks returns the total blocks of a certain phase step
func (p Params) getPhaseBlocks(phaseStep int) sdkmath.LegacyDec {
	// get the phase year coefficient
	yearCoefficient := p.Phases[phaseStep-1].YearCoefficient

	// calculate the block provisions according to the blocks per year
	// the value decimals should be truncated because in the blocks in the blockchain,
	// are going to be created one by one so there is no portion of block to be considered
	// ex.
	//    blocks per year = 101
	//    current block = 50
	//    so the changing point of the phase is in block number 50.5 which is not applicable
	//    the 0.5 provisions will be calculated in the BlockProvisions method of Minter
	phaseBlocks := yearCoefficient.Mul(sdkmath.LegacyMustNewDecFromStr(cast.ToString(p.BlocksPerYear))).TruncateDec()

	return phaseBlocks
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf(ErrTextInvalidParamType, i)
	}

	if strings.TrimSpace(v) == "" {
		return ErrMintDenomIsBlank
	}
	return sdk.ValidateDenom(v)
}

func validateExcludeAmount(i interface{}) error {
	v, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf(ErrTextInvalidParamType, i)
	}

	if v.LT(sdkmath.ZeroInt()) {
		return fmt.Errorf(ErrTextExcludeAmountMustBePositive, v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf(ErrTextInvalidParamType, i)
	}

	if v <= 0 {
		return fmt.Errorf(ErrTextBlocksPerYearMustBePositive, v)
	}

	return nil
}

func validatePhases(i interface{}) error {
	v, ok := i.([]Phase)
	if !ok {
		return fmt.Errorf(ErrTextInvalidParamType, i)
	}

	if len(v) == 0 {
		return fmt.Errorf(ErrTextPhasesShouldHaveValue, v)
	}

	for _, p := range v {
		if !p.YearCoefficient.GT(sdkmath.LegacyZeroDec()) {
			return fmt.Errorf(ErrTextYearCoefficientMustBePositive, p.YearCoefficient)
		}
		if IsEndPhase(p) {
			return fmt.Errorf(ErrTextEndPhaseParamNotAllowed, p.Inflation)
		}
	}

	return nil
}
