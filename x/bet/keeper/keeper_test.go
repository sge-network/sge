package keeper_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/simapp"

	sgetypes "github.com/sge-network/sge/types"
	"github.com/sge-network/sge/x/bet/keeper"
	"github.com/sge-network/sge/x/bet/types"
	marketkeeper "github.com/sge-network/sge/x/market/keeper"
	markettypes "github.com/sge-network/sge/x/market/types"
)

var (
	testMarketUID  = "5db09053-2901-4110-8fb5-c14e21f8d555"
	testOddsUID1   = "6db09053-2901-4110-8fb5-c14e21f8d666"
	testOddsUID2   = "5e31c60f-2025-48ce-ae79-1dc110f16358"
	testOddsUID3   = "6e31c60f-2025-48ce-ae79-1dc110f16354"
	testMarketOdds = []*markettypes.Odds{
		{UID: testOddsUID1, Meta: "Odds 1"},
		{UID: testOddsUID2, Meta: "Odds 2"},
		{UID: testOddsUID3, Meta: "Odds 3"},
	}
	testSelectedBetOdds = &types.BetOdds{
		UID:               testOddsUID1,
		MarketUID:         testMarketUID,
		Value:             "4.20",
		MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
	}
	testBetOdds = &[]types.BetOdds{
		{
			UID:               testOddsUID1,
			MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
		},
		{
			UID:               testOddsUID2,
			MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
		},
		{
			UID:               testOddsUID3,
			MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
		},
	}
	testCreator   string
	testBet       *types.MsgWager
	testAddMarket *markettypes.MsgAdd

	testMarket = markettypes.Market{
		UID:     testMarketUID,
		Creator: simapp.TestParamUsers["user1"].Address.String(),
		StartTS: 1111111111,
		EndTS:   uint64(time.Now().Unix()) + 5000,
		Odds:    testMarketOdds,
		Status:  markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
		PriceStats: &markettypes.PriceStats{
			ResolutionSgePrice: sdk.ZeroDec(),
		},
	}
)

func setupKeeperAndApp(t testing.TB) (*simapp.TestApp, *keeper.KeeperTest, sdk.Context) {
	tApp, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)

	return tApp, tApp.BetKeeper, ctx
}

func setupKeeper(t testing.TB) (*keeper.KeeperTest, sdk.Context) {
	_, k, ctx := setupKeeperAndApp(t)

	return k, ctx
}

func addTestMarket(t testing.TB, tApp *simapp.TestApp, ctx sdk.Context) {
	testCreator = simapp.TestParamUsers["user1"].Address.String()
	testAddMarketClaim := jwt.MapClaims{
		"uid":      testMarketUID,
		"start_ts": 1111111111,
		"end_ts":   uint64(ctx.BlockTime().Unix()) + 1000,
		"odds":     testMarketOdds,
		"exp":      9999999999,
		"iat":      7777777777,
		"meta":     "Winner of x:y",
		"status":   markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
	}
	testAddMarketTicket, err := createJwtTicket(testAddMarketClaim)
	require.Nil(t, err)

	testAddMarket = &markettypes.MsgAdd{
		Creator: testCreator,
		Ticket:  testAddMarketTicket,
	}
	wctx := sdk.WrapSDKContext(ctx)
	marketSrv := marketkeeper.NewMsgServerImpl(*tApp.MarketKeeper)
	resAddMarket, err := marketSrv.Add(wctx, testAddMarket)
	require.Nil(t, err)
	require.NotNil(t, resAddMarket)
}

func addTestMarketBatch(
	t testing.TB,
	tApp *simapp.TestApp,
	ctx sdk.Context,
	count int,
) (uids []string) {
	for i := 0; i < count; i++ {
		testCreator = simapp.TestParamUsers["user"+cast.ToString(i)].Address.String()
		uid := uuid.NewString()
		uids = append(uids, uid)
		testAddMarketClaim := jwt.MapClaims{
			"uid":      uid,
			"start_ts": 1111111111,
			"end_ts":   uint64(ctx.BlockTime().Unix()) + 1000,
			"odds":     testMarketOdds,
			"exp":      9999999999,
			"iat":      7777777777,
			"meta":     "Winner of x:y",
			"status":   markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
		}
		testAddMarketTicket, err := createJwtTicket(testAddMarketClaim)
		require.Nil(t, err)

		testAddMarket = &markettypes.MsgAdd{
			Creator: testCreator,
			Ticket:  testAddMarketTicket,
		}
		wctx := sdk.WrapSDKContext(ctx)
		marketSrv := marketkeeper.NewMsgServerImpl(*tApp.MarketKeeper)
		resAddMarket, err := marketSrv.Add(wctx, testAddMarket)
		require.Nil(t, err)
		require.NotNil(t, resAddMarket)
	}

	return uids
}

func placeTestBet(
	ctx sdk.Context,
	t testing.TB,
	tApp *simapp.TestApp,
	betUID string,
	selectedOdds *types.BetOdds,
) {
	testCreator = simapp.TestParamUsers["user1"].Address.String()
	wctx := sdk.WrapSDKContext(ctx)
	betSrv := keeper.NewMsgServerImpl(*tApp.BetKeeper)
	testKyc := &sgetypes.KycDataPayload{
		Approved: true,
		ID:       testCreator,
	}

	if selectedOdds == nil {
		selectedOdds = testSelectedBetOdds
	}

	testWagerClaim := jwt.MapClaims{
		"exp":           9999999999,
		"iat":           7777777777,
		"selected_odds": selectedOdds,
		"kyc_data":      testKyc,
		"all_odds":      testBetOdds,
	}
	testWagerTicket, err := createJwtTicket(testWagerClaim)
	require.Nil(t, err)

	testBet = &types.MsgWager{
		Creator: testCreator,
		Props: &types.WagerProps{
			UID:    betUID,
			Amount: sdkmath.NewInt(1000000),
			Ticket: testWagerTicket,
		},
	}
	resWagerBet, err := betSrv.Wager(wctx, testBet)
	require.Nil(t, err)
	require.NotNil(t, resWagerBet)
}

func createJwtTicket(claim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claim)
	return token.SignedString(simapp.TestOVMPrivateKeys[0])
}

func TestLogger(t *testing.T) {
	k, ctx := setupKeeper(t)
	logger := k.Logger(ctx)
	require.NotNil(t, logger)

	logger.Debug("test log")
}
