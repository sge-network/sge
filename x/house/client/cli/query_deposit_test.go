package cli_test

import (
	"fmt"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/house/client/cli"
	"github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

const testMarketUID = "5db09053-2901-4110-8fb5-c14e21f8d555"

func networkWithDepositObjects(t *testing.T, n int) (*network.Network, []types.Deposit, []types.Withdrawal) {
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

	// house module state
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		deposit := types.Deposit{
			Creator:               testAddress,
			DepositorAddress:      testAddress,
			MarketUID:             market.UID,
			ParticipationIndex:    cast.ToUint64(i + 1),
			Amount:                sdkmath.NewInt(10),
			WithdrawalCount:       0,
			TotalWithdrawalAmount: sdkmath.NewInt(0),
		}
		nullify.Fill(&deposit)

		withdrawal := types.Withdrawal{
			Creator:            testAddress,
			Address:            deposit.DepositorAddress,
			ID:                 cast.ToUint64(i + 1),
			MarketUID:          deposit.MarketUID,
			ParticipationIndex: deposit.ParticipationIndex,
			Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
			Amount:             deposit.Amount,
		}
		nullify.Fill(&withdrawal)

		state.DepositList = append(state.DepositList, deposit)
		state.WithdrawalList = append(state.WithdrawalList, withdrawal)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), state.DepositList, state.WithdrawalList
}

func TestQueryDeposit(t *testing.T) {
	net, objs, _ := networkWithDepositObjects(t, 5)

	t.Run("ListDeposits", func(t *testing.T) {
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryDeposits(), args)
				require.NoError(t, err)
				var resp types.QueryDepositsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Deposits), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Deposits),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryDeposits(), args)
				require.NoError(t, err)
				var resp types.QueryDepositsResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Deposits), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Deposits),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryDeposits(), args)
			require.NoError(t, err)
			var resp types.QueryDepositsResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Deposits),
			)
		})
	})

	t.Run("ListDepositsByDepositor", func(t *testing.T) {
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryDepositsByAccount(), args)
				require.NoError(t, err)
				var resp types.QueryDepositsByAccountResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Deposits), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Deposits),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryDepositsByAccount(), args)
				require.NoError(t, err)
				var resp types.QueryDepositsByAccountResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.Deposits), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.Deposits),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryDepositsByAccount(), args)
			require.NoError(t, err)
			var resp types.QueryDepositsByAccountResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.Deposits),
			)
		})
	})
}
