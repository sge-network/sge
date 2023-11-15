package reward_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/testutil/simapp"
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
		RewardList: []types.Reward{
			{
				UID: uuid.NewString(),
			},
			{
				UID: uuid.NewString(),
			},
		},
		RewardByRecCatList: []types.RewardByCategory{
			{
				UID: uuid.NewString(),
			},
			{
				UID: uuid.NewString(),
			},
		},
		RewardByCampaignList: []types.RewardByCampaign{
			{
				UID: uuid.NewString(),
			},
			{
				UID: uuid.NewString(),
			},
		},
	}

	tApp, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)
	reward.InitGenesis(ctx, *tApp.RewardKeeper, genesisState)
	got := reward.ExportGenesis(ctx, *tApp.RewardKeeper)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.CampaignList, got.CampaignList)
}
