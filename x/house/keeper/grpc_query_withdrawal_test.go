package keeper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/house/types"
)

func TestWithdrawQuerySingle(t *testing.T) {
	k, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNWithdrawals(k, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryWithdrawalRequest
		response *types.QueryWithdrawalResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryWithdrawalRequest{
				DepositorAddress:   msgs[0].Address,
				MarketUid:          msgs[0].MarketUID,
				ParticipationIndex: msgs[0].ParticipationIndex,
				Id:                 msgs[0].ID,
			},
			response: &types.QueryWithdrawalResponse{Withdrawal: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryWithdrawalRequest{
				DepositorAddress:   msgs[1].Address,
				MarketUid:          msgs[1].MarketUID,
				ParticipationIndex: msgs[1].ParticipationIndex,
				Id:                 msgs[1].ID,
			},
			response: &types.QueryWithdrawalResponse{Withdrawal: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryWithdrawalRequest{
				DepositorAddress:   msgs[0].Address,
				MarketUid:          uuid.NewString(),
				ParticipationIndex: msgs[0].ParticipationIndex,
				Id:                 msgs[0].ID,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.Withdrawal(wctx, tc.request)
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
