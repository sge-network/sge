package cli_test

import (
	"fmt"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tmcli "github.com/tendermint/tendermint/libs/cli"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/orderbook/client/cli"
	"github.com/sge-network/sge/x/orderbook/types"
)

const (
	testMarketUID          = "5db09053-2901-4110-8fb5-c14e21f8d555"
	testParticipationIndex = 1
)

func networkWithParticipationObjects(t *testing.T, n int) (*network.Network, []types.OrderBookParticipation) {
	t.Helper()
	cfg := network.DefaultConfig()

	// orderbook module state
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		participation := types.OrderBookParticipation{
			Index:                      cast.ToUint64(i + 1),
			OrderBookUID:               testMarketUID,
			ParticipantAddress:         testAddress,
			Liquidity:                  sdkmath.NewInt(1000),
			Fee:                        sdkmath.NewInt(10),
			CurrentRoundLiquidity:      sdkmath.NewInt(1000),
			ExposuresNotFilled:         1,
			TotalBetAmount:             sdkmath.NewInt(0),
			CurrentRoundTotalBetAmount: sdkmath.NewInt(0),
			MaxLoss:                    sdkmath.NewInt(0),
			CurrentRoundMaxLoss:        sdkmath.NewInt(0),
			CurrentRoundMaxLossOddsUID: "6db09053-2901-4110-8fb5-c14e21f8d666",
			ActualProfit:               sdkmath.NewInt(0),
			IsSettled:                  false,
		}
		nullify.Fill(&participation)

		state.OrderBookParticipationList = append(state.OrderBookParticipationList, participation)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), state.OrderBookParticipationList
}

func TestQueryParticipation(t *testing.T) {
	net, objs := networkWithParticipationObjects(t, 5)

	t.Run("ShowParticipation", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		common := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		for _, tc := range []struct {
			desc               string
			orderBookUID       string
			participationIndex uint64

			args []string
			err  error
			obj  types.OrderBookParticipation
		}{
			{
				desc:               "found",
				orderBookUID:       testMarketUID,
				participationIndex: objs[0].Index,

				args: common,
				obj:  objs[0],
			},
			{
				desc:               "not found",
				participationIndex: 10000,
				orderBookUID:       cast.ToString(100000),

				args: common,
				err:  status.Error(codes.NotFound, "not found"),
			},
		} {
			tc := tc
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					tc.orderBookUID,
					cast.ToString(tc.participationIndex),
				}
				args = append(args, tc.args...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookParticipation(), args)
				if tc.err != nil {
					stat, ok := status.FromError(tc.err)
					require.True(t, ok)
					require.ErrorIs(t, stat.Err(), tc.err)
				} else {
					require.NoError(t, err)
					var resp types.OrderBookParticipation
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

	t.Run("ListParticipation", func(t *testing.T) {
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookParticipations(), args)
				require.NoError(t, err)
				var resp types.QueryOrderBookParticipationsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.OrderBookParticipations), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.OrderBookParticipations),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookParticipations(), args)
				require.NoError(t, err)
				var resp types.QueryOrderBookParticipationsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.OrderBookParticipations), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.OrderBookParticipations),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookParticipations(), args)
			require.NoError(t, err)
			var resp types.QueryOrderBookParticipationsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.OrderBookParticipations),
			)
		})
	})
}
