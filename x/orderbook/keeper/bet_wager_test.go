package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/app/params"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	bettypes "github.com/sge-network/sge/x/bet/types"
	housetypes "github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestWager(t *testing.T) {
	ts := newTestBetSuite(t)
	ts.placeBetsAndTest()
}

type testBetSuite struct {
	t              *testing.T
	k              *keeper.KeeperTest
	ctx            sdk.Context
	tApp           simappUtil.TestApp
	betFee         sdkmath.Int
	market         markettypes.Market
	deposits       []housetypes.Deposit
	participations []types.OrderBookParticipation
}

func newTestBetSuite(t *testing.T) testBetSuite {
	tApp, k, ctx := setupKeeperAndApp(t)

	betFee := sdkmath.NewInt(10)

	marketUID := uuid.NewString()
	market := markettypes.Market{
		UID:     marketUID,
		StartTS: cast.ToUint64(time.Now().Unix()),
		EndTS:   cast.ToUint64(time.Now().Add(5 * time.Minute).Unix()),
		Odds: []*markettypes.Odds{
			{UID: uuid.NewString(), Meta: "test odds1"},
			{UID: uuid.NewString(), Meta: "test odds2"},
			{UID: uuid.NewString(), Meta: "test odds3"},
		},
		Status:  markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
		Creator: simappUtil.TestParamUsers["user1"].Address.String(),
		Meta:    "test market",
		BookUID: marketUID,
	}

	deposits := []housetypes.Deposit{
		{
			DepositorAddress: simappUtil.TestParamUsers["user2"].Address.String(),
			Amount:           sdkmath.NewInt(8000),
		},
		{
			DepositorAddress: simappUtil.TestParamUsers["user3"].Address.String(),
			Amount:           sdkmath.NewInt(10000),
		},
		{
			DepositorAddress: simappUtil.TestParamUsers["user4"].Address.String(),
			Amount:           sdkmath.NewInt(10000),
		},
	}

	participations := make([]types.OrderBookParticipation, len(deposits))

	return testBetSuite{t, k, ctx, *tApp, betFee, market, deposits, participations}
}

