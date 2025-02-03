package client

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	"github.com/cosmos/cosmos-sdk/testutil"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sge-network/sge/x/bet/types"
)

func (s *E2ETestSuite) TestBetsGRPCHandler() {
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
			"test GRPC Bet by UID",
			fmt.Sprintf("%s/sge/bet/%s/%s", baseURL, dummyBetCreator, dummyBetUID),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QueryBetResponse{},
			&types.QueryBetResponse{
				Bet:    genesis.BetList[0],
				Market: dummyMarket,
			},
		},
		{
			"test GRPC sorted Bet by creator",
			fmt.Sprintf("%s/sge/bet/creator/%s/bets", baseURL, dummyBetCreator),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QueryBetsByCreatorResponse{},
			&types.QueryBetsByCreatorResponse{
				Bet: genesis.BetList,
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			"test GRPC Bet by UID",
			fmt.Sprintf("%s/sge/bet/bets", baseURL),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QueryBetsResponse{},
			&types.QueryBetsResponse{
				Bet: genesis.BetList,
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			"test GRPC pending bets by market",
			fmt.Sprintf("%s/sge/bet/bets/pending/%s", baseURL, dummyMarketUID),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QueryPendingBetsResponse{},
			&types.QueryPendingBetsResponse{
				Bet: genesis.BetList,
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			"test GRPC settled bets by market",
			fmt.Sprintf("%s/sge/bet/bets/settled/%d", baseURL, 1),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&types.QuerySettledBetsOfHeightResponse{},
			&types.QuerySettledBetsOfHeightResponse{
				Bet: genesis.BetList,
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			s.Require().NoError(err)

			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}
