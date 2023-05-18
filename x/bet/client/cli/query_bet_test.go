package cli_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/nullify"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/client/cli"
	"github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

const testMarketUID = "5db09053-2901-4110-8fb5-c14e21f8d555"

func networkWithBetObjects(t *testing.T, n int) (*network.Network, []types.Bet) {
	t.Helper()
	cfg := network.DefaultConfig()

	// market module state
	marketState := markettypes.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[markettypes.ModuleName], &marketState))

	market := markettypes.Market{
		UID:     testMarketUID,
		Creator: simappUtil.TestParamUsers["user1"].Address.String(),
		StartTS: 1111111111,
		EndTS:   uint64(time.Now().Unix()) + 5000,
		Odds: []*markettypes.Odds{
			{UID: "6db09053-2901-4110-8fb5-c14e21f8d666", Meta: "Odds 1"},
			{UID: "5e31c60f-2025-48ce-ae79-1dc110f16358", Meta: "Odds 2"},
			{UID: "6e31c60f-2025-48ce-ae79-1dc110f16354", Meta: "Odds 3"},
		},
		Status: markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
	}
	nullify.Fill(&market)
	marketState.MarketList = []markettypes.Market{market}

	marketBuf, err := cfg.Codec.MarshalJSON(&marketState)
	require.NoError(t, err)
	cfg.GenesisState[markettypes.ModuleName] = marketBuf

	// bet module state
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		bet := types.Bet{
			Creator:           testAddress,
			UID:               uuid.NewString(),
			MarketUID:         market.UID,
			OddsValue:         "10",
			Amount:            sdk.NewInt(10),
			BetFee:            sdk.NewInt(1),
			MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
		}
		nullify.Fill(&bet)

		state.BetList = append(state.BetList, bet)
		state.PendingBetList = append(state.PendingBetList, types.PendingBet{UID: bet.UID, Creator: testAddress})
		state.SettledBetList = append(state.SettledBetList, types.SettledBet{UID: bet.UID, BettorAddress: testAddress})

		id := uint64(i + 1)
		state.Uid2IdList = append(state.Uid2IdList, types.UID2ID{UID: bet.UID, ID: id})
	}
	state.Stats = types.BetStats{Count: 5}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), state.BetList
}

func TestQueryBet(t *testing.T) {
	net, objs := networkWithBetObjects(t, 5)

	t.Run("ShowBet", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		common := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		for _, tc := range []struct {
			desc    string
			uid     string
			creator string

			args []string
			err  error
			obj  types.Bet
		}{
			{
				desc:    "found",
				creator: testAddress,
				uid:     objs[0].UID,

				args: common,
				obj:  objs[0],
			},
			{
				desc:    "not found",
				creator: "",
				uid:     cast.ToString(100000),

				args: common,
				err:  status.Error(codes.NotFound, "not found"),
			},
		} {
			tc := tc
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					tc.creator,
					tc.uid,
				}
				args = append(args, tc.args...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowBet(), args)
				if tc.err != nil {
					stat, ok := status.FromError(tc.err)
					require.True(t, ok)
					require.ErrorIs(t, stat.Err(), tc.err)
				} else {
					require.NoError(t, err)
					var resp types.QueryBetResponse
					require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
					require.NotNil(t, resp.Bet)
					require.Equal(t,
						nullify.Fill(&tc.obj),
						nullify.Fill(&resp.Bet),
					)
				}
			})
		}
	})

	t.Run("ListBet", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		request := func(next []byte, offset, limit uint64, total bool) []string {
			args := []string{
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListBet(), args)
				require.NoError(t, err)
				var resp types.QueryBetsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Bet), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Bet),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListBet(), args)
				require.NoError(t, err)
				var resp types.QueryBetsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Bet), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Bet),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListBet(), args)
			require.NoError(t, err)
			var resp types.QueryBetsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Bet),
			)
		})
	})

	t.Run("ListBetByCreator", func(t *testing.T) {
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListBetByCreator(), args)
				require.NoError(t, err)
				var resp types.QueryBetsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Bet), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Bet),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListBetByCreator(), args)
				require.NoError(t, err)
				var resp types.QueryBetsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Bet), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Bet),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListBetByCreator(), args)
			require.NoError(t, err)
			var resp types.QueryBetsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Bet),
			)
		})
	})

	t.Run("ListPendingBetOfMarket", func(t *testing.T) {
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPendingBets(), args)
				require.NoError(t, err)
				var resp types.QueryPendingBetsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Bet), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Bet),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPendingBets(), args)
				require.NoError(t, err)
				var resp types.QueryPendingBetsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Bet), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Bet),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListPendingBets(), args)
			require.NoError(t, err)
			var resp types.QueryPendingBetsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Bet),
			)
		})
	})

	t.Run("ListBetByUIDs", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		common := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		var items []string
		for _, v := range objs {
			items = append(items, testAddress+":"+v.UID)
		}

		for _, tc := range []struct {
			desc  string
			items []string

			args []string
			err  error
			obj  []types.Bet
		}{
			{
				desc:  "found",
				items: items,

				args: common,
				obj:  objs,
			},
		} {
			tc := tc
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					strings.Join(tc.items, ","),
				}
				args = append(args, tc.args...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListBetByUIDs(), args)
				if tc.err != nil {
					stat, ok := status.FromError(tc.err)
					require.True(t, ok)
					require.ErrorIs(t, stat.Err(), tc.err)
				} else {
					require.NoError(t, err)
					var resp types.QueryBetsByUIDsResponse
					require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
					require.NotNil(t, resp.Bets)
					require.Equal(t,
						nullify.Fill(&tc.obj),
						nullify.Fill(&resp.Bets),
					)
				}
			})
		}
	})
}
