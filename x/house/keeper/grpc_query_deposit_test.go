package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/house/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDepositsQueryPaginated(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDeposits(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryDepositsRequest {
		return &types.QueryDepositsRequest{
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
			resp, err := k.Deposits(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deposits), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Deposits),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.Deposits(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deposits), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Deposits),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.Deposits(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Deposits),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.Deposits(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}

func TestDepositsByAccountQueryPaginated(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNDeposits(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryDepositsByAccountRequest {
		return &types.QueryDepositsByAccountRequest{
			Address: testDepositorAddress,
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
			resp, err := k.DepositsByAccount(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deposits), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Deposits),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.DepositsByAccount(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deposits), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Deposits),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.DepositsByAccount(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Deposits),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.DepositsByAccount(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}
