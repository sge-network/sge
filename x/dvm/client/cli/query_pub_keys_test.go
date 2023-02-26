package cli_test

import (
	"fmt"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/sge-network/sge/testutil/network"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/dvm/client/cli"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/status"
)

func networkWithPublicKeys(t *testing.T) (*network.Network, *types.KeyVault) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	pubkeys := simappUtil.GenerateDvmPublicKeys(5)

	state.KeyVault = types.KeyVault{
		PublicKeys: pubkeys,
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), &state.KeyVault
}

func TestCmdPubKeysList(t *testing.T) {
	net, _ := networkWithPublicKeys(t)

	t.Run("PubKeysList", func(t *testing.T) {
		ctx := net.Validators[0].ClientCtx
		common := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		TestCases := []struct {
			desc string
			args []string
			err  error
		}{
			{
				desc: "success",
				args: common,
				err:  nil,
			},
		}
		for _, tc := range TestCases {
			t.Run(tc.desc, func(t *testing.T) {
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdPubKeysList(), tc.args)
				if tc.err != nil {
					stat, ok := status.FromError(tc.err)
					require.True(t, ok)
					require.Error(t, stat.Err(), "")
				} else {
					require.NoError(t, err)
					var resp types.QueryPubKeysResponse
					require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
					require.True(t, len(resp.List) > 0)
					t.Log(resp.List)
				}
			})
		}
	})
}
