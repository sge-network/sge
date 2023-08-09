package types_test

import (
	math "math"
	//#nosec
	"math/rand"
	"testing"

	"github.com/sge-network/sge/x/mint/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestPhaseInflation(t *testing.T) {
	params := types.DefaultParams()
	// Governing Mechanism:
	// sge tokenomics

	tests := []struct {
		phase        int
		expInflation sdk.Dec
	}{
		// phase 1
		{1, sdk.MustNewDecFromStr("0.229787234042553191")},
		// phase 2
		{2, sdk.MustNewDecFromStr("0.286259541984732824")},
		// phase 3
		{3, sdk.MustNewDecFromStr("0.150250417362270451")},
		// phase 4
		{4, sdk.MustNewDecFromStr("0.116459627329192547")},
		// phase 5
		{5, sdk.MustNewDecFromStr("0.088041085840058694")},
		// phase 6
		{6, sdk.MustNewDecFromStr("0.063246661981728742")},
		// phase 7
		{7, sdk.MustNewDecFromStr("0.040871934604904632")},
		// phase 8
		{8, sdk.MustNewDecFromStr("0.032042723631508678")},
		// phase 9
		{9, sdk.MustNewDecFromStr("0.019710906701708279")},
		// phase 10
		{10, sdk.MustNewDecFromStr("0.003903708523096942")},
		// end phase, inflation: 0%
		{11, sdk.MustNewDecFromStr("0")},
		// end phase, inflation: 0%
		{13, sdk.MustNewDecFromStr("0")},
		// end phase, inflation: 0%
		{23, sdk.MustNewDecFromStr("0")},
		// end phase, inflation: 0%
		{-1, sdk.MustNewDecFromStr("0")},
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

	blocksPerYear := int64(100)
	tests := []struct {
		currentBlock, currentPhase int64
		blocksYear                 int64
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
			YearCoefficient: sdk.MustNewDecFromStr("1"),
			Inflation:       sdk.MustNewDecFromStr("0.032042723631508678"),
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
			YearCoefficient: sdk.MustNewDecFromStr("1"),
			Inflation:       sdk.MustNewDecFromStr("0.00000000000000000"),
		},
		types.EndPhase(),
	}

	minter.PhaseStep = 2

	params.Phases = []types.Phase{
		{
			YearCoefficient: sdk.MustNewDecFromStr("2.5"),
			Inflation:       sdk.MustNewDecFromStr("0.00000000000000000"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.229787234042553191"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.286259541984732824"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.150250417362270451"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.116459627329192547"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.088041085840058694"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.063246661981728742"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.040871934604904632"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.032042723631508678"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.019710906701708279"),
		},
		{
			YearCoefficient: sdk.MustNewDecFromStr("0.5"),
			Inflation:       sdk.MustNewDecFromStr("0.003903708523096942"),
		},
	}
	phase, _ := minter.CurrentPhase(params, 250)
	require.Equal(t, params.Phases[0].Inflation, phase.Inflation)
	require.False(t, types.IsEndPhase(phase))
}

func TestBlockProvision(t *testing.T) {
	minter := types.InitialMinter(sdk.NewDecWithPrec(1, 1))
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
		minter.PhaseProvisions = sdk.NewDec(tc.phaseProvisions)
		provisions, _ := minter.BlockProvisions(params, 1)

		expProvisions := sdk.NewCoin(params.MintDenom,
			sdk.NewInt(tc.expProvisions))

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

	blocksPerYear := int64(51480)
	totalSupply := sdk.NewIntFromUint64(1150000000000000)
	params.ExcludeAmount = sdk.NewInt(445000000000025)
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

		currentPhaseProvision := sdk.NewInt(0)
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

		phaseDeviation := sdk.NewInt(tc.expProvision).Sub(currentPhaseProvision).Abs().Int64()
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

	expTotalSupply := sdk.NewInt(1600000000000000)
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
		sdk.NewInt(totalSupply),
		sdk.NewInt(excludeAmount),
		firstPhase,
	)
	t.Log(minter.PhaseProvisions)
	annualProvisions := minter.AnnualProvisions(firstPhase)

	require.Equal(t, firstPhase.Inflation.Mul(sdk.NewDec(totalSupply-excludeAmount)), annualProvisions)
}

// Benchmarking :)
// previously using sdkmath.Int operations:
// BenchmarkBlockProvision-4 5000000 220 ns/op
//
// using sdk.Dec operations: (current implementation)
// BenchmarkBlockProvision-4 3000000 429 ns/op
func BenchmarkBlockProvision(b *testing.B) {
	b.ReportAllocs()
	minter := types.InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := types.DefaultParams()

	s1 := rand.NewSource(100)
	//#nosec
	r1 := rand.New(s1)
	minter.PhaseProvisions = sdk.NewDec(r1.Int63n(1000000))

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
	minter := types.InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := types.DefaultParams()
	totalSupply := sdk.NewInt(100000000000000)
	phase := params.GetPhaseAtStep(1)

	// run the NextPhaseProvisions function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextPhaseProvisions(totalSupply, types.DefaultExcludeAmount, phase)
	}
}
