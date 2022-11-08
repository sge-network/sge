package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/mint/types"
	"github.com/stretchr/testify/require"
)

func TestNewGenesisState(t *testing.T) {
	defaultGs := types.DefaultGenesis()
	gs := types.NewGenesisState(defaultGs.Minter, defaultGs.Params)
	require.NoError(t, gs.Validate())
}

func TestGenesisStateValidate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
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
						{YearCoefficient: sdk.MustNewDecFromStr("0"), Inflation: sdk.MustNewDecFromStr("0")},
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
						{YearCoefficient: sdk.MustNewDecFromStr("0"), Inflation: sdk.MustNewDecFromStr("0")},
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
					Inflation:       sdk.NewDec(-10),
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
					ExcludeAmount: sdk.NewInt(-1),
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
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
