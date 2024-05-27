package client

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/testutil"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"

	"github.com/sge-network/sge/x/market/types"
)

func (s *E2ETestSuite) TestMarketsGRPCHandler() {
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
			"test GRPC Market by UID",
			fmt.Sprintf("%s/sge/market/%s", baseURL, dummyMarketUID),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QueryMarketResponse{},
			&types.QueryMarketResponse{
				Market: dummyMarket,
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
