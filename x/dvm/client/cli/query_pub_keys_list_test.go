package cli_test

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"
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

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithPublicKeys(t *testing.T) (*network.Network, *types.PublicKeys, *ed25519.PrivateKey) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	bs, err := x509.MarshalPKIXPublicKey(simappUtil.TestDVMPublicKey)
	if err != nil {
		panic(err)
	}
	state.PublicKeys = &types.PublicKeys{
		List: []string{string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: bs}))},
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	return network.New(t, cfg), state.PublicKeys, &simappUtil.TestDVMPrivateKey
}

func TestCmdPubKeysList(t *testing.T) {
	net, _, _ := networkWithPublicKeys(t)

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
					var resp types.QueryListPubKeyAllResponse
					require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
					require.True(t, len(resp.List) > 0)
					t.Log(resp.List)

				}
			})
		}
	})
}
