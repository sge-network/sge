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
	promoterAddr := "promoter"
	promoterUID := uuid.NewString()

	campaignUID1 := uuid.NewString()
	campaignUID2 := uuid.NewString()

	rewardID1 := uuid.NewString()
	rewardID2 := uuid.NewString()

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		PromoterList: []types.Promoter{
			{
				Creator:   promoterAddr,
				UID:       promoterUID,
				Addresses: []string{promoterAddr},
				Conf: types.PromoterConf{
					CategoryCap: []types.CategoryCap{
						{
							Category:  types.RewardCategory_REWARD_CATEGORY_SIGNUP,
							CapPerAcc: 1,
						},
					},
				},
			},
		},
		PromoterByAddressList: []types.PromoterByAddress{
			{
				PromoterUID: promoterUID,
				Address:     promoterAddr,
			},
		},
		CampaignList: []types.Campaign{
			{
				UID:      campaignUID1,
				Promoter: promoterAddr,
			},
			{
				UID:      campaignUID2,
				Promoter: promoterAddr,
			},
		},
		RewardList: []types.Reward{
			{
				UID:         rewardID1,
				CampaignUID: campaignUID1,
			},
			{
				UID:         rewardID2,
				CampaignUID: campaignUID2,
			},
		},
		RewardByCategoryList: []types.RewardByCategory{
			{
				UID: rewardID1,
			},
			{
				UID: rewardID2,
			},
		},
		RewardByCampaignList: []types.RewardByCampaign{
			{
				UID:         rewardID1,
				CampaignUID: campaignUID1,
			},
			{
				UID:         rewardID1,
				CampaignUID: campaignUID2,
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
