package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	"github.com/sge-network/sge/x/mint/types"
)

func TestNewGenesisState(t *testing.T) {
	defaultGs := types.DefaultGenesisState()
	gs := types.NewGenesisState(defaultGs.Minter, defaultGs.Params)
	require.NoError(t, types.ValidateGenesis(*gs))
}

func TestGenesisStateValidate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesisState(),
			valid:    true,
		},
		{
			desc: "invalid denom",
			genState: types.NewGenesisState(
				types.DefaultInitialMinter(),
				types.Params{
					MintDenom:     "1invalid_denom",
					BlocksPerYear: types.DefaultParams().BlocksPerYear,
					Phases:        types.DefaultParams().Phases,
					ExcludeAmount: types.DefaultParams().ExcludeAmount,
				},
			),
			valid: false,
		},
		{
			desc: "empty denom",
			genState: types.NewGenesisState(
				types.DefaultInitialMinter(),
				types.Params{
					MintDenom:     "",
					BlocksPerYear: types.DefaultParams().BlocksPerYear,
					Phases:        types.DefaultParams().Phases,
					ExcludeAmount: types.DefaultParams().ExcludeAmount,
				},
			),
			valid: false,
		},
		{
			desc: "zero blocks per year",
			genState: types.NewGenesisState(
				types.DefaultInitialMinter(),
				types.Params{
					MintDenom:     types.DefaultParams().MintDenom,
					BlocksPerYear: 0,
					Phases:        types.DefaultParams().Phases,
					ExcludeAmount: types.DefaultParams().ExcludeAmount,
				},
			),
			valid: false,
		},
		{
			desc: "phase with zero coefficient",
			genState: types.NewGenesisState(
				types.DefaultInitialMinter(),
				types.Params{
					MintDenom:     types.DefaultParams().MintDenom,
					BlocksPerYear: types.DefaultParams().BlocksPerYear,
					Phases: []types.Phase{
						{YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0"), Inflation: sdkmath.LegacyMustNewDecFromStr("0")},
					},
					ExcludeAmount: types.DefaultParams().ExcludeAmount,
				},
			),
			valid: false,
		},
		{
			desc: "phase with empty phases",
			genState: types.NewGenesisState(
				types.DefaultInitialMinter(),
				types.Params{
					MintDenom:     types.DefaultParams().MintDenom,
					BlocksPerYear: types.DefaultParams().BlocksPerYear,
					Phases:        []types.Phase{},
					ExcludeAmount: types.DefaultParams().ExcludeAmount,
				},
			),
			valid: false,
		},
		{
			desc: "phase with zero coefficient",
			genState: types.NewGenesisState(
				types.DefaultInitialMinter(),
				types.Params{
					MintDenom:     types.DefaultParams().MintDenom,
					BlocksPerYear: types.DefaultParams().BlocksPerYear,
					Phases: []types.Phase{
						{YearCoefficient: sdkmath.LegacyMustNewDecFromStr("0"), Inflation: sdkmath.LegacyMustNewDecFromStr("0")},
					},
					ExcludeAmount: types.DefaultParams().ExcludeAmount,
				},
			),
			valid: false,
		},
		{
			desc: "phase is equal to end-phase",
			genState: types.NewGenesisState(
				types.DefaultInitialMinter(),
				types.Params{
					MintDenom:     types.DefaultParams().MintDenom,
					BlocksPerYear: types.DefaultParams().BlocksPerYear,
					Phases:        []types.Phase{types.EndPhase()},
					ExcludeAmount: types.DefaultParams().ExcludeAmount,
				},
			),
			valid: false,
		},
		{
			desc: "negative inflation",
			genState: types.NewGenesisState(
				types.Minter{
					Inflation:       sdkmath.LegacyNewDec(-10),
					PhaseStep:       types.DefaultInitialMinter().PhaseStep,
					PhaseProvisions: types.DefaultInitialMinter().PhaseProvisions,
				},
				types.DefaultParams(),
			),
			valid: false,
		},
		{
			desc: "negative exclude amount",
			genState: types.NewGenesisState(
				types.DefaultInitialMinter(),
				types.Params{
					MintDenom:     types.DefaultParams().MintDenom,
					BlocksPerYear: types.DefaultParams().BlocksPerYear,
					Phases:        types.DefaultParams().Phases,
					ExcludeAmount: sdkmath.NewInt(-1),
				},
			),
			valid: false,
		},
		{
			desc: "valid genesis state",
			genState: types.NewGenesisState(
				types.DefaultInitialMinter(),
				types.DefaultParams(),
			),
			valid: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := types.ValidateGenesis(*tc.genState)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
