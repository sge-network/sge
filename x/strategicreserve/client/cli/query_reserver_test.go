package cli_test

import (
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/x/strategicreserve/client/cli"
	"github.com/stretchr/testify/require"
)

func TestQueryReserver(t *testing.T) {
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
			desc: "Success! Reserver returned",
			args: nil,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			_, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryReserver(), tc.args)
			if tc.err != nil {
				require.EqualError(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
