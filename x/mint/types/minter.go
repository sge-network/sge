package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation and phase
// provisions values.
func NewMinter(inflation, phaseProvisions sdk.Dec, phaseStep int32, truncatedTokens sdk.Dec) Minter {
	return Minter{
		Inflation:       inflation,
		PhaseProvisions: phaseProvisions,
		PhaseStep:       phaseStep,
		TruncatedTokens: truncatedTokens,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(inflation sdk.Dec) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(initialPhaseProvision),
		initialPhaseStep,
		sdk.NewDec(initialTruncatedTokens),
	)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 10%.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDec(initialInflation),
	)
}

// ValidateMinter validates minter
func ValidateMinter(minter Minter) error {
	if minter.Inflation.IsNegative() {
		return fmt.Errorf(ErrTextMintParamInflationShouldBePositive,
			minter.Inflation.String())
	}
	return nil
}

// CurrentPhase returns current phase of the inflation
func (m Minter) CurrentPhase(params Params, currentBlock int64) (Phase, int) {
	if currentBlock == 1 {
		return params.GetPhaseAtStep(1), 1
	}

	cumulativeBlock := sdk.NewDec(0)
	var currentStep int
	var found bool

	// add each phase blocks until reaching the range which the current block is in
	for i := 0; i < len(params.Phases); i++ {
		// add current phase blocks to cummulative blocks
		cumulativeBlock = cumulativeBlock.Add(params.getPhaseBlocks(i + 1))

		currentStep = i + 1

		// if the current block is less than or equal to cummulative blocks
		// this means that we are in the i+1 step which is set in above line
		if sdk.NewDec(currentBlock).LTE(cumulativeBlock) {
			found = true
			// it is the current phase
			// so there is no need for furthur phase blocks check
			break
		}
	}

	// if there is no detected phase,
	// this means that the rest of inflation is zero as end phase
	if !found {
		return EndPhase(), EndPhaseAlias
	}

	// the phase has found and we need to return the phase specifications
	return params.GetPhaseAtStep(currentStep), currentStep
}

// NextPhaseProvisions returns the phase provisions based on current total
// supply and inflation rate.
func (m Minter) NextPhaseProvisions(totalSupply sdk.Int, excludeAmount sdk.Int, phase Phase) sdk.Dec {
	// calculate annual provisions as normal
	annualProvisions := m.Inflation.MulInt(totalSupply.Sub(excludeAmount))

	// return this phase provisions according to the year coefficient
	// ex.
	//    year coefficient = 0.5
	//    blocks per year = 100
	// 	  this phase provisions is 100 * 0.5 => 50
	return annualProvisions.Mul(phase.YearCoefficient)
}

// BlockProvisions returns the provisions for a block based on the phase
// provisions rate.
func (m Minter) BlockProvisions(params Params, phaseStep int) (sdk.Coin, sdk.Dec) {
	// get total blocks in this phase
	blocksPerPhase := params.getPhaseBlocks(phaseStep).TruncateDec()

	// detect each block provisions then and the truncated value from previous block
	provisionAmt := m.PhaseProvisions.Quo(blocksPerPhase).Add(m.TruncatedTokens)

	// extract the integer and decimal part of provisions
	// the decimal part is the truncated value because of conversion to sdk.Int
	// so the decimal part is truncated and needs to be added in next block
	intPart := provisionAmt.TruncateDec()
	decPart := provisionAmt.Sub(intPart)

	return sdk.NewCoin(params.MintDenom, intPart.TruncateInt()), decPart
}

// AnnualProvisions returns annual provisions for the phase.
func (m Minter) AnnualProvisions(phase Phase) sdk.Dec {
	return m.PhaseProvisions.Quo(phase.YearCoefficient)
}
