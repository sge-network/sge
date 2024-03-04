package types_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/reward/types"
)

func TestGenesisState_Validate(t *testing.T) {
	dupplicateUID := uuid.NewString()
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
			desc: "valid genesis state",
			genState: &types.GenesisState{
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
				RewardByCategoryList: []types.RewardByCategory{
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
			},
			valid: true,
		},
		{
			desc: "duplicated campaign",
			genState: &types.GenesisState{
				CampaignList: []types.Campaign{
					{
						UID: dupplicateUID,
					},
					{
						UID: dupplicateUID,
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated reward",
			genState: &types.GenesisState{
				RewardList: []types.Reward{
					{
						UID: dupplicateUID,
					},
					{
						UID: dupplicateUID,
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated reward by category",
			genState: &types.GenesisState{
				RewardByCategoryList: []types.RewardByCategory{
					{
						UID: dupplicateUID,
					},
					{
						UID: dupplicateUID,
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated reward by campaign",
			genState: &types.GenesisState{
				RewardByCampaignList: []types.RewardByCampaign{
					{
						UID: dupplicateUID,
					},
					{
						UID: dupplicateUID,
					},
				},
			},
			valid: false,
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
