package bet_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/testutil/simapp"

	"github.com/sge-network/sge/x/bet"
	"github.com/sge-network/sge/x/bet/types"
)

func TestGenesis(t *testing.T) {
	tApp, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		BetList: []types.Bet{
			{
				UID:     "0",
				Creator: simapp.TestParamUsers["user1"].Address.String(),
			},
			{
				UID:     "1",
				Creator: simapp.TestParamUsers["user2"].Address.String(),
			},
		},
		Uid2IdList: []types.UID2ID{
			{
				UID: "0",
				ID:  1,
			},
			{
				UID: "1",
				ID:  2,
			},
		},
		PendingBetList: []types.PendingBet{
			{
				UID:     "1",
				Creator: simapp.TestParamUsers["user1"].Address.String(),
			},
		},
		SettledBetList: []types.SettledBet{
			{
				UID:           "1",
				BettorAddress: simapp.TestParamUsers["user1"].Address.String(),
			},
		},
		Stats: types.BetStats{
			Count: 2,
		},
	}

	bet.InitGenesis(ctx, *tApp.BetKeeper, genesisState)
	got := bet.ExportGenesis(ctx, *tApp.BetKeeper)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.BetList, got.BetList)
}
