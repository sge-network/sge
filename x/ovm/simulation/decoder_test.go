package simulation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/sge-network/sge/app"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/ovm/simulation"
	"github.com/sge-network/sge/x/ovm/types"
)

func TestDecodeStore(t *testing.T) {
	cdc := app.MakeEncodingConfig().Marshaler
	dec := simulation.NewDecodeStore(cdc)

	keyVault := types.KeyVault{
		PublicKeys: []string{"sample key"},
	}

	proposalStats := types.ProposalStats{
		PubkeysChangeCount: 1,
	}

	proposal := types.PublicKeysChangeProposal{
		Id:      1,
		Creator: sample.AccAddress(),
		StartTS: cast.ToInt64(time.Now().UTC()),
		Modifications: types.PubkeysChangeProposalPayload{
			PublicKeys:  []string{"new key"},
			LeaderIndex: 0,
		},
		Votes: []*types.Vote{
			{
				Vote:      types.ProposalVote_PROPOSAL_VOTE_YES,
				PublicKey: "sample key",
			},
		},
		Result:     types.ProposalResult_PROPOSAL_RESULT_UNSPECIFIED,
		ResultMeta: "sample metadata",
		Status:     types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
		FinishTS:   cast.ToInt64(time.Now().Add(1 * time.Hour).UTC().String()),
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.KeyVaultKey, Value: cdc.MustMarshal(&keyVault)},
			{Key: types.ProposalStatsKey, Value: cdc.MustMarshal(&proposalStats)},
			{Key: types.PubKeysChangeProposalListPrefix, Value: cdc.MustMarshal(&proposal)},
			{Key: []byte{0x99}, Value: []byte{0x99}},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"key_vault", fmt.Sprintf("%v\n%v", keyVault, keyVault)},
		{"proposals_stats", fmt.Sprintf("%v\n%v", proposalStats, proposalStats)},
		{"pubkeys_change_proposal", fmt.Sprintf("%v\n%v", proposal, proposal)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
