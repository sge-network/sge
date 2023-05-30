package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/app/params"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	bettypes "github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/strategicreserve/keeper"
	"github.com/sge-network/sge/x/strategicreserve/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestBetPlacement(t *testing.T) {
	ts := NewTestSuite(t)
	ts.placeBetsAndTest()
}

type TestSuite struct {
	t                *testing.T
	k                *keeper.KeeperTest
	ctx              sdk.Context
	tApp             simappUtil.TestApp
	betFee           sdk.Int
	participationFee sdk.Int
	market           markettypes.Market
}

func NewTestSuite(t *testing.T) TestSuite {
	tApp, k, ctx := setupKeeperAndApp(t)

	betFee := sdk.NewInt(10)
	participationFee := sdk.NewInt(10)

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

	return TestSuite{t, k, ctx, *tApp, betFee, participationFee, market}
}

func (ts *TestSuite) placeBetsAndTest() ([]bettypes.Bet, sdk.Dec, sdk.Dec) {

	ts.tApp.MarketKeeper.SetMarket(ts.ctx, ts.market)

	err := ts.k.InitiateOrderBook(ts.ctx, ts.market.UID, []string{
		ts.market.Odds[0].UID,
		ts.market.Odds[1].UID,
		ts.market.Odds[2].UID,
	})
	require.NoError(ts.t, err)

	_, err = ts.k.InitiateOrderBookParticipation(ts.ctx, simappUtil.TestParamUsers["user2"].Address, ts.market.BookUID, sdk.NewInt(8000), ts.participationFee)
	require.NoError(ts.t, err)

	_, err = ts.k.InitiateOrderBookParticipation(ts.ctx, simappUtil.TestParamUsers["user3"].Address, ts.market.BookUID, sdk.NewInt(10000), ts.participationFee)
	require.NoError(ts.t, err)

	oddsExposures, found := ts.k.GetOrderBookOddsExposure(ts.ctx, ts.market.BookUID, ts.market.Odds[0].UID)
	require.True(ts.t, found)
	require.Equal(ts.t, []uint64{1, 2}, oddsExposures.FulfillmentQueue)

	///// winner1 bet placement
	//
	//
	winner1BettorAddr := simappUtil.TestParamUsers["user4"].Address
	winner1Bal := ts.tApp.BankKeeper.GetBalance(ts.ctx, winner1BettorAddr, params.DefaultBondDenom)
	winner1BetID := uint64(1)
	winner1Bet, winner1PayoutProfit, winner1BetFulfillment := ts.placeTestBet(winner1BettorAddr, ts.market.UID, ts.market.Odds[0].UID, winner1BetID, ts.betFee, nil)
	winner1Bet.BetFulfillment = winner1BetFulfillment
	ts.tApp.BetKeeper.SetBet(ts.ctx, winner1Bet, winner1BetID)
	winner1BalAfterPlacement := ts.tApp.BankKeeper.GetBalance(ts.ctx, winner1BettorAddr, params.DefaultBondDenom)
	expWinner1BalanceAfterPlacement := winner1Bal.Amount.Sub(winner1Bet.Amount).Sub(winner1Bet.BetFee)
	require.Equal(ts.t, expWinner1BalanceAfterPlacement.Int64(), winner1BalAfterPlacement.Amount.Int64())

	oddsExposures, found = ts.k.GetOrderBookOddsExposure(ts.ctx, ts.market.BookUID, ts.market.Odds[0].UID)
	require.True(ts.t, found)
	require.Equal(ts.t, []uint64{2}, oddsExposures.FulfillmentQueue)

	///// winner2 bet placement
	//
	//
	winner2BettorAddr := simappUtil.TestParamUsers["user5"].Address
	winner2Bal := ts.tApp.BankKeeper.GetBalance(ts.ctx, winner2BettorAddr, params.DefaultBondDenom)
	winner2BetID := uint64(2)
	winner2Bet, winner2PayoutProfit, winner2BetFulfillment := ts.placeTestBet(winner2BettorAddr, ts.market.UID, ts.market.Odds[0].UID, winner2BetID, ts.betFee, nil)

	// oddsExposures, found := k.GetOrderBookOddsExposure(ctx, market.BookUID, oddsUID1)
	// require.True(t, found)

	// require.Equal(t, 1, len(oddsExposures.FulfillmentQueue))

	winner2Bet.BetFulfillment = winner2BetFulfillment
	ts.tApp.BetKeeper.SetBet(ts.ctx, winner2Bet, winner2BetID)
	winner2BalAfterPlacement := ts.tApp.BankKeeper.GetBalance(ts.ctx, winner1BettorAddr, params.DefaultBondDenom)
	expWinner2BalanceAfterPlacement := winner2Bal.Amount.Sub(winner2Bet.Amount).Sub(winner2Bet.BetFee)
	require.Equal(ts.t, expWinner2BalanceAfterPlacement.Int64(), winner2BalAfterPlacement.Amount.Int64())
	oddsExposures, found = ts.k.GetOrderBookOddsExposure(ts.ctx, ts.market.BookUID, ts.market.Odds[0].UID)
	require.True(ts.t, found)
	require.Equal(ts.t, 0, len(oddsExposures.FulfillmentQueue))

	///// failed winner bet placement
	// should fail because there is not participation to fulfill this bet.
	//
	failedWinnerBettorAddr := simappUtil.TestParamUsers["user6"].Address
	failedWinnerBetID := uint64(3)
	ts.placeTestBet(failedWinnerBettorAddr, ts.market.UID, ts.market.Odds[0].UID, failedWinnerBetID, ts.betFee, types.ErrInternalProcessingBet)

	///// loser bet placement
	//
	//
	loserBettorAddr := simappUtil.TestParamUsers["user7"].Address
	loserBal := ts.tApp.BankKeeper.GetBalance(ts.ctx, loserBettorAddr, params.DefaultBondDenom)
	loserBetID := uint64(3)
	loserBet, _, loserBetFulfillment := ts.placeTestBet(loserBettorAddr, ts.market.UID, ts.market.Odds[2].UID, loserBetID, ts.betFee, nil)
	loserBet.BetFulfillment = loserBetFulfillment
	ts.tApp.BetKeeper.SetBet(ts.ctx, loserBet, loserBetID)
	loserBalAfterPlacement := ts.tApp.BankKeeper.GetBalance(ts.ctx, loserBettorAddr, params.DefaultBondDenom)
	expLoserBalanceAfterPlacement := loserBal.Amount.Sub(loserBet.Amount).Sub(loserBet.BetFee)
	require.Equal(ts.t, expLoserBalanceAfterPlacement, loserBalAfterPlacement.Amount)

	return []bettypes.Bet{
		winner1Bet,
		winner2Bet,
		loserBet,
	}, winner1PayoutProfit, winner2PayoutProfit
}