func (ts *testBetSuite) placeBetsAndTest() ([]bettypes.Bet, sdk.Dec, sdk.Dec) {
	ts.tApp.MarketKeeper.SetMarket(ts.ctx, ts.market)

	err := ts.k.InitiateOrderBook(ts.ctx, ts.market.UID, []string{
		ts.market.Odds[0].UID,
		ts.market.Odds[1].UID,
		ts.market.Odds[2].UID,
	})
	require.NoError(ts.t, err)

	found := false
	participationIndex, err := ts.tApp.HouseKeeper.Deposit(
		ts.ctx,
		ts.deposits[0].DepositorAddress,
		ts.deposits[0].DepositorAddress,
		ts.market.BookUID,
		ts.deposits[0].Amount,
	)
	require.NoError(ts.t, err)
	ts.deposits[0], found = ts.tApp.HouseKeeper.GetDeposit(
		ts.ctx,
		ts.deposits[0].DepositorAddress,
		ts.market.UID,
		participationIndex,
	)
	require.True(ts.t, found)
	ts.participations[0], found = ts.k.GetOrderBookParticipation(
		ts.ctx,
		ts.market.UID,
		participationIndex,
	)
	require.True(ts.t, found)

	participationIndex, err = ts.tApp.HouseKeeper.Deposit(
		ts.ctx,
		ts.deposits[1].DepositorAddress,
		ts.deposits[1].DepositorAddress,
		ts.market.BookUID,
		ts.deposits[1].Amount,
	)
	require.NoError(ts.t, err)
	ts.deposits[1], found = ts.tApp.HouseKeeper.GetDeposit(
		ts.ctx,
		ts.deposits[1].DepositorAddress,
		ts.market.UID,
		participationIndex,
	)
	require.True(ts.t, found)
	ts.participations[1], found = ts.k.GetOrderBookParticipation(
		ts.ctx,
		ts.market.UID,
		participationIndex,
	)
	require.True(ts.t, found)

	participationIndex, err = ts.tApp.HouseKeeper.Deposit(
		ts.ctx,
		ts.deposits[2].DepositorAddress,
		ts.deposits[2].DepositorAddress,
		ts.market.BookUID,
		ts.deposits[2].Amount,
	)
	require.NoError(ts.t, err)
	ts.deposits[2], found = ts.tApp.HouseKeeper.GetDeposit(
		ts.ctx,
		ts.deposits[2].DepositorAddress,
		ts.market.UID,
		participationIndex,
	)
	require.True(ts.t, found)
	ts.participations[2], found = ts.k.GetOrderBookParticipation(
		ts.ctx,
		ts.market.UID,
		participationIndex,
	)
	require.True(ts.t, found)

	oddsExposures, found := ts.k.GetOrderBookOddsExposure(
		ts.ctx,
		ts.market.BookUID,
		ts.market.Odds[0].UID,
	)
	require.True(ts.t, found)
	require.Equal(ts.t, []uint64{1, 2, 3}, oddsExposures.FulfillmentQueue)

	defaultBetAmount := sdkmath.NewInt(400)

	betOdds := make(map[string]*bettypes.BetOddsCompact)
	var oddUIDS []string
	for _, odd := range ts.market.Odds {
		betOdds[odd.UID] = &bettypes.BetOddsCompact{UID: odd.UID, MaxLossMultiplier: sdk.MustNewDecFromStr("0.1")}
		oddUIDS = append(oddUIDS, odd.UID)
	}

	///// winner1 bet placement
	//
	//
	winner1BettorAddr := simappUtil.TestParamUsers["user5"].Address
	winner1Bal := ts.tApp.BankKeeper.GetBalance(ts.ctx, winner1BettorAddr, params.DefaultBondDenom)
	winner1BetID := uint64(1)
	winner1Bet, winner1PayoutProfit, winner1BetFulfillment := ts.placeTestBet(
		winner1BettorAddr,
		ts.market.UID,
		ts.market.Odds[0].UID,
		winner1BetID,
		defaultBetAmount,
		ts.betFee,
		nil,
		betOdds,
		oddUIDS,
	)
	winner1Bet.BetFulfillment = winner1BetFulfillment
	ts.tApp.BetKeeper.SetBet(ts.ctx, winner1Bet, winner1BetID)
	winner1BalAfterWager := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		winner1BettorAddr,
		params.DefaultBondDenom,
	)
	expWinner1BalanceAfterWager := winner1Bal.Amount.Sub(winner1Bet.Amount).Sub(winner1Bet.Fee)
	require.Equal(
		ts.t,
		expWinner1BalanceAfterWager.Int64(),
		winner1BalAfterWager.Amount.Int64(),
	)

	oddsExposures, found = ts.k.GetOrderBookOddsExposure(
		ts.ctx,
		ts.market.BookUID,
		ts.market.Odds[0].UID,
	)
	require.True(ts.t, found)
	require.Equal(ts.t, []uint64{2, 3}, oddsExposures.FulfillmentQueue)

	///// winner2 bet placement
	//
	//
	winner2BettorAddr := simappUtil.TestParamUsers["user6"].Address
	winner2Bal := ts.tApp.BankKeeper.GetBalance(ts.ctx, winner2BettorAddr, params.DefaultBondDenom)
	winner2BetID := uint64(2)
	winner2Bet, winner2PayoutProfit, winner2BetFulfillment := ts.placeTestBet(
		winner2BettorAddr,
		ts.market.UID,
		ts.market.Odds[0].UID,
		winner2BetID,
		defaultBetAmount,
		ts.betFee,
		nil,
		betOdds,
		oddUIDS,
	)

	winner2Bet.BetFulfillment = winner2BetFulfillment
	ts.tApp.BetKeeper.SetBet(ts.ctx, winner2Bet, winner2BetID)
	winner2BalAfterWager := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		winner1BettorAddr,
		params.DefaultBondDenom,
	)
	expWinner2BalanceAfterWager := winner2Bal.Amount.Sub(winner2Bet.Amount).Sub(winner2Bet.Fee)
	require.Equal(
		ts.t,
		expWinner2BalanceAfterWager.Int64(),
		winner2BalAfterWager.Amount.Int64(),
	)
	oddsExposures, found = ts.k.GetOrderBookOddsExposure(
		ts.ctx,
		ts.market.BookUID,
		ts.market.Odds[0].UID,
	)
	require.True(ts.t, found)
	require.Equal(ts.t, []uint64{3}, oddsExposures.FulfillmentQueue)

	///// failed winner bet placement
	// should fail because there is not participation to fulfill this bet.
	//
	failedWinnerBettorAddr := simappUtil.TestParamUsers["user7"].Address
	failedWinnerBetID := uint64(3)
	ts.placeTestBet(
		failedWinnerBettorAddr,
		ts.market.UID,
		ts.market.Odds[0].UID,
		failedWinnerBetID,
		sdkmath.NewInt(100000000000),
		ts.betFee,
		types.ErrInternalProcessingBet,
		betOdds,
		oddUIDS,
	)

	///// loser bet placement
	//
	//
	loserBettorAddr := simappUtil.TestParamUsers["user8"].Address
	loserBal := ts.tApp.BankKeeper.GetBalance(ts.ctx, loserBettorAddr, params.DefaultBondDenom)
	loserBetID := uint64(4)
	loserBet, _, loserBetFulfillment := ts.placeTestBet(
		loserBettorAddr,
		ts.market.UID,
		ts.market.Odds[2].UID,
		loserBetID,
		defaultBetAmount,
		ts.betFee,
		nil,
		betOdds,
		oddUIDS,
	)
	loserBet.BetFulfillment = loserBetFulfillment
	ts.tApp.BetKeeper.SetBet(ts.ctx, loserBet, loserBetID)
	loserBalAfterWager := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		loserBettorAddr,
		params.DefaultBondDenom,
	)
	expLoserBalanceAfterWager := loserBal.Amount.Sub(loserBet.Amount).Sub(loserBet.Fee)
	require.Equal(ts.t, expLoserBalanceAfterWager, loserBalAfterWager.Amount)

	return []bettypes.Bet{
		winner1Bet,
		winner2Bet,
		loserBet,
	}, winner1PayoutProfit, winner2PayoutProfit
}

