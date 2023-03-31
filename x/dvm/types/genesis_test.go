package types_test

import (
	"testing"
	"time"

	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

const (
	testAddress = "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"
)

func TestGenesisState_Validate(t *testing.T) {
	pubkeys := simappUtil.GenerateDvmPublicKeys(5)

	var votes []*types.Vote
	for _, v := range pubkeys {
		votes = append(votes, &types.Vote{PublicKey: v, Vote: types.ProposalVote_PROPOSAL_VOTE_YES})
	}

	validState := types.GenesisState{
		KeyVault: types.KeyVault{
			PublicKeys: pubkeys,
		},
		PubkeysChangeProposals: []types.PublicKeysChangeProposal{
			{
				Id:      1,
				Creator: testAddress,
				StartTS: time.Now().Unix(),
				Modifications: types.PubkeysChangeProposalPayload{
					PublicKeys:  pubkeys,
					LeaderIndex: 0,
				},
				Votes:  votes,
				Status: types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
			},
		},
		ProposalStats: types.ProposalStats{
			PubkeysChangeCount: 1,
		},
		Params: types.DefaultParams(),
	}

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is not valid",
			genState: types.DefaultGenesis(),
			valid:    false,
		},
		{
			desc:     "valid genesis state",
			genState: &validState,
			valid:    true,
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
