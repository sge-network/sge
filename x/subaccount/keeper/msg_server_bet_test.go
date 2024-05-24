package keeper_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/testutil/simapp"
	sgetypes "github.com/sge-network/sge/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	marketkeeper "github.com/sge-network/sge/x/market/keeper"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/subaccount/types"
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
		MaxLossMultiplier: sdkmath.LegacyMustNewDecFromStr("0.1"),
	}
	testBetOdds = &[]bettypes.BetOdds{
		{
			UID:               testOddsUID1,
			MaxLossMultiplier: sdkmath.LegacyMustNewDecFromStr("0.1"),
		},
		{
			UID:               testOddsUID2,
			MaxLossMultiplier: sdkmath.LegacyMustNewDecFromStr("0.1"),
		},
		{
			UID:               testOddsUID3,
			MaxLossMultiplier: sdkmath.LegacyMustNewDecFromStr("0.1"),
		},
	}
	testCreator string
)

var (
	micro       = sdkmath.NewInt(1_000_000)
	subAccFunds = sdkmath.NewInt(10_000).Mul(micro)
	subAccAddr  = types.NewAddressFromSubaccount(1)
)

func TestMsgServer_Bet(t *testing.T) {
	app, k, msgServer, ctx := setupMsgServerAndApp(t)

	parm := k.GetParams(ctx)
	parm.WagerEnabled = true
	k.SetParams(ctx, parm)

	subAccOwner := simapp.TestParamUsers["user2"].Address
	subAccFunder := simapp.TestParamUsers["user1"].Address

	initialOwnerBalance := app.BankKeeper.GetBalance(ctx, subAccOwner, params.DefaultBondDenom).Amount

	// do subaccount creation
	require.NoError(
		t,
		testutil.FundAccount(
			app.BankKeeper,
			ctx,
			subAccFunder,
			sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, subAccFunds)),
		),
	)

	_, err := msgServer.Create(sdk.WrapSDKContext(ctx), &types.MsgCreate{
		Creator: subAccFunder.String(),
		Owner:   subAccOwner.String(),
		LockedBalances: []types.LockedBalance{
			{
				UnlockTS: uint64(time.Now().Add(24 * time.Hour).Unix()),
				Amount:   subAccFunds,
			},
		},
	})
	require.NoError(t, err)

	// add market
	market := addTestMarket(t, app, ctx, true)

	// start betting using the subaccount
	betAmt := sdkmath.NewInt(1000).Mul(micro)
	halfBetAmt := betAmt.Quo(sdkmath.NewInt(2))

	_, err = msgServer.Wager(
		sdk.WrapSDKContext(ctx),
		testMsgWager(t, subAccOwner, betAmt, halfBetAmt, halfBetAmt),
	)
	require.NoError(t, err)

	afterBetOwnerBalance := app.BankKeeper.GetBalance(ctx, subAccOwner, params.DefaultBondDenom).Amount

	require.Equal(t, initialOwnerBalance.Sub(halfBetAmt).Int64(), afterBetOwnerBalance.Int64())

	// check subaccount accSumm
	accSumm, exists := k.GetAccountSummary(ctx, subAccAddr)
	require.True(t, exists)
	betFees := sdkmath.NewInt(100)

	require.Equal(t, accSumm.SpentAmount, sdkmath.ZeroInt())
	require.Equal(t, accSumm.WithdrawnAmount, halfBetAmt)

	t.Run("resolve market – bettor wins", func(t *testing.T) {
		ctx, _ := ctx.CacheContext()
		// resolve the market – bettor wins
		app.MarketKeeper.Resolve(ctx, *market, &markettypes.MarketResolutionTicketPayload{
			UID:            market.UID,
			ResolutionTS:   uint64(ctx.BlockTime().Unix()) + 10000,
			WinnerOddsUIDs: []string{testOddsUID1},
			Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
		})
		err := app.BetKeeper.BatchMarketSettlements(ctx)
		require.NoError(t, err)

		// now we check the subaccount accSumm
		accSumm, exists := k.GetAccountSummary(ctx, subAccAddr)
		require.True(t, exists)
		require.Equal(t, halfBetAmt.String(), accSumm.WithdrawnAmount.String())

		// now we want the user to have some balance which is the payout
		winningAmount := betAmt.Sub(betFees).ToLegacyDec().
			Mul(sdkmath.LegacyMustNewDecFromStr("4.2")).TruncateInt().
			Sub(halfBetAmt)

		ownerBalance := app.BankKeeper.GetBalance(ctx, subAccOwner, params.DefaultBondDenom).Amount
		require.Equal(t,
			initialOwnerBalance.Add(winningAmount).Int64(), // 4.2 - 1 = 3.2
			ownerBalance.Int64(),
		)
	})

	// resolve the market – bettor loses
	t.Run("resolve market – bettor loses", func(t *testing.T) {
		ctx, _ := ctx.CacheContext()
		// resolve the market – bettor loses
		app.MarketKeeper.Resolve(ctx, *market, &markettypes.MarketResolutionTicketPayload{
			UID:            market.UID,
			ResolutionTS:   uint64(ctx.BlockTime().Unix()) + 10000,
			WinnerOddsUIDs: []string{testOddsUID2},
			Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
		})
		err := app.BetKeeper.BatchMarketSettlements(ctx)
		require.NoError(t, err)

		// now we check the subaccount balance
		balance, exists := k.GetAccountSummary(ctx, subAccAddr)
		require.True(t, exists)
		require.Equal(t, sdkmath.ZeroInt(), balance.SpentAmount)
		require.Equal(t, sdkmath.ZeroInt(), balance.LostAmount)
		require.Equal(t, halfBetAmt, balance.WithdrawnAmount)
		// the owner has no balances
		ownerBalance := app.BankKeeper.GetBalance(ctx, subAccOwner, params.DefaultBondDenom).Amount
		require.Equal(t,
			initialOwnerBalance.Sub(halfBetAmt).Int64(),
			ownerBalance.Int64(),
		)
	})
	t.Run("resolve market – refund", func(t *testing.T) {
		ctx, _ := ctx.CacheContext()
		// resolve the market – refund
		app.MarketKeeper.Resolve(ctx, *market, &markettypes.MarketResolutionTicketPayload{
			UID:            market.UID,
			ResolutionTS:   uint64(ctx.BlockTime().Unix()) + 10000,
			WinnerOddsUIDs: []string{testOddsUID1},
			Status:         markettypes.MarketStatus_MARKET_STATUS_CANCELED,
		})
		err := app.BetKeeper.BatchMarketSettlements(ctx)
		require.NoError(t, err)

		// now we check the subaccount balance
		balance, exists := k.GetAccountSummary(ctx, subAccAddr)
		require.True(t, exists)
		require.Equal(t, balance.SpentAmount, sdkmath.ZeroInt())

		ownerBalance := app.BankKeeper.GetBalance(ctx, subAccOwner, params.DefaultBondDenom).Amount

		// the owner balance is zero
		require.Equal(t, initialOwnerBalance.Add(halfBetAmt).Int64(), ownerBalance.Int64())
	})
}

