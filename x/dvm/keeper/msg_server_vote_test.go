package keeper_test

import (
	"testing"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang-jwt/jwt/v4"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
)

func TestChangePubkeysVote(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

		createNActiveProposal(k, ctx, 1)
		creator := simappUtil.TestParamUsers["user1"]
		pubs, err := createNTestPubKeys(5)
		require.NoError(t, err)

		proposalTicket := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
			"public_keys":  pubs,
			"leader_index": 0,
			"exp":          jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})
		singedProposalTicket, err := proposalTicket.SignedString(simappUtil.TestDVMPrivateKeys[0])
		require.NoError(t, err)

		resp, err := msgk.SubmitPubkeysChangeProposal(wctx, &types.MsgSubmitPubkeysChangeProposalRequest{
			Creator: creator.Address.String(),
			Ticket:  singedProposalTicket,
		})
		require.NoError(t, err)
		require.Equal(t, true, resp.Success)

		voteTicket := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
			"proposal_id": 1,
			"vote":        types.ProposalVote_PROPOSAL_VOTE_YES,
			"exp":         jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})
		singedVoteTicket, err := voteTicket.SignedString(simappUtil.TestDVMPrivateKeys[0])
		require.NoError(t, err)

		respVote, err := msgk.VotePubkeysChange(wctx, &types.MsgVotePubkeysChangeRequest{
			Creator:       creator.Address.String(),
			Ticket:        singedVoteTicket,
			VoterKeyIndex: 0,
		})
		require.NoError(t, err)
		require.Equal(t, true, respVote.Success)
	})

	t.Run("duplicate vote", func(t *testing.T) {
		k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

		createNActiveProposal(k, ctx, 1)
		creator := simappUtil.TestParamUsers["user1"]
		pubs, err := createNTestPubKeys(5)
		require.NoError(t, err)

		proposalTicket := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
			"public_keys":  pubs,
			"leader_index": 0,
			"exp":          jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})
		singedProposalTicket, err := proposalTicket.SignedString(simappUtil.TestDVMPrivateKeys[0])
		require.NoError(t, err)

		resp, err := msgk.SubmitPubkeysChangeProposal(wctx, &types.MsgSubmitPubkeysChangeProposalRequest{
			Creator: creator.Address.String(),
			Ticket:  singedProposalTicket,
		})
		require.NoError(t, err)
		require.Equal(t, true, resp.Success)

		vote1Ticket := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
			"proposal_id": 1,
			"vote":        types.ProposalVote_PROPOSAL_VOTE_YES,
			"exp":         jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})
		singedVote1Ticket, err := vote1Ticket.SignedString(simappUtil.TestDVMPrivateKeys[0])
		require.NoError(t, err)

		respVote, err := msgk.VotePubkeysChange(wctx, &types.MsgVotePubkeysChangeRequest{
			Creator:       creator.Address.String(),
			Ticket:        singedVote1Ticket,
			VoterKeyIndex: 0,
		})
		require.NoError(t, err)
		require.Equal(t, true, respVote.Success)

		vote2Ticket := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
			"proposal_id": 1,
			"vote":        types.ProposalVote_PROPOSAL_VOTE_YES,
			"exp":         jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		})
		singedVote2Ticket, err := vote2Ticket.SignedString(simappUtil.TestDVMPrivateKeys[0])
		require.NoError(t, err)

		respVote, err = msgk.VotePubkeysChange(wctx, &types.MsgVotePubkeysChangeRequest{
			Creator:       creator.Address.String(),
			Ticket:        singedVote2Ticket,
			VoterKeyIndex: 0,
		})
		require.ErrorIs(t, sdkerrors.ErrInvalidRequest, err)
		require.Nil(t, respVote)
	})
}
