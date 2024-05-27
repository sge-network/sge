package client

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	e2esdk "github.com/sge-network/sge/tests/e2e/sdk"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/market/client/cli"
	"github.com/sge-network/sge/x/market/types"
)

var (
	genesis types.GenesisState

	dummyMarket    = types.Market{}
	dummyMarketUID = uuid.NewString()

	oddsUID1   = "9991c60f-2025-48ce-ae79-1dc110f16900"
	oddsValue1 = sdkmath.LegacyNewDecWithPrec(14, 1) // 1.4
	oddsUID2   = "9991c60f-2025-48ce-ae79-1dc110f16901"
	oddsValue2 = sdkmath.LegacyNewDecWithPrec(13, 1) // 1.3
	oddsUID3   = "9991c60f-2025-48ce-ae79-1dc110f16902"

	fees            = int64(10)
	marketUID       = uuid.NewString()
	marketStartDate = uint64(time.Now().Unix())
	marketEndDate   = uint64(time.Now().Add(time.Minute * time.Duration(120)).Unix())
)

type E2ETestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewE2ETestSuite(cfg network.Config) *E2ETestSuite {
	return &E2ETestSuite{cfg: cfg}
}

func (s *E2ETestSuite) SetupSuite() {
	s.T().Log("setting up e2e test suite")

	genesisState := s.cfg.GenesisState

	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[types.ModuleName], &genesis))

	dummyMarket = types.Market{
		Creator: sample.AccAddress(),
		UID:     dummyMarketUID,
		StartTS: uint64(time.Now().Unix()),
		EndTS:   uint64(time.Now().Add(time.Minute * time.Duration(120)).Unix()),
		Odds: []*types.Odds{
			{
				UID:  "9991c60f-2025-48ce-ae79-1dc110f16904",
				Meta: "Home",
			},
			{
				UID:  "9991c60f-2025-48ce-ae79-1dc110f16905",
				Meta: "Draw",
			},
			{
				UID:  "9991c60f-2025-48ce-ae79-1dc110f16906",
				Meta: "Away",
			},
		},
		WinnerOddsUIDs: []string{},
		Status:         types.MarketStatus_MARKET_STATUS_ACTIVE,
		Meta:           "dummy market",
		BookUID:        dummyMarketUID,
	}
	genesis.MarketList = []types.Market{dummyMarket}

	genesis.Stats = types.MarketStats{
		// ResolvedUnsettled: []string{dummyMarketUID},
	}

	// genesis.Params = types.Params{}

	genesisBz, err := s.cfg.Codec.MarshalJSON(&genesis)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = genesisBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)

	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *E2ETestSuite) TestMarketAddTxCmd() {
	val := s.network.Validators[0]

	clientCtx := val.ClientCtx

	valBalance := e2esdk.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(sdkmath.NewInt(400000000), valBalance)

	testCases := []struct {
		name           string
		ticketClaims   jwt.MapClaims
		args           []string
		expectErr      bool
		expectedCode   uint32
		expectedErrMsg string
		respType       proto.Message
	}{
		{
			"invalid transaction for val",
			jwt.MapClaims{
				"exp":      time.Now().Add(time.Minute * 5).Unix(),
				"iat":      time.Now().Unix(),
				"uid":      marketUID,
				"start_ts": marketStartDate,
				"end_ts":   marketEndDate,
				"odds":     []types.Odds{},
				"status":   types.MarketStatus_MARKET_STATUS_ACTIVE,
				"meta":     "sample 3way market",
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagGas, flags.GasFlagAuto),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(fees))).String()),
			},
			true,
			0,
			"",
			&sdk.TxResponse{},
		},
		{
			"valid transaction for val",
			jwt.MapClaims{
				"exp":      time.Now().Add(time.Minute * 5).Unix(),
				"iat":      time.Now().Unix(),
				"uid":      marketUID,
				"start_ts": marketStartDate,
				"end_ts":   marketEndDate,
				"odds": []types.Odds{
					{
						UID:  oddsUID1,
						Meta: "Home",
					},
					{
						UID:  oddsUID2,
						Meta: "Draw",
					},
					{
						UID:  oddsUID3,
						Meta: "Away",
					},
				},
				"status": types.MarketStatus_MARKET_STATUS_ACTIVE,
				"meta":   "sample 3way market",
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagGas, flags.GasFlagAuto),
				fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.15"),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(fees))).String()),
			},
			false,
			0,
			"",
			&sdk.TxResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Require().NoError(s.network.WaitForNextBlock())
		s.Run(tc.name, func() {
			ticket, err := simapp.CreateJwtTicket(tc.ticketClaims)
			require.Nil(s.T(), err)

			tc.args = append([]string{ticket}, tc.args...)
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdAdd(), tc.args)

			if tc.expectErr {
				s.Require().Error(err)
				s.T().Logf("error captured: %s", tc.name)
			} else {
				s.Require().NoError(err)

				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), tc.respType), bz.String())
				respType, ok := tc.respType.(*sdk.TxResponse)
				s.Require().True(ok)
				s.Require().Equal(tc.expectedCode, respType.Code)

				txResp, err := clitestutil.GetTxResponse(s.network, clientCtx, respType.TxHash)
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedCode, txResp.Code)
				if tc.expectedErrMsg != "" {
					s.Require().Contains(txResp.RawLog, tc.expectedErrMsg)
				}
			}
		})
	}
}
