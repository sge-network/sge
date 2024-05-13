package client

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"

	e2esdk "github.com/sge-network/sge/tests/e2e/sdk"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/testutil/simapp"
	sgetypes "github.com/sge-network/sge/types"
	"github.com/sge-network/sge/x/house/client/cli"
	"github.com/sge-network/sge/x/house/types"
	marketcli "github.com/sge-network/sge/x/market/client/cli"
	markettypes "github.com/sge-network/sge/x/market/types"
)

var (
	genesis        types.GenesisState
	dummyDepositor = sample.AccAddress()

	fees            = int64(10)
	marketUID       = uuid.NewString()
	marketStartDate = uint64(time.Now().Unix())
	marketEndDate   = uint64(time.Now().Add(time.Minute * time.Duration(120)).Unix())
	depositAmount   = sdkmath.NewInt(100)

	valBalance = sdkmath.NewInt(400000000)

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

	dummyMarketUID := uuid.NewString()
	marketGenesis.MarketList = []markettypes.Market{
		{
			Creator: sample.AccAddress(),
			UID:     dummyMarketUID,
			StartTS: uint64(time.Now().Unix()),
			EndTS:   uint64(time.Now().Add(time.Minute * time.Duration(120)).Unix()),
			Odds: []*markettypes.Odds{
				{
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16900",
					Meta: "Home",
				},
				{
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16901",
					Meta: "Draw",
				},
				{
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16901",
					Meta: "Away",
				},
			},
			WinnerOddsUIDs: []string{},
			Status:         markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
			Meta:           "dummy market",
			BookUID:        dummyMarketUID,
		},
	}
	marketGenesisBz, err := s.cfg.Codec.MarshalJSON(&marketGenesis)
	s.Require().NoError(err)
	genesisState[markettypes.ModuleName] = marketGenesisBz

	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[types.ModuleName], &genesis))

	genesis.DepositList = []types.Deposit{
		{
			Creator:               sample.AccAddress(),
			DepositorAddress:      dummyDepositor,
			MarketUID:             marketUID,
			ParticipationIndex:    0,
			Amount:                sdkmath.NewInt(10000),
			WithdrawalCount:       0,
			TotalWithdrawalAmount: sdkmath.ZeroInt(),
		},
	}

	genesis.WithdrawalList = []types.Withdrawal{
		{
			ID:                 1,
			Creator:            sample.AccAddress(),
			Address:            dummyDepositor,
			MarketUID:          marketUID,
			ParticipationIndex: 0,
			Amount:             sdkmath.NewInt(10000),
			Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
		},
	}

	genesis.Params = types.Params{
		MinDeposit:            sdkmath.NewInt(100),
		HouseParticipationFee: sdk.NewDecWithPrec(0, 2),
		MaxWithdrawalCount:    2,
	}

	genesisBz, err := s.cfg.Codec.MarshalJSON(&genesis)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = genesisBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)

	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *E2ETestSuite) TestDepositTxCmd() {
	s.T().Log("==== new house deposit command test started")
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
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16990",
					Meta: "Home",
				},
				{
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16991",
					Meta: "Draw",
				},
				{
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16992",
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

	expectedValBalance := valBalance.Sub(sdkmath.NewInt(fees))
	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(expectedValBalance, valBalance)

	acc1initialBalance := 1000
	acc1 = e2esdk.CreateAccount(val, "acc1")
	e2esdk.SendToken(val, acc1, acc1initialBalance)
	s.Require().NoError(s.network.WaitForNextBlock())

	expectedValBalance = valBalance.Sub(sdkmath.NewInt(int64(acc1initialBalance))).Sub(sdkmath.NewInt(fees))
	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(expectedValBalance, valBalance)

	testCases := []struct {
		name           string
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
			"valid transaction for val",
			marketUID,
			depositAmount,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       val.Address.String(),
				},
				"depositor_address": val.Address.String(),
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
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
		{
			"invalid transaction for acc1",
			marketUID,
			sdkmath.NewInt(100000000),
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       acc1.String(),
				},
				"depositor_address": acc1.String(),
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "acc1"),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagGas, flags.GasFlagAuto),
				fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.15"),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(fees))).String()),
			},
			true,
			sdkerrtypes.ErrInsufficientFunds.ABCICode(),
			"",
			&sdk.TxResponse{},
		},
		{
			"valid transaction for acc1",
			marketUID,
			depositAmount,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       acc1.String(),
				},
				"depositor_address": acc1.String(),
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "acc1"),
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

			tc.args = append([]string{
				tc.marketUID,
				tc.amount.String(),
				ticket,
			}, tc.args...)
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdDeposit(), tc.args)

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
				s.T().Logf("==== success deposit create test case finished: %s", tc.name)
			}
		})
	}

	depositFeeAmount := sdk.NewDecFromInt(depositAmount).Mul(genesis.Params.HouseParticipationFee).TruncateInt()

	expectedValBalance = valBalance.Sub(depositAmount).Sub(depositFeeAmount).Sub(sdkmath.NewInt(fees))
	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(expectedValBalance, valBalance)
	s.T().Logf("==== bank val balance checked after deposit to the market: %s", valBalance)

	acc1Balance = e2esdk.GetSGEBalance(clientCtx, acc1.String())
	expectedBalance := sdkmath.NewInt(int64(acc1initialBalance)).Sub(depositAmount).Sub(depositFeeAmount).Sub(sdkmath.NewInt(fees))
	s.Require().Equal(expectedBalance, acc1Balance)
	s.T().Logf("==== bank acc1 balance checked after deposit to the market: %s", acc1Balance)

	depositsBz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryDeposits(), []string{
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	})
	require.NoError(s.T(), err)
	deposits := types.QueryDepositsResponse{}
	err = clientCtx.Codec.UnmarshalJSON(depositsBz.Bytes(), &deposits)
	require.NoError(s.T(), err)

	require.Equal(s.T(), 3, len(deposits.Deposits))

	depositBz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryDepositsByAccount(), []string{
		acc1.String(),
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	})
	require.NoError(s.T(), err)
	deposit := types.QueryDepositsByAccountResponse{}
	err = clientCtx.Codec.UnmarshalJSON(depositBz.Bytes(), &deposit)
	require.NoError(s.T(), err)

	require.Equal(s.T(), []types.Deposit{
		{
			Creator:               acc1.String(),
			DepositorAddress:      acc1.String(),
			MarketUID:             marketUID,
			ParticipationIndex:    2,
			Amount:                depositAmount,
			WithdrawalCount:       0,
			TotalWithdrawalAmount: sdkmath.ZeroInt(),
		},
	}, deposit.Deposits)
	require.Equal(s.T(), acc1.String(), deposit.Deposits[0].DepositorAddress)
	s.T().Log("==== state modifications checked after deposit to the market")

	s.T().Log("==== new deposit test passed successfully")
}

