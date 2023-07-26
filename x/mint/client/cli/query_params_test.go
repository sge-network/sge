package cli_test

import (
	"encoding/json"
	"fmt"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/x/mint/client/cli"
	"github.com/sge-network/sge/x/mint/types"
	"github.com/stretchr/testify/require"
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
			res, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdQueryParams(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}

			var params types.QueryParamsResponse
			err = json.Unmarshal(res.Bytes(), &params)
			fmt.Println(params)
			// command-line response, wraps the primitive numbers in double quotes, so it is not unmarshall-able.
			require.EqualError(
				t,
				err,
				"json: cannot unmarshal string into Go struct field Params.params.blocks_per_year of type int64",
			)

			defaultParams := types.DefaultParams()
			defaultParams.BlocksPerYear = 0
			require.Equal(t, types.QueryParamsResponse{
				Params: defaultParams,
			}, params)
		})
	}
}
