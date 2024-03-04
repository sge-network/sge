package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/orderbook/types"
)

func TestOrderBookParticipationExposuresQueryPaginated(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNParticipationExposure(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryOrderBookParticipationExposuresRequest {
		return &types.QueryOrderBookParticipationExposuresRequest{
			OrderBookUid: testOrderBookUID,
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
			resp, err := k.OrderBookParticipationExposures(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ParticipationExposures), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ParticipationExposures),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.OrderBookParticipationExposures(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ParticipationExposures), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ParticipationExposures),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.OrderBookParticipationExposures(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ParticipationExposures),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.OrderBookParticipationExposures(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}

func TestParticipationExposuresQueryPaginated(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNParticipationExposure(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryParticipationExposuresRequest {
		return &types.QueryParticipationExposuresRequest{
			OrderBookUid:       testOrderBookUID,
			ParticipationIndex: testParticipationIndex,
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
			resp, err := k.ParticipationExposures(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ParticipationExposures), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ParticipationExposures),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.ParticipationExposures(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ParticipationExposures), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ParticipationExposures),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.ParticipationExposures(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ParticipationExposures),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.ParticipationExposures(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}

func TestHistoricalParticipationExposuresQueryPaginated(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNParticipationExposure(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryHistoricalParticipationExposuresRequest {
		return &types.QueryHistoricalParticipationExposuresRequest{
			OrderBookUid: testOrderBookUID,
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
			resp, err := k.HistoricalParticipationExposures(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ParticipationExposures), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ParticipationExposures),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.HistoricalParticipationExposures(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ParticipationExposures), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ParticipationExposures),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.HistoricalParticipationExposures(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ParticipationExposures),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.HistoricalParticipationExposures(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}
