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

func TestWithdrawByAccountQueryPaginated(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNWithdrawals(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryWithdrawalsByAccountRequest {
		return &types.QueryWithdrawalsByAccountRequest{
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
			resp, err := k.WithdrawalsByAccount(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Withdrawals), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Withdrawals),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.WithdrawalsByAccount(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Withdrawals), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Withdrawals),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.WithdrawalsByAccount(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Withdrawals),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.WithdrawalsByAccount(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest))
	})
}
