package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/ovm/types"
)

const (
	testAddress = "cosmos1s4ycalgh3gjemd4hmqcvcgmnf647rnd0tpg2w9"
)

func TestGenesisState_Validate(t *testing.T) {
	pubkeys := simapp.GenerateOvmPublicKeys(types.MaxPubKeysCount)

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

	corruptPublicKeys := make([]string, len(pubkeys))
	copy(corruptPublicKeys, pubkeys)
	corruptPublicKeys[0] = "-----BEGIN PUBLIC KEY-----\n \n-----END PUBLIC KEY-----"
	invalidPublicKeysState := types.GenesisState{
		KeyVault: types.KeyVault{
			PublicKeys: corruptPublicKeys,
		},
		PubkeysChangeProposals: validState.PubkeysChangeProposals,
		ProposalStats:          validState.ProposalStats,
		Params:                 types.DefaultParams(),
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
		{
			desc:     "invalid genesis state",
			genState: &invalidPublicKeysState,
			valid:    false,
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
