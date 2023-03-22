package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestActivePubkeysChangeProposalQuerySingle(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNActiveProposal(k, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryPublicKeysChangeProposalRequest
		response *types.QueryPublicKeysChangeProposalResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryPublicKeysChangeProposalRequest{
				Id:     msgs[0].Id,
				Status: types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
			},
			response: &types.QueryPublicKeysChangeProposalResponse{Proposal: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryPublicKeysChangeProposalRequest{
				Id:     msgs[1].Id,
				Status: types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
			},
			response: &types.QueryPublicKeysChangeProposalResponse{Proposal: msgs[1]},
		},
		{
			desc: "ActiveNotFound",
			request: &types.QueryPublicKeysChangeProposalRequest{
				Id:     msgs[1].Id,
				Status: types.ProposalStatus_PROPOSAL_STATUS_FINISHED,
			},
			err: status.Error(codes.NotFound, "not found")},
		{
			desc: "KeyNotFound",
			request: &types.QueryPublicKeysChangeProposalRequest{
				Id:     100000,
				Status: types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.PublicKeysChangeProposal(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
