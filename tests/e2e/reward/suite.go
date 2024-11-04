package client

import (
	"fmt"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	sdkcli "github.com/sge-network/sge/tests/e2e/sdk"
	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/simapp"
	sgetypes "github.com/sge-network/sge/types"
	"github.com/sge-network/sge/x/reward/client/cli"
	"github.com/sge-network/sge/x/reward/types"
	subaccountcli "github.com/sge-network/sge/x/subaccount/client/cli"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

var (
	genesis           types.GenesisState
	fees              = int64(10)
	promoterUID       = uuid.NewString()
	campaignStartDate = uint64(time.Now().Unix())
	campaignEndDate   = time.Now().Add(time.Minute * time.Duration(120)).Unix()
	campaignFunds     = sdkmath.NewInt(100000000)
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

	genesisBz, err := s.cfg.Codec.MarshalJSON(&genesis)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = genesisBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *E2ETestSuite) TestNewCampaignTxCmd() {
	s.T().Log("==== new campaign create command test started")
	val := s.network.Validators[0]

	clientCtx := val.ClientCtx
	{
		ticket, err := simapp.CreateJwtTicket(jwt.MapClaims{
			"exp": time.Now().Add(time.Minute * 5).Unix(),
			"iat": time.Now().Unix(),
			"uid": promoterUID,
			"conf": types.PromoterConf{
				CategoryCap: []types.CategoryCap{{
					Category:  types.RewardCategory_REWARD_CATEGORY_SIGNUP,
					CapPerAcc: 1,
				}},
			},
		})
		require.Nil(s.T(), err)

		args := []string{
			ticket,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(10))).String()),
		}
		bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdCreatePromoter(), args)
		s.Require().NoError(err)
		respType := sdk.TxResponse{}
		s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), &respType), bz.String())
		s.Require().Equal(uint32(0), respType.Code)
	}

	valBalance := sdkcli.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(sdkmath.NewInt(400000000), valBalance)

	campaignUID := uuid.NewString()
	testCases := []struct {
		name           string
		uid            string
		totalFunds     sdkmath.Int
		ticketClaims   jwt.MapClaims
		args           []string
		expectErr      bool
		expectedCode   uint32
		expectedErrMsg string
		respType       proto.Message
	}{
		{
			"not enough balance to charge pool",
			campaignUID,
			sdkmath.NewInt(4000000000),
			jwt.MapClaims{
				"exp":                time.Now().Add(time.Minute * 5).Unix(),
				"iat":                time.Now().Unix(),
				"promoter":           val.Address,
				"start_ts":           campaignStartDate,
				"end_ts":             campaignEndDate,
				"category":           types.RewardCategory_REWARD_CATEGORY_SIGNUP,
				"reward_type":        types.RewardType_REWARD_TYPE_SIGNUP,
				"reward_amount_type": types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED,
				"reward_amount": types.RewardAmount{
					SubaccountAmount: sdkmath.NewInt(100),
					UnlockPeriod:     uint64(1000),
				},
				"is_active": true,
				"meta":      "sample signup campaign",
				"cap_count": 1,
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(10))).String()),
			},
			false,
			types.ErrInFundingCampaignPool.ABCICode(),
			"",
			&sdk.TxResponse{},
		},
		{
			"valid transaction",
			campaignUID,
			campaignFunds,
			jwt.MapClaims{
				"exp":                time.Now().Add(time.Minute * 5).Unix(),
				"iat":                time.Now().Unix(),
				"promoter":           val.Address,
				"start_ts":           campaignStartDate,
				"end_ts":             campaignEndDate,
				"category":           types.RewardCategory_REWARD_CATEGORY_SIGNUP,
				"reward_type":        types.RewardType_REWARD_TYPE_SIGNUP,
				"reward_amount_type": types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED,
				"reward_amount": types.RewardAmount{
					SubaccountAmount: sdkmath.NewInt(100),
					UnlockPeriod:     uint64(1000),
				},
				"is_active": true,
				"meta":      "sample signup campaign",
				"cap_count": 1,
			},
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(10))).String()),
			},
			false,
			0,
			"",
			&sdk.TxResponse{},
		},
	}

	for _, tc := range testCases {
		s.Require().NoError(s.network.WaitForNextBlock())
		s.Run(tc.name, func() {
			clientCtx := val.ClientCtx

			ticket, err := simapp.CreateJwtTicket(tc.ticketClaims)
			require.Nil(s.T(), err)

			tc.args = append([]string{
				tc.uid,
				tc.totalFunds.String(),
				ticket,
			}, tc.args...)
			bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdCreateCampaign(), tc.args)

			if tc.expectErr {
				s.Require().Error(err)
				s.T().Logf("error captured: %s", tc.name)
			} else {
				s.Require().NoError(err)

				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), tc.respType), bz.String())

				txResp, err := clitestutil.GetTxResponse(s.network, clientCtx, tc.respType.(*sdk.TxResponse).TxHash)
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedCode, txResp.Code)
				if tc.expectedErrMsg != "" {
					s.Require().Contains(txResp.RawLog, tc.expectedErrMsg)
				}
				s.T().Logf("==== success campaign create test case finished: %s", tc.name)
			}
		})
	}

	s.Require().NoError(s.network.WaitForNextBlock())
	expectedBalance := valBalance.Sub(campaignFunds).Sub(sdk.NewInt(fees * 3))
	valBalance = sdkcli.GetSGEBalance(clientCtx, val.Address.String())
	s.Require().Equal(expectedBalance, valBalance)
	s.T().Logf("==== bank balance checked after creating the campaign: %s", valBalance)

	{
		ticket, err := simapp.CreateJwtTicket(jwt.MapClaims{
			"exp": time.Now().Add(time.Minute * 5).Unix(),
			"iat": time.Now().Unix(),
			"common": types.RewardPayloadCommon{
				Receiver:  val.Address.String(),
				SourceUID: "",
				Meta:      "signup reward campaign",
				KycData: &sgetypes.KycDataPayload{
					Ignore:   false,
					Approved: true,
					ID:       val.Address.String(),
				},
			},
		})
		require.Nil(s.T(), err)

		rewardUID := uuid.NewString()
		args := []string{
			rewardUID,
			campaignUID,
			ticket,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdkmath.NewInt(10))).String()),
		}
		bz, err := clitestutil.ExecTestCLICmd(clientCtx, cli.CmdGrantReward(), args)
		s.Require().NoError(err)
		respType := sdk.TxResponse{}

		s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), &respType), bz.String())

		txResp, err := clitestutil.GetTxResponse(s.network, clientCtx, respType.TxHash)
		s.Require().NoError(err)
		s.Require().Equal(uint32(0), txResp.Code)
		s.T().Logf("==== success reward grant test case finished: %s", rewardUID)
	}

	valBalance = sdkcli.GetSGEBalance(clientCtx, val.Address.String())
	expectedBalance = expectedBalance.Sub(sdkmath.NewInt(10))
	s.Require().Equal(expectedBalance, valBalance)
	s.T().Logf("==== bank balance checked after creating the campaign: %s", valBalance)

	args := []string{
		val.Address.String(),
		fmt.Sprintf("--%s=json", flags.FlagOutput),
	}
	bz, err := clitestutil.ExecTestCLICmd(clientCtx, subaccountcli.QuerySubaccount(), args)
	if err != nil {
		panic(err)
	}
	respType := subaccounttypes.QuerySubaccountResponse{}
	err = clientCtx.Codec.UnmarshalJSON(bz.Bytes(), &respType)
	if err != nil {
		panic(err)
	}

	subaccountBalance := sdkcli.GetSGEBalance(clientCtx, respType.Address)
	s.Require().Equal(sdkmath.NewInt(100), subaccountBalance)
	s.T().Logf("==== bank balance checked after creating the campaign: %s", subaccountBalance)

	s.T().Log("==== new campaign create command test passed successfully")
}
