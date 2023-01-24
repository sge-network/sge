package types_test

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestParamsYML(t *testing.T) {
	param := types.Params{
		MintDenom:     params.DefaultBondDenom,
		BlocksPerYear: 10,
		Phases: []types.Phase{
			{Inflation: sdk.NewDec(10), YearCoefficient: sdk.NewDec(1)},
		},
		ExcludeAmount: sdk.NewInt(100),
	}

	ymlStr := param.String()
	require.Equal(t, "mintdenom: usge\nblocks_per_year: 10\nphases:\n- inflation: \"10.000000000000000000\"\n  year_coefficient: \"1.000000000000000000\"\nexclude_amount: \"100\"\n", ymlStr)
}

func TestIsPhaseAtStep(t *testing.T) {
	defaultParams := types.DefaultParams()
	p := types.NewParams(defaultParams.MintDenom, defaultParams.BlocksPerYear, defaultParams.ExcludeAmount, defaultParams.Phases)
	firstDefaultStep := p.GetPhaseAtStep(1)
	require.Equal(t, defaultParams.Phases[0], firstDefaultStep)
}

func TestIsPhaseEndPhase(t *testing.T) {
	defaultParams := types.DefaultParams()
	p := types.NewParams(defaultParams.MintDenom, defaultParams.BlocksPerYear, defaultParams.ExcludeAmount, defaultParams.Phases)
	require.False(t, p.IsEndPhaseByStep(1))
	require.False(t, p.IsEndPhaseByStep(2))
	require.False(t, p.IsEndPhaseByStep(3))
	require.False(t, p.IsEndPhaseByStep(4))
	require.False(t, p.IsEndPhaseByStep(5))
	require.False(t, p.IsEndPhaseByStep(6))
	require.False(t, p.IsEndPhaseByStep(7))
	require.False(t, p.IsEndPhaseByStep(8))
	require.False(t, p.IsEndPhaseByStep(9))
	require.False(t, p.IsEndPhaseByStep(10))
	require.True(t, p.IsEndPhaseByStep(11))
}

func TestParamKeyTable(t *testing.T) {
	types.ParamKeyTable()
}

func TestParamSetPairs(t *testing.T) {
	params := types.DefaultParams()
	patamSetPairs := params.ParamSetPairs()

	pt := reflect.TypeOf(params)
	paramFiledNames := make([]string, pt.NumField())
	for i := range paramFiledNames {
		paramFiledNames[i] = pt.Field(i).Name
	}

	var paramSetFiledNames []string
	for _, v := range patamSetPairs {
		paramSetFiledNames = append(paramSetFiledNames, string(v.Key))
	}

	require.True(t, reflect.DeepEqual(paramFiledNames, paramSetFiledNames))
}

func TestNonPhase(t *testing.T) {
	params := types.DefaultParams()

	phase := params.GetPhaseAtStep(0)
	require.Equal(t, sdk.NewDec(0), phase.YearCoefficient)
}
