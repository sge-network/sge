package cli_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"

	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/testutil/simapp"

	bettypes "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/orderbook/client/cli"
	"github.com/sge-network/sge/x/orderbook/types"
)

func networkWithParticipationBetObjects(t *testing.T, n int) (*network.Network, []types.ParticipationBetPair) {
	t.Helper()
	cfg := network.DefaultConfig()

	// bet module state
	betState := bettypes.GenesisState{}
	betState.BetList = []bettypes.Bet{}
	betState.Uid2IdList = []bettypes.UID2ID{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[bettypes.ModuleName], &betState))

	// orderbook module state
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		bet := bettypes.Bet{
			UID:     uuid.NewString(),
			Creator: simapp.TestParamUsers["user"+cast.ToString(i+1)].Address.String(),
		}
		nullify.Fill(&bet)
		betState.BetList = append(betState.BetList, bet)
		betState.Uid2IdList = append(betState.Uid2IdList, bettypes.UID2ID{UID: bet.UID, ID: cast.ToUint64(i + 1)})

		betPair := types.ParticipationBetPair{
			OrderBookUID:       testMarketUID,
			ParticipationIndex: testParticipationIndex,
			BetUID:             bet.UID,
		}
		nullify.Fill(&betPair)

		state.ParticipationBetPairExposureList = append(state.ParticipationBetPairExposureList, betPair)
	}

	betBuf, err := cfg.Codec.MarshalJSON(&betState)
	require.NoError(t, err)
	cfg.GenesisState[bettypes.ModuleName] = betBuf

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), state.ParticipationBetPairExposureList
}

func TestQueryParticipationBetPair(t *testing.T) {
	net, objs := networkWithParticipationBetObjects(t, 5)

	t.Run("ListParticipation", func(t *testing.T) {
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryParticipationBets(), args)
				require.NoError(t, err)
				var resp types.QueryParticipationFulfilledBetsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.ParticipationBets), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.ParticipationBets),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryParticipationBets(), args)
				require.NoError(t, err)
				var resp types.QueryParticipationFulfilledBetsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.ParticipationBets), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.ParticipationBets),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryParticipationBets(), args)
			require.NoError(t, err)
			var resp types.QueryParticipationFulfilledBetsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.ParticipationBets),
			)
		})
	})
}
