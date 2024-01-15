package keeper_test

import (
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/orderbook/types"
)

func TestOrderBookQuerySingle(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNOrderBook(tApp, k, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryOrderBookRequest
		response *types.QueryOrderBookResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryOrderBookRequest{
				OrderBookUid: msgs[0].UID,
			},
			response: &types.QueryOrderBookResponse{OrderBook: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryOrderBookRequest{
				OrderBookUid: msgs[1].UID,
			},
			response: &types.QueryOrderBookResponse{OrderBook: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryOrderBookRequest{
				OrderBookUid: cast.ToString(100000),
			},
			err: status.Error(codes.NotFound, "order book 100000 not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.OrderBook(wctx, tc.request)
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

func TestOrderBooksQueryPaginated(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNOrderBook(tApp, k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryOrderBooksRequest {
		return &types.QueryOrderBooksRequest{
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
			resp, err := k.OrderBooks(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Orderbooks), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Orderbooks),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.OrderBooks(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Orderbooks), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Orderbooks),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.OrderBooks(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Orderbooks),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.OrderBooks(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}
