package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	sgetypes "github.com/sge-network/sge/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	marketkeeper "github.com/sge-network/sge/x/market/keeper"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/subaccount/types"
	"github.com/stretchr/testify/require"
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
	testSelectedBetOdds = &bettypes.BetOdds{
		UID:               testOddsUID1,
		MarketUID:         testMarketUID,
		Value:             "4.20",
		MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
	}
	testCreator string
)

var (
	subAccOwner  = sample.NativeAccAddress()
	subAccFunder = sample.NativeAccAddress()
	micro        = sdk.NewInt(1_000_000)
	subAccFunds  = sdk.NewInt(10_000).Mul(micro)
)

func TestMsgServer_Bet(t *testing.T) {
	app, k, msgServer, ctx := setupMsgServerAndApp(t)

	// do subaccount creation
	require.NoError(
		t,
		simapp.FundAccount(
			app.BankKeeper,
			ctx,
			subAccFunder,
			sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).LockedBalanceDenom, subAccFunds)),
		),
	)

	_, err := msgServer.CreateSubAccount(sdk.WrapSDKContext(ctx), &types.MsgCreateSubAccount{
		Sender:          subAccFunder.String(),
		SubAccountOwner: subAccOwner.String(),
		LockedBalances: []*types.LockedBalance{
			{
				UnlockTime: time.Now().Add(24 * time.Hour),
				Amount:     subAccFunds,
			},
		},
	})
	require.NoError(t, err)

	// add market
	addTestMarket(t, app, ctx)

	// start betting using the subaccount
	betAmt := sdk.NewInt(1000).Mul(micro)
	_, err = msgServer.PlaceBet(
		sdk.WrapSDKContext(ctx),
		&types.MsgPlaceBet{Msg: testBet(t, subAccOwner, betAmt)},
	)
	require.NoError(t, err)

	t.Run("resolve market – better wins", func(t *testing.T) {
		ctx, _ = ctx.CacheContext()
		// resolve the market – better wins

	})
	// resolve the market – better loses
	t.Run("resolve market – better loses", func(t *testing.T) {
		ctx, _ = ctx.CacheContext()
		// resolve the market – better loses

	})
	t.Run("resolve market – refund", func(t *testing.T) {
		ctx, _ = ctx.CacheContext()
		// resolve the market – refund

	})
}

func addTestMarket(t testing.TB, tApp *simappUtil.TestApp, ctx sdk.Context) {
	testCreator = simappUtil.TestParamUsers["user1"].Address.String()
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

	testAddMarket := &markettypes.MsgAddMarket{
		Creator: testCreator,
		Ticket:  testAddMarketTicket,
	}
	wctx := sdk.WrapSDKContext(ctx)
	marketSrv := marketkeeper.NewMsgServerImpl(*tApp.MarketKeeper)
	resAddMarket, err := marketSrv.AddMarket(wctx, testAddMarket)
	require.Nil(t, err)
	require.NotNil(t, resAddMarket)

	// add liquidity
	err = simapp.FundAccount(
		tApp.BankKeeper,
		ctx,
		simappUtil.TestParamUsers["user1"].Address,
		sdk.NewCoins(sdk.NewCoin(tApp.SubaccountKeeper.GetParams(ctx).LockedBalanceDenom, sdk.NewInt(1_000_000).Mul(micro))),
	)
	require.NoError(t, err)
	_, err = tApp.OrderbookKeeper.InitiateOrderBookParticipation(
		ctx,
		simappUtil.TestParamUsers["user1"].Address,
		resAddMarket.Data.UID,
		sdk.NewInt(1_000_000).Mul(micro),
		sdk.NewInt(1),
	)
	require.NoError(t, err)
}

func createJwtTicket(claim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claim)
	return token.SignedString(simappUtil.TestOVMPrivateKeys[0])
}

func testBet(t testing.TB, better sdk.AccAddress, amount sdk.Int) *bettypes.MsgPlaceBet {
	ticket, err := createJwtTicket(jwt.MapClaims{
		"exp":           9999999999,
		"iat":           7777777777,
		"selected_odds": testSelectedBetOdds,
		"kyc_data": &sgetypes.KycDataPayload{
			Approved: true,
			ID:       better.String(),
		},
		"odds_type": 1,
	})
	require.NoError(t, err)

	return &bettypes.MsgPlaceBet{
		Creator: better.String(),
		Bet: &bettypes.PlaceBetFields{
			UID:    uuid.NewString(),
			Amount: amount,
			Ticket: ticket,
		},
	}
}
