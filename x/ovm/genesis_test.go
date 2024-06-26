package ovm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/testutil/simapp"

	"github.com/sge-network/sge/x/ovm"
	"github.com/sge-network/sge/x/ovm/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		KeyVault: types.KeyVault{
			PublicKeys: []string{"Key1"},
		},
	}

	tApp, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)

	ovm.InitGenesis(ctx, *tApp.OVMKeeper, genesisState)
	got := ovm.ExportGenesis(ctx, *tApp.OVMKeeper)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
