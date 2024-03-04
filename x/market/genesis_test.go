package market_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/testutil/simapp"
	market "github.com/sge-network/sge/x/market"
	"github.com/sge-network/sge/x/market/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		MarketList: []types.Market{
			{
				UID: "0",
			},
			{
				UID: "1",
			},
		},
	}

	tApp, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)

	market.InitGenesis(ctx, *tApp.MarketKeeper, genesisState)
	got := market.ExportGenesis(ctx, *tApp.MarketKeeper)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.MarketList, got.MarketList)
}
