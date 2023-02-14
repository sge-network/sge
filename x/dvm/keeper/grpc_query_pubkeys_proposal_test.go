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
		request  *types.QueryActivePublicKeysChangeProposalRequest
		response *types.QueryActivePublicKeysChangeProposalResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryActivePublicKeysChangeProposalRequest{
				Id: msgs[0].Id,
			},
			response: &types.QueryActivePublicKeysChangeProposalResponse{Proposal: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryActivePublicKeysChangeProposalRequest{
				Id: msgs[1].Id,
			},
			response: &types.QueryActivePublicKeysChangeProposalResponse{Proposal: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryActivePublicKeysChangeProposalRequest{
				Id: 100000,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.ActivePublicKeysChangeProposal(wctx, tc.request)
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
