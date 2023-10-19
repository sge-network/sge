package cli_test

import (
	"fmt"
	"testing"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/orderbook/client/cli"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func networkWithParticipationExposureObjects(t *testing.T, n int) (*network.Network, []types.ParticipationExposure) {
	t.Helper()
	cfg := network.DefaultConfig()

	// orderbook module state
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		exposure := types.ParticipationExposure{
			OddsUID:            uuid.NewString(),
			ParticipationIndex: testParticipationIndex,
			OrderBookUID:       testMarketUID,
			Exposure:           sdk.NewInt(1000),
			BetAmount:          sdk.NewInt(1000),
			IsFulfilled:        false,
			Round:              1,
		}
		nullify.Fill(&exposure)

		state.ParticipationExposureList = append(state.ParticipationExposureList, exposure)
		state.HistoricalParticipationExposureList = append(state.HistoricalParticipationExposureList, exposure)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), state.ParticipationExposureList
}

func TestQueryParticipationExposure(t *testing.T) {
	net, objs := networkWithParticipationExposureObjects(t, 5)

	t.Run("ListOrderBookParticipationExposures", func(t *testing.T) {
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookParticipationExposures(), args)
				require.NoError(t, err)
				var resp types.QueryOrderBookParticipationExposuresResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.ParticipationExposures), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.ParticipationExposures),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookParticipationExposures(), args)
				require.NoError(t, err)
				var resp types.QueryOrderBookParticipationExposuresResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.ParticipationExposures), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.ParticipationExposures),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryOrderBookParticipationExposures(), args)
			require.NoError(t, err)
			var resp types.QueryOrderBookParticipationExposuresResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.ParticipationExposures),
			)
		})
	})

	t.Run("ListParticipationExposures", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		request := func(next []byte, offset, limit uint64, total bool) []string {
			args := []string{
				testMarketUID,
				cast.ToString(testParticipationIndex),
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryParticipationExposures(), args)
				require.NoError(t, err)
				var resp types.QueryParticipationExposuresResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.ParticipationExposures), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.ParticipationExposures),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryParticipationExposures(), args)
				require.NoError(t, err)
				var resp types.QueryParticipationExposuresResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.ParticipationExposures), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.ParticipationExposures),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryParticipationExposures(), args)
			require.NoError(t, err)
			var resp types.QueryParticipationExposuresResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.ParticipationExposures),
			)
		})
	})

	t.Run("ListHistoricalParticipationExposures", func(t *testing.T) {
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryHistoricalParticipationExposures(), args)
				require.NoError(t, err)
				var resp types.QueryHistoricalParticipationExposuresResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.ParticipationExposures), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.ParticipationExposures),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryHistoricalParticipationExposures(), args)
				require.NoError(t, err)
				var resp types.QueryHistoricalParticipationExposuresResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.ParticipationExposures), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.ParticipationExposures),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryHistoricalParticipationExposures(), args)
			require.NoError(t, err)
			var resp types.QueryHistoricalParticipationExposuresResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.ParticipationExposures),
			)
		})
	})
}
