package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/orderbook/client/cli"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func networkWithOddsExposureObjects(t *testing.T, n int) (*network.Network, []types.OrderBookOddsExposure) {
	t.Helper()
	cfg := network.DefaultConfig()

	// orderbook module state
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		exposure := types.OrderBookOddsExposure{
			OrderBookUID:     testMarketUID,
			OddsUID:          uuid.NewString(),
			FulfillmentQueue: []uint64{},
		}
		nullify.Fill(&exposure)

		state.OrderBookExposureList = append(state.OrderBookExposureList, exposure)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), state.OrderBookExposureList
}

func TestQueryOddsExposure(t *testing.T) {
	net, objs := networkWithOddsExposureObjects(t, 5)

	t.Run("ShowParticipation", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		common := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		for _, tc := range []struct {
			desc         string
			orderBookUID string
			oddsUID      string

			args []string
			err  error
			obj  types.OrderBookOddsExposure
		}{
			{
				desc:         "found",
				orderBookUID: testMarketUID,
				oddsUID:      objs[0].OddsUID,

				args: common,
				obj:  objs[0],
			},
			{
				desc:         "not found",
				oddsUID:      "dummy",
				orderBookUID: cast.ToString(100000),

				args: common,
				err:  status.Error(codes.NotFound, "not found"),
			},
		} {
			tc := tc
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					tc.orderBookUID,
					cast.ToString(tc.oddsUID),
				}
				args = append(args, tc.args...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookExposure(), args)
				if tc.err != nil {
					stat, ok := status.FromError(tc.err)
					require.True(t, ok)
					require.ErrorIs(t, stat.Err(), tc.err)
				} else {
					require.NoError(t, err)
					var resp types.OrderBookOddsExposure
					require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
					require.NotNil(t, resp)
					require.Equal(t,
						nullify.Fill(&tc.obj),
						nullify.Fill(&resp),
					)
				}
			})
		}
	})

	t.Run("ListOrderBookExposures", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		request := func(next []byte, offset, limit uint64, total bool) []string {
			args := []string{
				testMarketUID,
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookExposures(), args)
				require.NoError(t, err)
				var resp types.QueryOrderBookExposuresResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.OrderBookExposures), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.OrderBookExposures),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookExposures(), args)
				require.NoError(t, err)
				var resp types.QueryOrderBookExposuresResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.OrderBookExposures), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.OrderBookExposures),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookExposures(), args)
			require.NoError(t, err)
			var resp types.QueryOrderBookExposuresResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.OrderBookExposures),
			)
		})
	})
}
