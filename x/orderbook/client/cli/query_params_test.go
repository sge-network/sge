package cli_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"

	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/x/orderbook/client/cli"
	"github.com/sge-network/sge/x/orderbook/types"
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
			// command-line response, wraps the primitive numbers in double quotes, so it is not unmarshall-able.
			require.EqualError(
				t,
				err,
				"json: cannot unmarshal string into Go struct field Params.params.max_order_book_participations of type uint64",
				"json: cannot unmarshal string into Go struct field Params.params.batch_settlement_count of type uint64",
				"json: cannot unmarshal string into Go struct field Params.params.requeue_threshold of type uint64",
			)

			defaultParams := types.DefaultParams()
			defaultParams.BatchSettlementCount = 0
			defaultParams.MaxOrderBookParticipations = 0
			defaultParams.RequeueThreshold = 0
			require.Equal(t, types.QueryParamsResponse{
				Params: defaultParams,
			}, params)
		})
	}
}
