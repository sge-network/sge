package cli_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/dvm/client/cli"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithPubkeysChangeProposalObjects(t *testing.T, n int) (*network.Network, []types.PublicKeysChangeProposal) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		proposal := types.PublicKeysChangeProposal{
			Id: uint64(i),
			Modifications: types.PubkeysChangeProposalPayload{
				Additions: []string{},
				Deletions: []string{},
			},
			Votes:   []*types.Vote{},
			StartTS: time.Now().Unix(),
		}
		nullify.Fill(&proposal)
		state.ActivePubkeysChangeProposals = append(state.ActivePubkeysChangeProposals, proposal)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.ActivePubkeysChangeProposals
}

func TestCmdQueryPubkeysChangeProposal(t *testing.T) {
	net, objs := networkWithPubkeysChangeProposalObjects(t, 5)
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}

	t.Run("ShowActiveProposal", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		for _, tc := range []struct {
			desc string
			id   string

			args []string
			err  error
			obj  types.PublicKeysChangeProposal
		}{
			{
				desc: "found",
				id:   cast.ToString(objs[0].Id),

				args: common,
				obj:  objs[0],
			},
			{
				desc: "not found",
				id:   cast.ToString(100000),

				args: common,
				err:  status.Error(codes.NotFound, "not found"),
			},
		} {
			tc := tc
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					tc.id,
				}
				args = append(args, tc.args...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdActivePubkeysChangeProposal(), args)
				if tc.err != nil {
					stat, ok := status.FromError(tc.err)
					require.True(t, ok)
					require.ErrorIs(t, stat.Err(), tc.err)
				} else {
					require.NoError(t, err)
					var resp types.QueryActivePublicKeysChangeProposalResponse
					require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
					require.NotNil(t, resp.Proposal)
					require.Equal(t,
						nullify.Fill(&tc.obj),
						nullify.Fill(&resp.Proposal),
					)
				}
			})
		}
	})
}