func addTestMarket(t testing.TB, tApp *simapp.TestApp, ctx sdk.Context, prefund bool) *markettypes.Market {
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

	testAddMarket := &markettypes.MsgAdd{
		Creator: testCreator,
		Ticket:  testAddMarketTicket,
	}
	wctx := sdk.WrapSDKContext(ctx)
	marketSrv := marketkeeper.NewMsgServerImpl(*tApp.MarketKeeper)
	resAddMarket, err := marketSrv.Add(wctx, testAddMarket)
	require.Nil(t, err)
	require.NotNil(t, resAddMarket)

	if prefund {
		// add liquidity
		err = testutil.FundAccount(
			tApp.BankKeeper,
			ctx,
			simapp.TestParamUsers["user1"].Address,
			sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdkmath.NewInt(1_000_000).Mul(micro))),
		)
		require.NoError(t, err)
		_, err = tApp.OrderbookKeeper.InitiateOrderBookParticipation(
			ctx,
			simapp.TestParamUsers["user1"].Address,
			resAddMarket.Data.UID,
			sdkmath.NewInt(1_000_000).Mul(micro),
			sdkmath.NewInt(1),
		)
		require.NoError(t, err)
	}
	return resAddMarket.Data
}

func createJwtTicket(claim jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claim)
	return token.SignedString(simapp.TestOVMPrivateKeys[0])
}

func testMsgWager(t testing.TB, bettor sdk.AccAddress, amount sdkmath.Int, mainAccDeduct, subAccDeduct sdkmath.Int) *types.MsgWager {
	betTicket, err := createJwtTicket(jwt.MapClaims{
		"exp":           9999999999,
		"iat":           7777777777,
		"selected_odds": testSelectedBetOdds,
		"kyc_data": &sgetypes.KycDataPayload{
			Approved: true,
			ID:       bettor.String(),
		},
		"all_odds": testBetOdds,
	})
	require.NoError(t, err)

	subAccTicket, err := createJwtTicket(jwt.MapClaims{
		"exp": 9999999999,
		"iat": 7777777777,
		"msg": bettypes.MsgWager{
			Creator: bettor.String(),
			Props: &bettypes.WagerProps{
				UID:    uuid.NewString(),
				Amount: amount,
				Ticket: betTicket,
			},
		},
		"mainacc_deduct_amount": mainAccDeduct,
		"subacc_deduct_amount":  subAccDeduct,
	})
	require.NoError(t, err)

	return &types.MsgWager{
		Creator: bettor.String(),
		Ticket:  subAccTicket,
	}
}