func (s *E2ETestSuite) TestDepositWithAuthzTxCmd() {
	s.T().Log("==== new house deposit with authorization command test started")
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
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16990",
					Meta: "Home",
				},
				{
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16991",
					Meta: "Draw",
				},
				{
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16992",
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

	expectedValBalance := valBalance.Sub(sdkmath.NewInt(fees))
	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(expectedValBalance, valBalance)

	e2esdk.SetGenericAuthorization(val, acc1, sdk.MsgTypeURL(&types.MsgDeposit{}))
	s.Require().NoError(s.network.WaitForNextBlock())
	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())

	testCases := []struct {
		name           string
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
			marketUID,
			depositAmount,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       val.Address.String(),
				},
				"depositor_address": val.Address.String(),
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "acc1"),
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

			tc.args = append([]string{
				tc.marketUID,
				tc.amount.String(),
				ticket,
			}, tc.args...)
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdDeposit(), tc.args)

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
				s.T().Logf("==== success deposit create test case finished: %s", tc.name)
			}
		})
	}

	depositFeeAmount := sdk.NewDecFromInt(depositAmount).Mul(genesis.Params.HouseParticipationFee).TruncateInt()

	expectedValBalance = valBalance.Sub(depositAmount).Sub(depositFeeAmount)
	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(expectedValBalance, valBalance)
	s.T().Logf("==== bank val balance checked after deposit to the market: %s", valBalance)

	expectedAcc1Balance := acc1Balance.Sub(sdkmath.NewInt(fees))
	acc1Balance = e2esdk.GetSGEBalance(clientCtx, acc1.String())
	s.Require().Equal(expectedAcc1Balance, acc1Balance)
	s.T().Logf("==== bank acc1 balance checked after deposit to the market: %s", acc1Balance)

	depositsBz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryDeposits(), []string{
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	})
	require.NoError(s.T(), err)
	deposits := types.QueryDepositsResponse{}
	err = clientCtx.Codec.UnmarshalJSON(depositsBz.Bytes(), &deposits)
	require.NoError(s.T(), err)

	require.Equal(s.T(), 4, len(deposits.Deposits))

	depositBz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryDepositsByAccount(), []string{
		val.Address.String(),
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	})
	require.NoError(s.T(), err)
	deposit := types.QueryDepositsByAccountResponse{}
	err = clientCtx.Codec.UnmarshalJSON(depositBz.Bytes(), &deposit)
	require.NoError(s.T(), err)

	require.Equal(s.T(), []types.Deposit{
		{
			Creator:               val.Address.String(),
			DepositorAddress:      val.Address.String(),
			MarketUID:             marketUID,
			ParticipationIndex:    1,
			Amount:                depositAmount,
			WithdrawalCount:       0,
			TotalWithdrawalAmount: sdkmath.ZeroInt(),
		},
		{
			Creator:               acc1.String(),
			DepositorAddress:      val.Address.String(),
			MarketUID:             marketUID,
			ParticipationIndex:    3,
			Amount:                depositAmount,
			WithdrawalCount:       0,
			TotalWithdrawalAmount: sdkmath.ZeroInt(),
		},
	}, deposit.Deposits)
	require.Equal(s.T(), val.Address.String(), deposit.Deposits[0].DepositorAddress)
	s.T().Log("==== state modifications checked after deposit to the market")

	s.T().Log("==== new deposit with authorization test passed successfully")
}

