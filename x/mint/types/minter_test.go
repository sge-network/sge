package types_test

import (
	math "math"
	//#nosec
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/mint/types"
)

func TestPhaseInflation(t *testing.T) {
	params := types.DefaultParams()
	// Governing Mechanism:
	// sge tokenomics

	tests := []struct {
		phase        int
		expInflation sdkmath.LegacyDec
	}{
		// phase 1
		{1, sdkmath.LegacyMustNewDecFromStr("0.229787234042553191")},
		// phase 2
		{2, sdkmath.LegacyMustNewDecFromStr("0.286259541984732824")},
		// phase 3
		{3, sdkmath.LegacyMustNewDecFromStr("0.150250417362270451")},
		// phase 4
		{4, sdkmath.LegacyMustNewDecFromStr("0.116459627329192547")},
		// phase 5
		{5, sdkmath.LegacyMustNewDecFromStr("0.088041085840058694")},
		// phase 6
		{6, sdkmath.LegacyMustNewDecFromStr("0.063246661981728742")},
		// phase 7
		{7, sdkmath.LegacyMustNewDecFromStr("0.040871934604904632")},
		// phase 8
		{8, sdkmath.LegacyMustNewDecFromStr("0.032042723631508678")},
		// phase 9
		{9, sdkmath.LegacyMustNewDecFromStr("0.019710906701708279")},
		// phase 10
		{10, sdkmath.LegacyMustNewDecFromStr("0.003903708523096942")},
		// end phase, inflation: 0%
		{11, sdkmath.LegacyMustNewDecFromStr("0")},
		// end phase, inflation: 0%
		{13, sdkmath.LegacyMustNewDecFromStr("0")},
		// end phase, inflation: 0%
		{23, sdkmath.LegacyMustNewDecFromStr("0")},
		// end phase, inflation: 0%
		{-1, sdkmath.LegacyMustNewDecFromStr("0")},
	}
	for i, tc := range tests {
		inflation := params.GetPhaseAtStep(tc.phase).Inflation

		require.True(t, inflation.Equal(tc.expInflation),
			"Test Index: %v\nInflation:  %v\nExpected: %v\n", i, inflation, tc.expInflation)
	}
}

func TestNextPhase(t *testing.T) {
	minter := types.DefaultInitialMinter()
	params := types.DefaultParams()

	blocksPerYear := uint64(100)
	tests := []struct {
		currentBlock, currentPhase int64
		blocksYear                 uint64
		expPhase                   int
	}{
		{1, 0, blocksPerYear, 1},
		{20, 1, blocksPerYear, 1},
		{49, 1, blocksPerYear, 1},
		{50, 1, blocksPerYear, 1},
		{51, 1, blocksPerYear, 2},
		{200, 4, blocksPerYear, 4},
		{201, 4, blocksPerYear, 5},
		{300, 6, blocksPerYear, 6},
		{301, 6, blocksPerYear, 7},
		{400, 7, blocksPerYear, 8},
		{500, 8, blocksPerYear, 10},
		{601, 10, blocksPerYear, types.EndPhaseAlias},
		// increase blocks_per_year
		{101, 1, blocksPerYear * 2, 2},
		{201, 2, blocksPerYear * 2, 3},
		{299, 3, blocksPerYear * 2, 3},
		{300, 3, blocksPerYear * 2, 3},
		{301, 3, blocksPerYear * 2, 4},
		{401, 4, blocksPerYear * 2, 5},
		{1001, 10, blocksPerYear * 2, types.EndPhaseAlias},
	}
	for i, tc := range tests {
		params.BlocksPerYear = tc.blocksYear

		_, phaseStep := minter.CurrentPhase(params, tc.currentBlock)

		require.True(t, phaseStep == tc.expPhase,
			"Test Index: %v\nPhase Step:  %v\nExpected: %v\n", i, phaseStep, tc.expPhase)
	}
}

