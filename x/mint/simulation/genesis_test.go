package simulation_test

import (
	"encoding/json"
	//#nosec
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sge-network/sge/x/mint/simulation"
	"github.com/sge-network/sge/x/mint/types"
)

var params = types.Params{
	MintDenom:     "usge",
	BlocksPerYear: 100,
	Phases:        types.DefaultGenesis().Params.Phases,
}

// TestRandomizedGenState tests the normal scenario of applying RandomizedGenState.
// Abonormal scenarios are not tested here.
func TestRandomizedGenState(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	s := rand.NewSource(1)
	//#nosec
	r := rand.New(s)

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: 1000,
		GenState:     make(map[string]json.RawMessage),
	}

	simulation.RandomizedGenState(&simState)

	var mintGenesis types.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &mintGenesis)

	require.Equal(t, int64(6311520), mintGenesis.Params.BlocksPerYear)
	require.Equal(t, "usge", mintGenesis.Params.MintDenom)
	bp, _ := mintGenesis.Minter.BlockProvisions(mintGenesis.Params, 1)
	require.Equal(t, "0usge", bp.String())
	require.Equal(t, "0.000000000000000000", mintGenesis.Minter.NextPhaseProvisions(sdk.OneInt(), types.DefaultExcludeAmount, types.NonePhase()).String())
	require.Equal(t, "0.229787234042553191", params.GetPhaseAtStep(1).Inflation.String())
	require.Equal(t, "0.286259541984732824", params.GetPhaseAtStep(2).Inflation.String())
	require.Equal(t, "0.150250417362270451", params.GetPhaseAtStep(3).Inflation.String())
	require.Equal(t, "0.116459627329192547", params.GetPhaseAtStep(4).Inflation.String())
	require.Equal(t, "0.088041085840058694", params.GetPhaseAtStep(5).Inflation.String())
	require.Equal(t, "0.063246661981728742", params.GetPhaseAtStep(6).Inflation.String())
	require.Equal(t, "0.040871934604904632", params.GetPhaseAtStep(7).Inflation.String())
	require.Equal(t, "0.032042723631508678", params.GetPhaseAtStep(8).Inflation.String())
	require.Equal(t, "0.019710906701708279", params.GetPhaseAtStep(9).Inflation.String())
	require.Equal(t, "0.003903708523096942", params.GetPhaseAtStep(10).Inflation.String())
	require.Equal(t, "0.000000000000000000", params.GetPhaseAtStep(types.EndPhaseAlias).Inflation.String())
	require.Equal(t, "0.170000000000000000", mintGenesis.Minter.Inflation.String())
	_, nextPhaseIndex := mintGenesis.Minter.CurrentPhase(mintGenesis.Params, 1)
	require.Equal(t, 1, nextPhaseIndex)
	require.Equal(t, int32(0), mintGenesis.Minter.PhaseStep)
	require.Equal(t, "0.000000000000000000", mintGenesis.Minter.PhaseProvisions.String())
}

// TestRandomizedGenState tests abnormal scenarios of applying RandomizedGenState.
func TestRandomizedGenState1(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	s := rand.NewSource(1)
	//#nosec
	r := rand.New(s)
	// all these tests will panic
	tests := []struct {
		simState module.SimulationState
		panicMsg string
	}{
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{}, "invalid memory address or nil pointer dereference"},
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{
				AppParams: make(simtypes.AppParams),
				Cdc:       cdc,
				Rand:      r,
			}, "assignment to entry in nil map"},
	}

	for _, tt := range tests {
		require.Panicsf(t, func() { simulation.RandomizedGenState(&tt.simState) }, tt.panicMsg)
	}
}
