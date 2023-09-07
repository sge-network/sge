package reward_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/stretchr/testify/require"

	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/reward"
	"github.com/sge-network/sge/x/reward/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		CampaignList: []types.Campaign{
			{
				UID: uuid.NewString(),
			},
			{
				UID: uuid.NewString(),
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	tApp, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)
	reward.InitGenesis(ctx, *tApp.RewardKeeper, genesisState)
	got := reward.ExportGenesis(ctx, *tApp.RewardKeeper)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.CampaignList, got.CampaignList)
	// this line is used by starport scaffolding # genesis/test/assert
}