func (s *E2ETestSuite) TestWithdrawTxCmd() {
	s.T().Log("==== new house deposit command test started")
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
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16990",
					Meta: "Home",
				},
				{
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16991",
					Meta: "Draw",
				},
				{
					UID:  "9991c60f-2025-48ce-ae79-1dc110f16992",
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

	expectedValBalance := valBalance.Sub(sdkmath.NewInt(fees))
	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(expectedValBalance, valBalance)
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

		bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdDeposit(), []string{
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
		s.T().Logf("==== success deposit create test case finished")
	}

	depositFeeAmount := sdk.NewDecFromInt(depositAmount).Mul(genesis.Params.HouseParticipationFee).TruncateInt()

	expectedValBalance = valBalance.Sub(depositAmount).Sub(depositFeeAmount).Sub(sdkmath.NewInt(fees))
	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(expectedValBalance, valBalance)
	s.T().Logf("==== bank val balance checked after deposit to the market: %s", valBalance)

	depositsBz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryDeposits(), []string{
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	})
	require.NoError(s.T(), err)
	deposits := types.QueryDepositsResponse{}
	err = clientCtx.Codec.UnmarshalJSON(depositsBz.Bytes(), &deposits)
	require.NoError(s.T(), err)

	require.Equal(s.T(), 5, len(deposits.Deposits))

	depositBz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryDepositsByAccount(), []string{
		acc1.String(),
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	})
	require.NoError(s.T(), err)
	deposit := types.QueryDepositsByAccountResponse{}
	err = clientCtx.Codec.UnmarshalJSON(depositBz.Bytes(), &deposit)
	require.NoError(s.T(), err)

	require.Equal(s.T(), []types.Deposit{
		{
			Creator:               acc1.String(),
			DepositorAddress:      acc1.String(),
			MarketUID:             marketUID,
			ParticipationIndex:    2,
			Amount:                depositAmount,
			WithdrawalCount:       0,
			TotalWithdrawalAmount: sdkmath.ZeroInt(),
		},
	}, deposit.Deposits)
	require.Equal(s.T(), acc1.String(), deposit.Deposits[0].DepositorAddress)
	s.T().Log("==== state modifications checked after deposit to the market")

	firstWithdrawAmount := depositAmount.Sub(sdkmath.NewInt(10))
	testCases := []struct {
		name               string
		marketUID          string
		participationIndex uint64
		ticketClaims       jwt.MapClaims
		mode               types.WithdrawalMode
		amount             sdkmath.Int
		args               []string
		expectErr          bool
		expectedCode       uint32
		expectedErrMsg     string
		respType           proto.Message
	}{
		{
			"invalid transaction with large amount for val",
			marketUID,
			4,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       val.Address.String(),
				},
			},
			types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL,
			sdkmath.NewInt(10000000000),
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagGas, flags.GasFlagAuto),
				fmt.Sprintf("--%s=%s", flags.FlagGasAdjustment, "1.15"),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(fees))).String()),
			},
			true,
			0,
			"withdrawal is more than unused amount",
			&sdk.TxResponse{},
		},
		{
			"valid transaction for val",
			marketUID,
			4,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       val.Address.String(),
				},
			},
			types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL,
			firstWithdrawAmount,
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
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
		{
			"second valid transaction for val with the rest of the amount",
			marketUID,
			4,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       val.Address.String(),
				},
			},
			types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
			sdkmath.Int{},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
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

			tc.args = append([]string{
				tc.marketUID,
				cast.ToString(tc.participationIndex),
				ticket,
				cast.ToString(int32(tc.mode)),
				tc.amount.String(),
			}, tc.args...)
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdWithdraw(), tc.args)

			if tc.expectErr {
				s.Require().Error(err)
				s.T().Logf("error captured: %s", tc.name)
			} else {
				s.Require().NoError(err)

				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), tc.respType), bz.String())
				respType := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, respType.Code)

				txResp, err := clitestutil.GetTxResponse(s.network, clientCtx, respType.TxHash)
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedCode, txResp.Code)
				if tc.expectedErrMsg != "" {
					s.Require().Contains(txResp.RawLog, tc.expectedErrMsg)
				}
				s.T().Logf("==== success deposit create test case finished: %s", tc.name)
			}
		})
	}

	s.T().Log("==== new deposit test passed successfully")
}

