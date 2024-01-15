package keeper_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/reward/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestCampaignQuerySingle(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCampaign(k, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryCampaignRequest
		response *types.QueryCampaignResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryCampaignRequest{
				Uid: msgs[0].UID,
			},
			response: &types.QueryCampaignResponse{Campaign: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryCampaignRequest{
				Uid: msgs[1].UID,
			},
			response: &types.QueryCampaignResponse{Campaign: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryCampaignRequest{
				Uid: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.Campaign(wctx, tc.request)
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

func TestCampaignQueryPaginated(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNCampaign(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryCampaignsRequest {
		return &types.QueryCampaignsRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := k.Campaigns(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Campaign), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Campaign),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.Campaigns(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Campaign), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Campaign),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.Campaigns(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Campaign),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.Campaigns(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
