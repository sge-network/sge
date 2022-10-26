package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/sportevent/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestSportEventQuerySingle(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNSportEvent(k, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QuerySportEventRequest
		response *types.QuerySportEventResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QuerySportEventRequest{
				Uid: msgs[0].UID,
			},
			response: &types.QuerySportEventResponse{SportEvent: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QuerySportEventRequest{
				Uid: msgs[1].UID,
			},
			response: &types.QuerySportEventResponse{SportEvent: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QuerySportEventRequest{
				Uid: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.SportEvent(wctx, tc.request)
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

func TestSportEventQueryPaginated(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNSportEvent(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QuerySportEventListAllRequest {
		return &types.QuerySportEventListAllRequest{
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
			resp, err := k.SportEventAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.SportEvent), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.SportEvent),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.SportEventAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.SportEvent), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.SportEvent),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.SportEventAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.SportEvent),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.SportEventAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}
