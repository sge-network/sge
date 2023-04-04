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
		items[i].Modifications = types.PubkeysChangeProposalPayload{PublicKeys: pubKeys, LeaderIndex: 0}
		items[i].StartTS = ctx.BlockTime().Unix()
		items[i].Status = types.ProposalStatus_PROPOSAL_STATUS_ACTIVE

		keeper.SetPubkeysChangeProposal(ctx, items[i])
	}
	return items
}

func TestGetActivePubkeysChangeProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNActiveProposal(k, ctx, 10)
	_, found := k.GetPubkeysChangeProposal(ctx, types.ProposalStatus_PROPOSAL_STATUS_ACTIVE, 5000000)
	require.False(t, found)

	for _, item := range items {
		rst, found := k.GetPubkeysChangeProposal(ctx, types.ProposalStatus_PROPOSAL_STATUS_ACTIVE, item.Id)
		require.True(t, found)
		require.EqualValues(t, item, rst)
	}
}

func TestRemoveActivePubkeysChangeProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNActiveProposal(k, ctx, 10)
	for _, item := range items {
		k.RemoveProposal(ctx,
			types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
			item.Id,
		)
		_, found := k.GetPubkeysChangeProposal(ctx,
			types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
			item.Id,
		)
		require.False(t, found)
	}
}

func TestGetAllActivePubkeysChangeProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNActiveProposal(k, ctx, 10)

	Markets, err := k.GetAllPubkeysChangeProposalsByStatus(ctx, types.ProposalStatus_PROPOSAL_STATUS_ACTIVE)
	require.NoError(t, err)

	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(Markets),
	)
}

func TestFinishProposals(t *testing.T) {
	k, ctx := setupKeeper(t)
	ctx = ctx.WithBlockTime(time.Now())
	items := createNActiveProposal(k, ctx, 10)

	proposals, err := k.GetAllPubkeysChangeProposalsByStatus(ctx, types.ProposalStatus_PROPOSAL_STATUS_ACTIVE)
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

		k.SetPubkeysChangeProposal(ctx, proposal)
	}

	proposals, err = k.GetAllPubkeysChangeProposalsByStatus(ctx, types.ProposalStatus_PROPOSAL_STATUS_ACTIVE)
	require.NoError(t, err)
	require.Equal(t, len(items), len(proposals))

	err = k.FinishProposals(ctx)
	require.NoError(t, err)

	proposals, err = k.GetAllPubkeysChangeProposalsByStatus(ctx, types.ProposalStatus_PROPOSAL_STATUS_ACTIVE)
	require.NoError(t, err)
	require.Equal(t, 0, len(proposals))

	finishedProposals, err := k.GetAllPubkeysChangeProposalsByStatus(ctx, types.ProposalStatus_PROPOSAL_STATUS_FINISHED)
	require.NoError(t, err)
	require.Equal(t, len(items), len(finishedProposals))

	finishedProposals, err = k.GetAllPubkeysChangeProposals(ctx)
	require.NoError(t, err)

	require.Equal(t, len(items), len(finishedProposals))

	keyVault, found := k.GetKeyVault(ctx)
	require.True(t, found)
	require.Equal(t, 5, len(keyVault.PublicKeys))
}

func TestFinishProposal(t *testing.T) {
	k, ctx := setupKeeper(t)
	ctx = ctx.WithBlockTime(time.Now())
	proposals := createNActiveProposal(k, ctx, 3)

	voteAll := func(proposal *types.PublicKeysChangeProposal, vote types.ProposalVote) {
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
		proposal *types.PublicKeysChangeProposal

		err error
	}{
		{
			timeDiff: types.MaxValidProposalMinutes / 2 * time.Minute,
			proposal: &proposals[0],
			vote:     types.ProposalVote_PROPOSAL_VOTE_YES,
			result:   types.ProposalResult_PROPOSAL_RESULT_APPROVED,
		},
		{
			timeDiff: types.MaxValidProposalMinutes / 2 * time.Minute,
			proposal: &proposals[1],
			vote:     types.ProposalVote_PROPOSAL_VOTE_NO,
			result:   types.ProposalResult_PROPOSAL_RESULT_REJECTED,
		},
		{
			timeDiff: types.MaxValidProposalMinutes * 2 * time.Minute,
			proposal: &proposals[2],
			vote:     types.ProposalVote_PROPOSAL_VOTE_YES,
			result:   types.ProposalResult_PROPOSAL_RESULT_EXPIRED,
		},
	} {
		ctx = ctx.WithBlockTime(ctx.BlockTime().Add(tc.timeDiff))

		voteAll(tc.proposal, tc.vote)

		k.SetPubkeysChangeProposal(ctx, *tc.proposal)

		err := k.FinishProposals(ctx)
		require.NoError(t, err)

		finishedProposal, found := k.GetPubkeysChangeProposal(ctx, types.ProposalStatus_PROPOSAL_STATUS_FINISHED, tc.proposal.Id)
		require.True(t, found)
		require.Equal(t, tc.result, finishedProposal.Result)
	}
}
