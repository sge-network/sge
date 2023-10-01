package types_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/reward/types"
)

func TestGenesisState_Validate(t *testing.T) {
	dupplicateCampaign := uuid.NewString()
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
			},
			valid: true,
		},
		{
			desc: "duplicated campaign",
			genState: &types.GenesisState{
				CampaignList: []types.Campaign{
					{
						UID: dupplicateCampaign,
					},
					{
						UID: dupplicateCampaign,
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