func (ts *TestSuite) placeTestBet(bettorAddr sdk.AccAddress, marketUID, oddsUID string, betID uint64, betFee sdk.Int, expErr error) (bettypes.Bet, sdk.Dec, []*bettypes.BetFulfillment) {
	bet := bettypes.Bet{
		UID:               uuid.NewString(),
		MarketUID:         marketUID,
		OddsUID:           oddsUID,
		OddsType:          bettypes.OddsType_ODDS_TYPE_DECIMAL,
		OddsValue:         "1.5",
		Amount:            sdk.NewInt(100),
		BetFee:            betFee,
		Status:            bettypes.Bet_STATUS_PENDING,
		Creator:           bettorAddr.String(),
		CreatedAt:         cast.ToInt64(ts.ctx.BlockTime().Unix()),
		MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
	}
	ts.tApp.BetKeeper.SetBet(ts.ctx, bet, betID)
	ts.tApp.BetKeeper.SetPendingBet(ts.ctx, &bettypes.PendingBet{Creator: bet.Creator, UID: bet.UID}, betID, marketUID)

	payoutProfit, err := bettypes.CalculatePayoutProfit(bet.OddsType, bet.OddsValue, bet.Amount)
	require.NoError(ts.t, err)

	betFulfillment, err := ts.k.ProcessBetPlacement(
		ts.ctx, bet.UID, bet.MarketUID, bet.OddsUID, bet.MaxLossMultiplier, bet.Amount, payoutProfit,
		bettorAddr, bet.BetFee, bet.OddsType, bet.OddsValue, 1,
	)
	if expErr != nil {
		require.ErrorIs(ts.t, expErr, err)
	} else {
		require.NoError(ts.t, err)
	}

	return bet, payoutProfit, betFulfillment
}
