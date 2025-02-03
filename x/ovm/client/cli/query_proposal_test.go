package cli_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	tmcli "github.com/cometbft/cometbft/libs/cli"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/ovm/client/cli"
	"github.com/sge-network/sge/x/ovm/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithPubkeysChangeProposalObjects(
	t *testing.T,
	n int,
) (*network.Network, []types.PublicKeysChangeProposal) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		proposal := types.PublicKeysChangeProposal{
			Id: uint64(i),
			Modifications: types.PubkeysChangeProposalPayload{
				PublicKeys:  []string{},
				LeaderIndex: 0,
			},
			Votes:   []*types.Vote{},
			StartTS: time.Now().Unix(),
		}
		nullify.Fill(&proposal)
		state.PubkeysChangeProposals = append(state.PubkeysChangeProposals, proposal)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.PubkeysChangeProposals
}

func TestCmdQueryPubkeysChangeProposal(t *testing.T) {
	net, objs := networkWithPubkeysChangeProposalObjects(t, 5)
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}

	t.Run("ShowActiveProposal", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		for _, tc := range []struct {
			desc   string
			id     string
			status types.ProposalStatus

			args []string
			err  error
			obj  types.PublicKeysChangeProposal
		}{
			{
				desc:   "found",
				id:     cast.ToString(objs[0].Id),
				status: types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,

				args: common,
				obj:  objs[0],
			},
			{
				desc:   "not found",
				id:     cast.ToString(100000),
				status: types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,

				args: common,
				err:  status.Error(codes.NotFound, "not found"),
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				args := []string{
					cast.ToString(cast.ToInt32(tc.status)),
					tc.id,
				}
				args = append(args, tc.args...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdPubkeysChangeProposal(), args)
				if tc.err != nil {
					stat, ok := status.FromError(tc.err)
					require.True(t, ok)
					require.ErrorIs(t, stat.Err(), tc.err)
				} else {
					require.NoError(t, err)
					var resp types.QueryPublicKeysChangeProposalResponse
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