func (ts *testBetSuite) placeTestBet(
	bettorAddr sdk.AccAddress,
	marketUID, oddsUID string,
	betID uint64,
	amount sdkmath.Int,
	fee sdkmath.Int,
	expErr error,
	odds map[string]*bettypes.BetOddsCompact,
	oddUIDS []string,
) (bettypes.Bet, sdk.Dec, []*bettypes.BetFulfillment) {
	bet := bettypes.Bet{
		UID:               uuid.NewString(),
		MarketUID:         marketUID,
		OddsUID:           oddsUID,
		OddsType:          bettypes.OddsType_ODDS_TYPE_DECIMAL,
		OddsValue:         "1.1",
		Amount:            amount,
		Fee:               fee,
		Status:            bettypes.Bet_STATUS_PENDING,
		Creator:           bettorAddr.String(),
		CreatedAt:         cast.ToInt64(ts.ctx.BlockTime().Unix()),
		MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
	}

	payoutProfit, err := bettypes.CalculatePayoutProfit(bet.OddsType, bet.OddsValue, bet.Amount)
	require.NoError(ts.t, err)

	betFeeCollectorBalanceBeforeWager := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		ts.tApp.AccountKeeper.GetModuleAddress(bettypes.BetFeeCollectorFunder{}.GetModuleAcc()),
		params.DefaultBondDenom,
	)
	liquidityPoolBalanceBeforeWager := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		ts.tApp.AccountKeeper.GetModuleAddress(types.OrderBookLiquidityFunder{}.GetModuleAcc()),
		params.DefaultBondDenom,
	)

	betFulfillment, err := ts.k.ProcessWager(
		ts.ctx, bet.UID, bet.MarketUID, bet.OddsUID, bet.MaxLossMultiplier, bet.Amount, payoutProfit,
		bettorAddr, bet.Fee, bet.OddsType, bet.OddsValue, 1, odds, oddUIDS,
	)
	if expErr != nil {
		require.ErrorIs(ts.t, expErr, err)
	} else {
		require.NoError(ts.t, err)
		ts.tApp.BetKeeper.SetBet(ts.ctx, bet, betID)
		ts.tApp.BetKeeper.SetPendingBet(ts.ctx, &bettypes.PendingBet{Creator: bet.Creator, UID: bet.UID}, betID, marketUID)

		betFeeCollectorBalanceAfterWager := ts.tApp.BankKeeper.GetBalance(ts.ctx, ts.tApp.AccountKeeper.GetModuleAddress(bettypes.BetFeeCollectorFunder{}.GetModuleAcc()), params.DefaultBondDenom)
		require.Equal(ts.t, bet.Fee.Int64(), betFeeCollectorBalanceAfterWager.Sub(betFeeCollectorBalanceBeforeWager).Amount.Int64())

		liquidityPoolBalanceAfterWager := ts.tApp.BankKeeper.GetBalance(ts.ctx, ts.tApp.AccountKeeper.GetModuleAddress(types.OrderBookLiquidityFunder{}.GetModuleAcc()), params.DefaultBondDenom)
		require.Equal(ts.t, bet.Amount.Int64(), liquidityPoolBalanceAfterWager.Sub(liquidityPoolBalanceBeforeWager).Amount.Int64())
	}

	return bet, payoutProfit, betFulfillment
}