func TestNextPhaseAfterAppendToEndPhase(t *testing.T) {
	minter := types.DefaultInitialMinter()

	params := types.DefaultParams()
	params.BlocksPerYear = 100

	minter.PhaseStep = int32(len(params.Phases) + 1)

	params.Phases = append(
		params.Phases,
		types.Phase{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("1"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.032042723631508678"),
		},
	)
	phase, _ := minter.CurrentPhase(params, 351)
	require.Equal(t, params.Phases[len(params.Phases)-1].Inflation, phase.Inflation)

	require.False(t, types.IsEndPhase(phase))
}

func TestNextPhaseAfterRelaceEndPhase(t *testing.T) {
	minter := types.DefaultInitialMinter()
	params := types.DefaultParams()

	params.BlocksPerYear = 100
	params.Phases = []types.Phase{
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("1"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.00000000000000000"),
		},
		types.EndPhase(),
	}

	minter.PhaseStep = 2

	params.Phases = []types.Phase{
		{
			YearCoefficient: sdkmath.LegacyMustNewDecFromStr("2.5"),
			Inflation:       sdkmath.LegacyMustNewDecFromStr("0.00000000000000000"),
		},
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
	phase, _ := minter.CurrentPhase(params, 250)
	require.Equal(t, params.Phases[0].Inflation, phase.Inflation)
	require.False(t, types.IsEndPhase(phase))
}

func TestBlockProvision(t *testing.T) {
	minter := types.InitialMinter(sdkmath.LegacyNewDecWithPrec(1, 1))
	params := types.DefaultParams()

	tests := []struct {
		phaseProvisions int64
		expProvisions   int64
	}{
		{types.YearSeconds / 5, 2},
		{types.YearSeconds/5 + 1, 2},
		{(types.YearSeconds / 5) * 2, 4},
		{(types.YearSeconds / 5) / 2, 1},
	}

	for i, tc := range tests {
		minter.PhaseProvisions = sdkmath.LegacyNewDec(tc.phaseProvisions)
		provisions, _ := minter.BlockProvisions(params, 1)

		expProvisions := sdk.NewCoin(params.MintDenom,
			sdkmath.NewInt(tc.expProvisions))

		require.True(t, expProvisions.IsEqual(provisions),
			"test: %v\n\tExp: %v\n\tGot: %v\n",
			i, tc.expProvisions, provisions)
	}
}

func TestMaxComulativeBlocksPerYear(t *testing.T) {
	minter := types.DefaultInitialMinter()
	params := types.DefaultParams()
	params.BlocksPerYear = math.MaxInt64
	// if comulative overflows, this line will panic
	minter.CurrentPhase(params, math.MaxInt64)

	params.BlocksPerYear = 1000
	_, phaseStep := minter.CurrentPhase(params, math.MaxInt64)
	require.Equal(t, types.EndPhaseAlias, phaseStep)
}

func TestBlockProvisions(t *testing.T) {
	minter := types.DefaultInitialMinter()
	params := types.DefaultParams()

	blocksPerYear := uint64(51480)
	totalSupply := sdkmath.NewIntFromUint64(1150000000000000)
	params.ExcludeAmount = sdkmath.NewInt(445000000000025)
	tests := []struct {
		phaseindex   int
		expProvision int64
	}{
		{0, 81000000000000},
		{1, 112500000000000},
		{2, 67500000000000},
		{3, 56250000000000},
		{4, 45000000000000},
		{5, 33750000000000},
		{6, 22500000000000},
		{7, 18000000000000},
		{8, 11250000000000},
		{9, 2250000000000},
		{10, 0},
	}
	currentBlock := int64(1)
	currentPhaseStep := 1
	for i, tc := range tests {
		phaseStep := tc.phaseindex + 1
		var phase types.Phase
		if len(params.Phases) > i {
			phase = params.Phases[tc.phaseindex]
		} else {
			phase = types.EndPhase()
		}

		minter.Inflation = phase.Inflation
		params.BlocksPerYear = blocksPerYear
		minter.PhaseProvisions = minter.NextPhaseProvisions(totalSupply, params.ExcludeAmount, phase)

		currentPhaseProvision := sdkmath.NewInt(0)
		if !types.IsEndPhase(phase) {
			for phaseStep == currentPhaseStep {
				blockIProvision, truncatedToken := minter.BlockProvisions(params, phaseStep)
				minter.TruncatedTokens = truncatedToken

				totalSupply = totalSupply.Add(blockIProvision.Amount)
				currentPhaseProvision = currentPhaseProvision.Add(blockIProvision.Amount)

				currentBlock++
				_, currentPhaseStep = minter.CurrentPhase(params, currentBlock)
			}
			t.Logf("current block: %v", currentBlock)
			t.Logf("Total Supply after phase %v: %v", phaseStep, totalSupply)
		}

		phaseDeviation := sdkmath.NewInt(tc.expProvision).Sub(currentPhaseProvision).Abs().Int64()
		allowedPhaseDeviation := int64(4)
		require.LessOrEqual(
			t,
			phaseDeviation,
			allowedPhaseDeviation,
			"Test Index: %v\nPhaseProvisions: %v\nPhaseProvision: %v\nExpected: %v\n",
			i,
			minter.PhaseProvisions,
			currentPhaseProvision,
			tc.expProvision,
		)
	}

	expTotalSupply := sdkmath.NewInt(1600000000000000)
	deviation := totalSupply.Sub(expTotalSupply).Abs().Int64()
	allowedTotalSupplyDeviation := int64(17)
	require.LessOrEqual(t, deviation, allowedTotalSupplyDeviation,
		"TotalSupply: %v\nExpected: %v\n", totalSupply, expTotalSupply)
}

func TestAnnualProvisions(t *testing.T) {
	minter := types.DefaultInitialMinter()
	params := types.DefaultParams()

	totalSupply := int64(150)
	excludeAmount := int64(50)
	firstPhase := params.Phases[0]
	minter.PhaseStep = 1
	minter.Inflation = firstPhase.Inflation
	minter.PhaseProvisions = minter.NextPhaseProvisions(
		sdkmath.NewInt(totalSupply),
		sdkmath.NewInt(excludeAmount),
		firstPhase,
	)
	t.Log(minter.PhaseProvisions)
	annualProvisions := minter.AnnualProvisions(firstPhase)

	require.Equal(t, firstPhase.Inflation.Mul(sdkmath.LegacyNewDec(totalSupply-excludeAmount)), annualProvisions)
}

// Benchmarking :)
// previously using sdkmath.Int operations:
// BenchmarkBlockProvision-4 5000000 220 ns/op
//
// using sdkmath.LegacyDec operations: (current implementation)
// BenchmarkBlockProvision-4 3000000 429 ns/op
func BenchmarkBlockProvision(b *testing.B) {
	b.ReportAllocs()
	minter := types.InitialMinter(sdkmath.LegacyNewDecWithPrec(1, 1))
	params := types.DefaultParams()

	s1 := rand.NewSource(100)
	//#nosec
	r1 := rand.New(s1)
	minter.PhaseProvisions = sdkmath.LegacyNewDec(r1.Int63n(1000000))

	// run the BlockProvision function b.N times
	for n := 0; n < b.N; n++ {
		minter.BlockProvisions(params, 1)
	}
}

// Next inflation benchmarking
// BenchmarkPhaseInflation-4 1000000 1828 ns/op
func BenchmarkPhaseInflation(b *testing.B) {
	params := types.DefaultParams()
	b.ReportAllocs()
	phase := 4

	// run the PhaseInflationRate function b.N times
	for n := 0; n < b.N; n++ {
		params.GetPhaseAtStep(phase)
	}
}

// Next phase provisions benchmarking
// BenchmarkNextPhaseProvisions-4 5000000 251 ns/op
func BenchmarkNextPhaseProvisions(b *testing.B) {
	b.ReportAllocs()
	minter := types.InitialMinter(sdkmath.LegacyNewDecWithPrec(1, 1))
	params := types.DefaultParams()
	totalSupply := sdkmath.NewInt(100000000000000)
	phase := params.GetPhaseAtStep(1)

	// run the NextPhaseProvisions function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextPhaseProvisions(totalSupply, types.DefaultExcludeAmount, phase)
	}
}
