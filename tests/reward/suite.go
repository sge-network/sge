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
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/sge-network/sge/testutil/network"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/reward/client/cli"
	"github.com/sge-network/sge/x/reward/types"
)

var (
	genesis           types.GenesisState
	promoterUID       = uuid.NewString()
	campaignStartDate = uint64(time.Now().Unix())
	campaignEndDate   = uint64(time.Now().Add(5 * time.Minute).Unix())
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

	// promoterAddr = simapp.TestParamUsers["user1"].Address.String()

	// genesis.PromoterList = []types.Promoter{
	// 	{
	// 		Creator:   promoterAddr,
	// 		UID:       uuid.NewString(),
	// 		Addresses: []string{promoterAddr},
	// 		Conf: types.PromoterConf{
	// 			CategoryCap: []types.CategoryCap{{
	// 				Category:  types.RewardCategory_REWARD_CATEGORY_SIGNUP,
	// 				CapPerAcc: 1,
	// 			}},
	// 		},
	// 	},
	// }

	// genesis.PromoterByAddressList = []types.PromoterByAddress{
	// 	{
	// 		PromoterUID: genesis.PromoterList[0].UID,
	// 		Address:     promoterAddr,
	// 	},
	// }

	// genesis.CampaignList = []types.Campaign{
	// 	{
	// 		Creator:          promoterAddr,
	// 		UID:              uuid.NewString(),
	// 		Promoter:         promoterAddr,
	// 		StartTS:          campaignStartDate,
	// 		EndTS:            campaignEndDate,
	// 		RewardCategory:   types.RewardCategory_REWARD_CATEGORY_SIGNUP,
	// 		RewardType:       types.RewardType_REWARD_TYPE_SIGNUP,
	// 		RewardAmountType: types.RewardAmountType_REWARD_AMOUNT_TYPE_FIXED,
	// 		RewardAmount: &types.RewardAmount{
	// 			SubaccountAmount: sdkmath.NewInt(100),
	// 		},
	// 		Pool: ,
	// 	},
	// }

	genesisBz, err := s.cfg.Codec.MarshalJSON(&genesis)
	s.Require().NoError(err)
	genesisState[types.ModuleName] = genesisBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())
}

// func (s *E2ETestSuite) TearDownSuite() {
// 	s.T().Log("tearing down e2e test suite")
// 	s.network.Cleanup()
// }

func (s *E2ETestSuite) TestNewCampaignTxCmd() {
	val := s.network.Validators[0]

	{
		clientCtx := val.ClientCtx

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

	testCases := []struct {
		name         string
		uid          string
		totalFunds   sdkmath.Int
		ticketClaims jwt.MapClaims
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			"valid transaction",
			uuid.NewString(),
			sdkmath.NewInt(100),
			jwt.MapClaims{
				"exp":                time.Now().Add(time.Minute * 5).Unix(),
				"iat":                time.Now().Unix(),
				"promoter":           val.Address,
				"start_ts":           cast.ToString(campaignStartDate),
				"end_ts":             cast.ToString(campaignEndDate),
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
			false, 0, &sdk.TxResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc

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
			} else {
				s.Require().NoError(err)

				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), tc.respType), bz.String())
				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code)
			}
		})
	}
}