// 100 Deposits case

func TestWagerLargeNumbers(t *testing.T) {
	ts := newTestBetSuiteForLargeNumbers(t)
	ts.bulkDepositPlaceBetsAndTest()
}

func newTestBetSuiteForLargeNumbers(t *testing.T) testBetSuite {
	tApp, k, ctx := setupKeeperAndApp(t)

	betFee := sdkmath.NewInt(10)

	params := tApp.HouseKeeper.GetParams(ctx)
	params.HouseParticipationFee = sdk.NewDec(0)
	tApp.HouseKeeper.SetParams(ctx, params)

	marketUID := uuid.NewString()
	market := markettypes.Market{
		UID:     marketUID,
		StartTS: cast.ToUint64(time.Now().Unix()),
		EndTS:   cast.ToUint64(time.Now().Add(5 * time.Minute).Unix()),
		Odds: []*markettypes.Odds{
			{UID: uuid.NewString(), Meta: "test odds1"},
			{UID: uuid.NewString(), Meta: "test odds2"},
			{UID: uuid.NewString(), Meta: "test odds3"},
			{UID: uuid.NewString(), Meta: "test odds4"},
			{UID: uuid.NewString(), Meta: "test odds5"},
			{UID: uuid.NewString(), Meta: "test odds6"},
		},
		Status:  markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
		Creator: simappUtil.TestParamUsers["user1"].Address.String(),
		Meta:    "test market",
		BookUID: marketUID,
	}

	deposits := []housetypes.Deposit{}
	for i := 0; i < 100; i++ {
		deposits = append(deposits, housetypes.Deposit{
			DepositorAddress: simappUtil.TestParamUsers["user2"].Address.String(),
			Amount:           sdkmath.NewInt(10000000),
		})
	}

	participations := make([]types.OrderBookParticipation, len(deposits))

	return testBetSuite{t, k, ctx, *tApp, betFee, market, deposits, participations}
}

