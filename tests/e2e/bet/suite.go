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
	sgetypes "github.com/sge-network/sge/types"
	"github.com/sge-network/sge/x/bet/client/cli"
	"github.com/sge-network/sge/x/bet/types"
	housecli "github.com/sge-network/sge/x/house/client/cli"
	marketcli "github.com/sge-network/sge/x/market/client/cli"
	markettypes "github.com/sge-network/sge/x/market/types"
)

var (
	genesis types.GenesisState

	dummyMarket     = markettypes.Market{}
	dummyBetCreator = sample.AccAddress()
	dummyBetUID     = uuid.NewString()
	dummyMarketUID  = uuid.NewString()

	oddsUID1          = "9991c60f-2025-48ce-ae79-1dc110f16900"
	oddsValue1        = sdkmath.LegacyNewDecWithPrec(14, 1) // 1.4
	oddsUID2          = "9991c60f-2025-48ce-ae79-1dc110f16901"
	oddsValue2        = sdkmath.LegacyNewDecWithPrec(13, 1) // 1.3
	oddsUID3          = "9991c60f-2025-48ce-ae79-1dc110f16902"
	maxLossMultiplier = sdkmath.LegacyNewDecWithPrec(10, 2)

	betFees         = sdkmath.NewInt(10)
	fees            = int64(10)
	marketUID       = uuid.NewString()
	marketStartDate = uint64(time.Now().Unix())
	marketEndDate   = uint64(time.Now().Add(time.Minute * time.Duration(120)).Unix())
	depositAmount   = sdkmath.NewInt(100000)

	wagerDate   = time.Now().Add(time.Minute * time.Duration(180)).Unix()
	wagerAmount = sdkmath.NewInt(100)

	acc1        = sdk.AccAddress{}
	acc1Balance = sdkmath.ZeroInt()

	acc2        = sdk.AccAddress{}
	acc2Balance = sdkmath.ZeroInt()
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

	var marketGenesis markettypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[markettypes.ModuleName], &marketGenesis))

	dummyMarket = markettypes.Market{
		Creator: sample.AccAddress(),
		UID:     dummyMarketUID,
		StartTS: uint64(time.Now().Unix()),
		EndTS:   uint64(time.Now().Add(time.Minute * time.Duration(120)).Unix()),
		Odds: []*markettypes.Odds{
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
		Status:         markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
		Meta:           "dummy market",
		BookUID:        dummyMarketUID,
	}
	marketGenesis.MarketList = []markettypes.Market{dummyMarket}
	marketGenesisBz, err := s.cfg.Codec.MarshalJSON(&marketGenesis)
	s.Require().NoError(err)
	genesisState[markettypes.ModuleName] = marketGenesisBz

	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[types.ModuleName], &genesis))

	genesis.BetList = []types.Bet{
		{
			Creator:           dummyBetCreator,
			UID:               dummyBetUID,
			OddsUID:           marketGenesis.MarketList[0].OddsUIDS()[0],
			OddsValue:         "0.5",
			Fee:               sdkmath.NewInt(10),
			Status:            types.Bet_STATUS_SETTLED,
			Result:            types.Bet_RESULT_WON,
			CreatedAt:         wagerDate,
			SettlementHeight:  1,
			MaxLossMultiplier: sdkmath.LegacyNewDecWithPrec(10, 2),
			Meta: types.MetaData{
				SelectedOddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
				SelectedOddsValue: "0.5",
				IsMainMarket:      false,
			},
			BetFulfillment: []*types.BetFulfillment{},
			MarketUID:      dummyMarketUID,
			Amount:         wagerAmount,
		},
	}

	genesis.Uid2IdList = []types.UID2ID{
		{
			UID: dummyBetUID,
			ID:  1,
		},
	}

	genesis.PendingBetList = []types.PendingBet{
		{
			UID:     dummyBetUID,
			Creator: dummyBetCreator,
		},
	}

	genesis.SettledBetList = []types.SettledBet{
		{
			UID:           dummyBetUID,
			BettorAddress: dummyBetCreator,
		},
	}

	genesis.Stats = types.BetStats{
		Count: 1,
	}

	genesis.Params = types.Params{
		BatchSettlementCount:  uint32(1000),
		MaxBetByUidQueryCount: uint32(1000),
		Constraints: types.Constraints{
			MinAmount: sdkmath.NewInt(100),
			Fee:       betFees,
		},
	}

	genesisBz, err := s.cfg.Codec.MarshalJSON(&genesis)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = genesisBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)

	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *E2ETestSuite) TestWagerTxCmd() {
	val := s.network.Validators[0]

	clientCtx := val.ClientCtx
	{
		ticket, err := simapp.CreateJwtTicket(jwt.MapClaims{
			"exp":      time.Now().Add(time.Minute * 5).Unix(),
			"iat":      time.Now().Unix(),
			"uid":      marketUID,
			"start_ts": marketStartDate,
			"end_ts":   marketEndDate,
			"odds": []markettypes.Odds{
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
			"status": markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
			"meta":   "sample 3way market",
		})
		require.Nil(s.T(), err)

		args := []string{
			ticket,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(fees))).String()),
		}
		bz, err := clitestutil.ExecTestCLICmd(clientCtx, marketcli.CmdAdd(), args)
		s.Require().NoError(err)
		respType := sdk.TxResponse{}
		s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), &respType), bz.String())
		s.Require().Equal(uint32(0), respType.Code)
	}
	s.Require().NoError(s.network.WaitForNextBlock())

	{
		ticket, err := simapp.CreateJwtTicket(jwt.MapClaims{
			"exp": time.Now().Add(time.Minute * 5).Unix(),
			"iat": time.Now().Unix(),
			"kyc_data": sgetypes.KycDataPayload{
				Ignore:   false,
				Approved: true,
				ID:       val.Address.String(),
			},
			"depositor_address": val.Address.String(),
		})
		require.Nil(s.T(), err)

		bz, err := clitestutil.ExecTestCLICmd(clientCtx, housecli.CmdDeposit(), []string{
			marketUID,
			depositAmount.String(),
			ticket,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			fmt.Sprintf("--%s=%s", flags.FlagGas, flags.GasFlagAuto),
			fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.15"),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(fees))).String()),
		})
		s.Require().NoError(err)

		var respType sdk.TxResponse
		s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), &respType), bz.String())
		s.Require().Equal(uint32(0), respType.Code)

		txResp, err := clitestutil.GetTxResponse(s.network, clientCtx, respType.TxHash)
		s.Require().NoError(err)
		s.Require().Equal(uint32(0), txResp.Code)
	}

	acc1initialBalance := 1000
	acc1 = e2esdk.CreateAccount(val, "acc1")
	e2esdk.SendToken(val, acc1, acc1initialBalance)
	s.Require().NoError(s.network.WaitForNextBlock())
	acc1Balance = e2esdk.GetSGEBalance(clientCtx, acc1.String())

	acc2initialBalance := 1000
	acc2 = e2esdk.CreateAccount(val, "acc2")
	e2esdk.SendToken(val, acc2, acc2initialBalance)
	s.Require().NoError(s.network.WaitForNextBlock())
	acc2Balance = e2esdk.GetSGEBalance(clientCtx, acc2.String())

	testCases := []struct {
		name           string
		uid            string
		marketUID      string
		amount         sdkmath.Int
		ticketClaims   jwt.MapClaims
		args           []string
		expectErr      bool
		expectedCode   uint32
		expectedErrMsg string
		respType       proto.Message
	}{
		{
			"valid transaction for acc1",
			uuid.NewString(),
			marketUID,
			wagerAmount,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       acc1.String(),
				},
				"selected_odds": types.BetOdds{
					UID:               oddsUID1,
					MarketUID:         marketUID,
					Value:             oddsValue1.String(),
					MaxLossMultiplier: maxLossMultiplier,
				},
				"all_odds": []types.BetOddsCompact{
					{UID: oddsUID1, MaxLossMultiplier: maxLossMultiplier},
					{UID: oddsUID2, MaxLossMultiplier: maxLossMultiplier},
					{UID: oddsUID3, MaxLossMultiplier: maxLossMultiplier},
				},
				"meta": types.MetaData{
					SelectedOddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
					SelectedOddsValue: oddsValue1.String(),
					IsMainMarket:      false,
				},
				"context": "sample context",
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "acc1"),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagGas, flags.GasFlagAuto),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(fees))).String()),
			},
			false,
			0,
			"",
			&sdk.TxResponse{},
		},
		{
			"valid transaction for acc2",
			uuid.NewString(),
			marketUID,
			wagerAmount,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       acc2.String(),
				},
				"selected_odds": types.BetOdds{
					UID:               oddsUID2,
					MarketUID:         marketUID,
					Value:             oddsValue2.String(),
					MaxLossMultiplier: maxLossMultiplier,
				},
				"all_odds": []types.BetOddsCompact{
					{UID: oddsUID1, MaxLossMultiplier: maxLossMultiplier},
					{UID: oddsUID2, MaxLossMultiplier: maxLossMultiplier},
					{UID: oddsUID3, MaxLossMultiplier: maxLossMultiplier},
				},
				"meta": types.MetaData{
					SelectedOddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
					SelectedOddsValue: oddsValue2.String(),
					IsMainMarket:      false,
				},
				"context": "sample context",
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "acc2"),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagGas, flags.GasFlagAuto),
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

			tc.args = append([]string{
				tc.uid,
				tc.amount.String(),
				ticket,
			}, tc.args...)
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdWager(), tc.args)

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

	{
		ticket, err := simapp.CreateJwtTicket(jwt.MapClaims{
			"exp":              time.Now().Add(time.Minute * 5).Unix(),
			"iat":              time.Now().Unix(),
			"uid":              marketUID,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{oddsUID1},
			"status":           markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
		})
		require.Nil(s.T(), err)

		bz, err := clitestutil.ExecTestCLICmd(clientCtx, marketcli.CmdResolve(), []string{
			ticket,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			fmt.Sprintf("--%s=%s", flags.FlagGas, flags.GasFlagAuto),
			fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.15"),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(fees))).String()),
		})
		s.Require().NoError(err)

		var respType sdk.TxResponse
		s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), &respType), bz.String())
		s.Require().Equal(uint32(0), respType.Code)

		s.Require().NoError(s.network.WaitForNextBlock())

		txResp, err := clitestutil.GetTxResponse(s.network, clientCtx, respType.TxHash)
		s.Require().NoError(err)
		s.Require().Equal(uint32(0), txResp.Code)
		s.T().Logf("==== success deposit create test case finished")
	}

	bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdListBet(), []string{fmt.Sprintf("--%s=json", flags.FlagOutput)})
	s.Require().NoError(err)
	respType := types.QueryBetsResponse{}
	err = clientCtx.Codec.UnmarshalJSON(bz.Bytes(), &respType)
	s.Require().NoError(err)

	s.Require().Equal(3, len(respType.Bet))

	for _, b := range respType.Bet {
		s.Require().Equal(types.Bet_STATUS_SETTLED, b.Status)
	}

	payout1 := oddsValue1.Mul(wagerAmount.Sub(betFees).ToLegacyDec()).TruncateInt().Sub(wagerAmount)
	expectedAcc1Balance := acc1Balance.Add(payout1).Sub(sdkmath.NewInt(fees))
	acc1Balance = e2esdk.GetSGEBalance(clientCtx, acc1.String())
	s.Require().Equal(expectedAcc1Balance, acc1Balance)

	expectedAcc2Balance := acc2Balance.Sub(wagerAmount).Sub(sdkmath.NewInt(fees))
	acc2Balance = e2esdk.GetSGEBalance(clientCtx, acc2.String())
	s.Require().Equal(expectedAcc2Balance, acc2Balance)
}
