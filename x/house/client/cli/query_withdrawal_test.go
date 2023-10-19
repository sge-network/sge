package cli_test

import (
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/house/client/cli"
	"github.com/sge-network/sge/x/house/types"
	"github.com/stretchr/testify/require"
)

func TestQueryWithdraw(t *testing.T) {
	net, _, objs := networkWithDepositObjects(t, 5)

	t.Run("ListWithdrawalsByAccount", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		request := func(next []byte, offset, limit uint64, total bool) []string {
			args := []string{
				testAddress,
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			}
			if next == nil {
				args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
			} else {
				args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
			}
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
			if total {
				args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
			}
			return args
		}
		t.Run("ByOffset", func(t *testing.T) {
			step := 2
			for i := 0; i < len(objs); i += step {
				args := request(nil, uint64(i), uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryWithdrawalsByAccount(), args)
				require.NoError(t, err)
				var resp types.QueryWithdrawalsByAccountResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Withdrawals), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Withdrawals),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryWithdrawalsByAccount(), args)
				require.NoError(t, err)
				var resp types.QueryWithdrawalsByAccountResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Withdrawals), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Withdrawals),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryWithdrawalsByAccount(), args)
			require.NoError(t, err)
			var resp types.QueryWithdrawalsByAccountResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Withdrawals),
			)
		})
	})
}
