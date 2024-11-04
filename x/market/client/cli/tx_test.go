package cli_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/sge-network/sge/testutil/network"

	"github.com/sge-network/sge/x/market/client/cli"
)

func TestGetTxCmd(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	for _, tc := range []struct {
		desc string
		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			var args []string

			args = append(args, tc.args...)
			res, err := clitestutil.ExecTestCLICmd(ctx, cli.GetTxCmd(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}

			require.True(t, strings.HasPrefix(string(res.Bytes()), "market transactions subcommands"))
		})
	}
}
