package sportevent_test

import (
	"testing"

	"github.com/sge-network/sge/testutil/nullify"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/sportevent"
	"github.com/sge-network/sge/x/sportevent/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		SportEventList: []types.SportEvent{
			{
				UID: "0",
			},
			{
				UID: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	tApp, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	sportevent.InitGenesis(ctx, tApp.SporteventKeeper, genesisState)
	got := sportevent.ExportGenesis(ctx, tApp.SporteventKeeper)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.SportEventList, got.SportEventList)
	// this line is used by starport scaffolding # genesis/test/assert
}
