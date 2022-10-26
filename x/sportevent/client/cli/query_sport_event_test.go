package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/sportevent/client/cli"
	"github.com/sge-network/sge/x/sportevent/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithSportEventObjects(t *testing.T, n int) (*network.Network, []types.SportEvent) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		sportEvent := types.SportEvent{
			UID:            strconv.Itoa(i),
			WinnerOddsUIDs: map[string][]byte{},
		}
		nullify.Fill(&sportEvent)
		state.SportEventList = append(state.SportEventList, sportEvent)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.SportEventList
}

func TestQuerySportEventCLI(t *testing.T) {
	net, objs := networkWithSportEventObjects(t, 5)
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}

	t.Run("ShowSportEvent", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		for _, tc := range []struct {
			desc  string
			idUID string

			args []string
			err  error
			obj  types.SportEvent
		}{
			{
				desc:  "found",
				idUID: objs[0].UID,

				args: common,
				obj:  objs[0],
			},
			{
				desc:  "not found",
				idUID: strconv.Itoa(100000),

				args: common,
				err:  status.Error(codes.NotFound, "not found"),
			},
		} {
			tc := tc
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					tc.idUID,
				}
				args = append(args, tc.args...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowSportEvent(), args)
				if tc.err != nil {
					stat, ok := status.FromError(tc.err)
					require.True(t, ok)
					require.ErrorIs(t, stat.Err(), tc.err)
				} else {
					require.NoError(t, err)
					var resp types.QuerySportEventResponse
					require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
					require.NotNil(t, resp.SportEvent)
					require.Equal(t,
						nullify.Fill(&tc.obj),
						nullify.Fill(&resp.SportEvent),
					)
				}
			})
		}
	})

	t.Run("ListSportEvent", func(t *testing.T) {
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListSportEvent(), args)
				require.NoError(t, err)
				var resp types.QuerySportEventListAllResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.SportEvent), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.SportEvent),
				)
			}
		})
		t.Run("ByKey", func(t *testing.T) {
			step := 2
			var next []byte
			for i := 0; i < len(objs); i += step {
				args := request(next, 0, uint64(step), false)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListSportEvent(), args)
				require.NoError(t, err)
				var resp types.QuerySportEventListAllResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.LessOrEqual(t, len(resp.SportEvent), step)
				require.Subset(t,
					nullify.Fill(objs),
					nullify.Fill(resp.SportEvent),
				)
				next = resp.Pagination.NextKey
			}
		})
		t.Run("Total", func(t *testing.T) {
			args := request(nil, 0, uint64(len(objs)), true)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListSportEvent(), args)
			require.NoError(t, err)
			var resp types.QuerySportEventListAllResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.NoError(t, err)
			require.Equal(t, len(objs), int(resp.Pagination.Total))
			require.ElementsMatch(t,
				nullify.Fill(objs),
				nullify.Fill(resp.SportEvent),
			)
		})
	})
}
