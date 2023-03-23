package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/testutil/nullify"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/types"
)

func TestBetQuerySingle(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNBet(tApp, k, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryBetRequest
		response *types.QueryBetResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryBetRequest{
				Creator: testCreator,
				Uid:     msgs[0].UID,
			},
			response: &types.QueryBetResponse{Bet: msgs[0], Market: testMarket},
		},
		{
			desc: "Second",
			request: &types.QueryBetRequest{
				Creator: testCreator,
				Uid:     msgs[1].UID,
			},
			response: &types.QueryBetResponse{Bet: msgs[1], Market: testMarket},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryBetRequest{
				Uid: cast.ToString(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.Bet(wctx, tc.request)
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

func TestBetQueryPaginated(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNBet(tApp, k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryBetsRequest {
		return &types.QueryBetsRequest{
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
			resp, err := k.Bets(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bet), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bet),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.Bets(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bet), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bet),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.Bets(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Bet),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.Bets(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}

// TestBetQueryPaginatedReverse test if IDs are sorted reveresely
func TestBetQueryPaginatedReverse(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNBet(tApp, k, ctx, 100)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryBetsRequest {
		return &types.QueryBetsRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
				Reverse:    true,
			},
		}
	}
	t.Run("Sorted", func(t *testing.T) {
		resp, err := k.Bets(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))

		lastBetID := resp.Pagination.Total + 1
		for _, b := range resp.Bet {
			uuid2ID, found := k.GetBetID(ctx, b.UID)
			require.True(t, found)
			require.Less(t, uuid2ID.ID, lastBetID)
			lastBetID = uuid2ID.ID
		}
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.Bets(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}

func TestBetsByCreatorQueryPaginated(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNBet(tApp, k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryBetsByCreatorRequest {
		return &types.QueryBetsByCreatorRequest{
			Creator: simappUtil.TestParamUsers["user1"].Address.String(),
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
			resp, err := k.BetsByCreator(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bet), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bet),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.BetsByCreator(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bet), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bet),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.BetsByCreator(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Bet),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.Bets(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}

func TestBetByCreatorQueryPaginatedReverse(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNBet(tApp, k, ctx, 100)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryBetsByCreatorRequest {
		return &types.QueryBetsByCreatorRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
				Reverse:    true,
			},
		}
	}
	t.Run("Sorted", func(t *testing.T) {
		resp, err := k.BetsByCreator(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))

		lastBetID := resp.Pagination.Total + 1
		for _, b := range resp.Bet {
			uuid2ID, found := k.GetBetID(ctx, b.UID)
			require.True(t, found)
			require.Less(t, uuid2ID.ID, lastBetID)
			lastBetID = uuid2ID.ID
		}
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.Bets(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}

func TestBetByUIDsQuery(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNBet(tApp, k, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryBetsByUIDsRequest
		response *types.QueryBetsByUIDsResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryBetsByUIDsRequest{
				Items: []string{
					simappUtil.TestParamUsers["user1"].Address.String() + ":" + msgs[0].UID,
				},
			},
			response: &types.QueryBetsByUIDsResponse{Bets: []types.Bet{msgs[0]}},
		},
		{
			desc: "Second",
			request: &types.QueryBetsByUIDsRequest{
				Items: []string{
					simappUtil.TestParamUsers["user1"].Address.String() + ":" + msgs[1].UID,
				},
			},
			response: &types.QueryBetsByUIDsResponse{Bets: []types.Bet{msgs[1]}},
		},
		{
			desc: "NotFound",
			request: &types.QueryBetsByUIDsRequest{
				Items: []string{
					simappUtil.TestParamUsers["user1"].Address.String() + ":" + "100000",
				},
			},
			response: &types.QueryBetsByUIDsResponse{Bets: []types.Bet{}, NotFoundBetUids: []string{"100000"}},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.BetsByUIDs(wctx, tc.request)
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

func TestPendingBetsOfMarketQueryPaginated(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNBet(tApp, k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryPendingBetsRequest {
		return &types.QueryPendingBetsRequest{
			MarketUid: testMarketUID,
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
			resp, err := k.PendingBets(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bet), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bet),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.PendingBets(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Bet), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Bet),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.PendingBets(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Bet),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.Bets(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}
