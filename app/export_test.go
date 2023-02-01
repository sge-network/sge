package app_test

import (
	"testing"

	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestExport(t *testing.T) {
	tObjs, _, err := simappUtil.GetTestObjectsWithOptions(
		simappUtil.Options{
			CreateGenesisValidators: false,
		},
	)
	require.NoError(t, err)

	exp, err := tObjs.ExportAppStateAndValidators(true, []string{})
	require.NoError(t, err)

	_, err = exp.AppState.MarshalJSON()
	require.NoError(t, err)

	require.Equal(t, []tmtypes.GenesisValidator(nil), exp.Validators)
	require.Equal(t, int64(0), exp.Height)

	exp, err = tObjs.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err)
}
