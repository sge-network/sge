package client

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/testutil"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sge-network/sge/x/house/types"
)

func (s *E2ETestSuite) TestDepositsGRPCHandler() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress

	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			"test GRPC all deposits",
			fmt.Sprintf("%s/sge/house/deposits", baseURL),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QueryDepositsResponse{},
			&types.QueryDepositsResponse{
				Deposits: genesis.DepositList,
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			"test GRPC deposits of certain creator",
			fmt.Sprintf("%s/sge/house/deposits/%s", baseURL, dummyDepositor),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QueryDepositsByAccountResponse{},
			&types.QueryDepositsByAccountResponse{
				Deposits: genesis.DepositList,
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}
