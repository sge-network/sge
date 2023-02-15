package keeper_test

import (
	"crypto/x509"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/nullify"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/dvm/keeper"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func createNActiveProposal(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PublicKeysChangeProposal {
	items := make([]types.PublicKeysChangeProposal, n)

	pubKeys, err := createNTestPubKeys(5)
	if err != nil {
		panic(err)
	}

	for i := range items {
		items[i].Id = uint64(i)
		items[i].Creator = simappUtil.TestParamUsers["user"+cast.ToString(i)].Address.String()
		items[i].Modifications = types.PubkeysChangeProposalPayload{Additions: pubKeys}
		items[i].StartTS = ctx.BlockTime().Unix()

		keeper.SetActivePubkeysChangeProposal(ctx, items[i])
	}
	return items
}

func TestGetActivePubkeysChangeProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNActiveProposal(k, ctx, 10)
	_, found := k.GetActivePubkeysChangeProposal(ctx, 5000000)
	require.False(t, found)

	for _, item := range items {
		rst, found := k.GetActivePubkeysChangeProposal(ctx, item.Id)
		require.True(t, found)
		require.EqualValues(t, item, rst)
	}
}

func TestRemoveActivePubkeysChangeProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNActiveProposal(k, ctx, 10)
	for _, item := range items {
		k.RemoveActiveProposal(ctx,
			item.Id,
		)
		_, found := k.GetActivePubkeysChangeProposal(ctx,
			item.Id,
		)
		require.False(t, found)
	}
}

func TestGetAllActivePubkeysChangeProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNActiveProposal(k, ctx, 10)

	sportEvents, err := k.GetAllActivePubkeysChangeProposals(ctx)
	require.NoError(t, err)

	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(sportEvents),
	)
}

func TestFinishProposals(t *testing.T) {
	k, ctx := setupKeeper(t)
	ctx = ctx.WithBlockTime(time.Now())
	items := createNActiveProposal(k, ctx, 10)

	proposals, err := k.GetAllActivePubkeysChangeProposals(ctx)
	require.NoError(t, err)
	require.Equal(t, len(items), len(proposals))

	// set the block time equal to the half of the valid time range.
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(types.MaxValidProposalMinutes / 2 * time.Minute))

	for _, proposal := range items {
		for _, p := range simappUtil.TestDVMPublicKeys {
			bs, err := x509.MarshalPKIXPublicKey(p)
			if err != nil {
				panic(err)
			}

			proposal.Votes = append(proposal.Votes,
				types.NewVote(
					string(utils.NewPubKeyMemory(bs)),
					types.ProposalVote_PROPOSAL_VOTE_YES,
				),
			)
		}

		k.SetActivePubkeysChangeProposal(ctx, proposal)
	}

	proposals, err = k.GetAllActivePubkeysChangeProposals(ctx)
	require.NoError(t, err)
	require.Equal(t, len(items), len(proposals))

	err = k.FinishProposals(ctx)
	require.NoError(t, err)

	proposals, err = k.GetAllActivePubkeysChangeProposals(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(proposals))

	finishedProposals, err := k.GetAllFinishedPubkeysChangeProposals(ctx)
	require.NoError(t, err)
	require.Equal(t, len(items), len(finishedProposals))

	keyVault, found := k.GetKeyVault(ctx)
	require.True(t, found)
	// 55 = 5*10 +  5(existing public keys in the genesis)
	require.Equal(t, 55, len(keyVault.PublicKeys))
}

func TestFinishProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	ctx = ctx.WithBlockTime(time.Now())
	proposal := createNActiveProposal(k, ctx, 1)[0]

	voteAll := func(vote types.ProposalVote) {
		for _, p := range simappUtil.TestDVMPublicKeys {
			bs, err := x509.MarshalPKIXPublicKey(p)
			if err != nil {
				panic(err)
			}

			proposal.Votes = append(proposal.Votes,
				types.NewVote(
					string(utils.NewPubKeyMemory(bs)),
					vote,
				),
			)
		}
	}

	for _, tc := range []struct {
		timeDiff time.Duration
		vote     types.ProposalVote
		result   types.ProposalResult

		err error
	}{
		{
			timeDiff: types.MaxValidProposalMinutes / 2 * time.Minute,
			vote:     types.ProposalVote_PROPOSAL_VOTE_YES,
			result:   types.ProposalResult_PROPOSAL_RESULT_APPROVED,
		},
		{
			timeDiff: types.MaxValidProposalMinutes / 2 * time.Minute,
			vote:     types.ProposalVote_PROPOSAL_VOTE_NO,
			result:   types.ProposalResult_PROPOSAL_RESULT_REJECTED,
		},
		{
			timeDiff: types.MaxValidProposalMinutes * 2 * time.Minute,
			vote:     types.ProposalVote_PROPOSAL_VOTE_YES,
			result:   types.ProposalResult_PROPOSAL_RESULT_EXPIRED,
		},
	} {
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(tc.timeDiff))

		voteAll(tc.vote)

		k.SetActivePubkeysChangeProposal(ctx, proposal)

		err := k.FinishProposals(ctx)
		require.NoError(t, err)

		finishedProposal, found := k.GetFinishedPubkeysChangeProposal(ctx, proposal.Id)
		require.True(t, found)
		require.Equal(t, tc.result, finishedProposal.Result)
	}
}
