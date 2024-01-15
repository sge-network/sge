package cli_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/x/market/client/cli"
	"github.com/sge-network/sge/x/market/types"
)

func TestQueryParams(t *testing.T) {
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
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			var args []string

			args = append(args, tc.args...)
			res, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryParams(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}

			var params types.QueryParamsResponse
			err = json.Unmarshal(res.Bytes(), &params)
			require.NoError(t, err)

			defaultParams := types.DefaultParams()
			require.Equal(t, types.QueryParamsResponse{
				Params: defaultParams,
			}, params)
		})
	}
}