func (s *E2ETestSuite) TestWithdrawWithAuthzTxCmd() {
	s.T().Log("==== new house deposit command test started")
	val := s.network.Validators[0]

	clientCtx := val.ClientCtx
	// {
	// 	ticket, err := simapp.CreateJwtTicket(jwt.MapClaims{
	// 		"exp":      time.Now().Add(time.Minute * 5).Unix(),
	// 		"iat":      time.Now().Unix(),
	// 		"uid":      marketUID,
	// 		"start_ts": marketStartDate,
	// 		"end_ts":   marketEndDate,
	// 		"odds": []markettypes.Odds{
	// 			{
	// 				UID:  "9991c60f-2025-48ce-ae79-1dc110f16990",
	// 				Meta: "Home",
	// 			},
	// 			{
	// 				UID:  "9991c60f-2025-48ce-ae79-1dc110f16991",
	// 				Meta: "Draw",
	// 			},
	// 			{
	// 				UID:  "9991c60f-2025-48ce-ae79-1dc110f16992",
	// 				Meta: "Away",
	// 			},
	// 		},
	// 		"status": markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
	// 		"meta":   "sample 3way market",
	// 	})
	// 	require.Nil(s.T(), err)

	// 	args := []string{
	// 		ticket,
	// 		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
	// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
	// 		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(fees))).String()),
	// 	}
	// 	bz, err := clitestutil.ExecTestCLICmd(clientCtx, marketcli.CmdAdd(), args)
	// 	s.Require().NoError(err)
	// 	respType := sdk.TxResponse{}
	// 	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), &respType), bz.String())
	// 	s.Require().Equal(uint32(0), respType.Code)
	// }
	// s.Require().NoError(s.network.WaitForNextBlock())

	acc2initialBalance := 1000
	acc2 = e2esdk.CreateAccount(val, "acc2")
	e2esdk.SendToken(val, acc2, acc2initialBalance)
	s.Require().NoError(s.network.WaitForNextBlock())
	acc2Balance = e2esdk.GetSGEBalance(clientCtx, acc2.String())

	e2esdk.SetGenericAuthorization(val, acc2, sdk.MsgTypeURL(&types.MsgWithdraw{}))
	s.Require().NoError(s.network.WaitForNextBlock())

	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())

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

		bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdDeposit(), []string{
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
		s.T().Logf("==== success deposit create test case finished")
	}

	depositFeeAmount := sdk.NewDecFromInt(depositAmount).Mul(genesis.Params.HouseParticipationFee).TruncateInt()

	expectedValBalance := valBalance.Sub(depositAmount).Sub(depositFeeAmount).Sub(sdkmath.NewInt(fees))
	valBalance = e2esdk.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(expectedValBalance, valBalance)
	s.T().Logf("==== bank val balance checked after deposit to the market: %s", valBalance)

	depositsBz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryDeposits(), []string{
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	})
	require.NoError(s.T(), err)
	deposits := types.QueryDepositsResponse{}
	err = clientCtx.Codec.UnmarshalJSON(depositsBz.Bytes(), &deposits)
	require.NoError(s.T(), err)

	require.Equal(s.T(), 6, len(deposits.Deposits))

	depositBz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.GetCmdQueryDepositsByAccount(), []string{
		val.Address.String(),
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	})
	require.NoError(s.T(), err)
	deposit := types.QueryDepositsByAccountResponse{}
	err = clientCtx.Codec.UnmarshalJSON(depositBz.Bytes(), &deposit)
	require.NoError(s.T(), err)

	require.Equal(s.T(), val.Address.String(), deposit.Deposits[0].DepositorAddress)
	s.T().Log("==== state modifications checked after deposit to the market")

	firstWithdrawAmount := depositAmount.Sub(sdkmath.NewInt(10))
	testCases := []struct {
		name               string
		marketUID          string
		participationIndex uint64
		ticketClaims       jwt.MapClaims
		mode               types.WithdrawalMode
		amount             sdkmath.Int
		args               []string
		expectErr          bool
		expectedCode       uint32
		expectedErrMsg     string
		respType           proto.Message
	}{
		{
			"valid transaction for acc2",
			marketUID,
			5,
			jwt.MapClaims{
				"exp": time.Now().Add(time.Minute * 5).Unix(),
				"iat": time.Now().Unix(),
				"kyc_data": sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       val.Address.String(),
				},
				"depositor_address": val.Address.String(),
			},
			types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL,
			firstWithdrawAmount,
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "acc2"),
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

			tc.args = append([]string{
				tc.marketUID,
				cast.ToString(tc.participationIndex),
				ticket,
				cast.ToString(int32(tc.mode)),
				tc.amount.String(),
			}, tc.args...)
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdWithdraw(), tc.args)

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
				s.T().Logf("==== success deposit create test case finished: %s", tc.name)
			}
		})
	}

	expectedAcc2Balance := acc2Balance.Sub(sdkmath.NewInt(fees))
	acc2Balance = e2esdk.GetSGEBalance(clientCtx, acc2.String())
	s.Require().Equal(expectedAcc2Balance, acc2Balance)

	s.T().Log("==== new deposit test passed successfully")
}