func (ts *testBetSuite) bulkDepositPlaceBetsAndTest() {
	ts.tApp.MarketKeeper.SetMarket(ts.ctx, ts.market)

	betAmount := sdkmath.NewInt(25354548)

	oddsUIDs := []string{}
	for _, o := range ts.market.Odds {
		oddsUIDs = append(oddsUIDs, o.UID)
	}
	err := ts.k.InitiateOrderBook(ts.ctx, ts.market.UID, oddsUIDs)
	require.NoError(ts.t, err)

	for i := 0; i < len(ts.deposits); i++ {
		found := false
		participationIndex, err := ts.tApp.HouseKeeper.Deposit(
			ts.ctx,
			ts.deposits[i].DepositorAddress,
			ts.deposits[i].DepositorAddress,
			ts.market.BookUID,
			ts.deposits[i].Amount,
		)
		require.NoError(ts.t, err)
		ts.deposits[i], found = ts.tApp.HouseKeeper.GetDeposit(
			ts.ctx,
			ts.deposits[i].DepositorAddress,
			ts.market.UID,
			participationIndex,
		)
		require.True(ts.t, found)
		ts.participations[i], found = ts.k.GetOrderBookParticipation(
			ts.ctx,
			ts.market.UID,
			participationIndex,
		)
		require.True(ts.t, found)
	}

	multipliers := []sdk.Dec{
		sdk.MustNewDecFromStr("0.8649932157"),
		sdk.MustNewDecFromStr("0.8658575751"),
		sdk.MustNewDecFromStr("0.9368906369"),
		sdk.MustNewDecFromStr("0.81411403723"),
		sdk.MustNewDecFromStr("1.00"),
		sdk.MustNewDecFromStr("0.92010393313"),
	}
	betOdds := make(map[string]*bettypes.BetOddsCompact)
	var oddUIDS []string
	for i, odd := range ts.market.Odds {
		betOdds[odd.UID] = &bettypes.BetOddsCompact{UID: odd.UID, MaxLossMultiplier: multipliers[i]}
		oddUIDS = append(oddUIDS, odd.UID)
	}

	bettorAddr := simappUtil.TestParamUsers["user5"].Address
	bet := bettypes.Bet{
		UID:               uuid.NewString(),
		MarketUID:         ts.market.UID,
		OddsUID:           ts.market.Odds[1].UID,
		OddsType:          bettypes.OddsType_ODDS_TYPE_DECIMAL,
		OddsValue:         "4.415",
		Amount:            betAmount,
		Fee:               sdk.ZeroInt(),
		Status:            bettypes.Bet_STATUS_PENDING,
		Creator:           bettorAddr.String(),
		CreatedAt:         cast.ToInt64(ts.ctx.BlockTime().Unix()),
		MaxLossMultiplier: betOdds[ts.market.Odds[1].UID].MaxLossMultiplier,
	}

	payoutProfit, err := bettypes.CalculatePayoutProfit(bet.OddsType, bet.OddsValue, bet.Amount)
	require.NoError(ts.t, err)

	betFulfillment, err := ts.k.ProcessWager(
		ts.ctx, bet.UID, bet.MarketUID, bet.OddsUID, bet.MaxLossMultiplier, bet.Amount, payoutProfit,
		bettorAddr, bet.Fee, bet.OddsType, bet.OddsValue, 1, betOdds, oddUIDS,
	)
	require.NoError(ts.t, err)

	expected := []*bettypes.BetFulfillment{
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 1, BetAmount: sdkmath.NewInt(2535454), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 2, BetAmount: sdkmath.NewInt(2535454), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 3, BetAmount: sdkmath.NewInt(2535454), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 4, BetAmount: sdkmath.NewInt(2535453), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 5, BetAmount: sdkmath.NewInt(2535453), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 6, BetAmount: sdkmath.NewInt(2535453), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 7, BetAmount: sdkmath.NewInt(2535453), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 8, BetAmount: sdkmath.NewInt(2535454), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 9, BetAmount: sdkmath.NewInt(2535453), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 10, BetAmount: sdkmath.NewInt(2535454), PayoutProfit: sdkmath.NewInt(8658575)},
		{ParticipantAddress: simappUtil.TestParamUsers["user2"].Address.String(), ParticipationIndex: 11, BetAmount: sdkmath.NewInt(13), PayoutProfit: sdkmath.NewInt(31)},
	}
	require.Equal(ts.t, expected, betFulfillment)

	// fulfilled bet amount sum should be equal to the actual bet amount
	betAmountSum := sdkmath.NewInt(0)
	for _, bf := range betFulfillment {
		betAmountSum = betAmountSum.Add(bf.BetAmount)
	}
	require.Equal(ts.t, betAmount, betAmountSum)

	// fulfilled bet payout profit sum should be equal to the calculated payout profit
	payoutProfitSum := sdkmath.NewInt(0)
	for _, bf := range betFulfillment {
		payoutProfitSum = payoutProfitSum.Add(bf.PayoutProfit)
	}
	require.Equal(ts.t, payoutProfit.TruncateInt(), payoutProfitSum)
}
