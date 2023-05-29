package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/app/params"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/strategicreserve/keeper"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

type TestSuite struct {
	t      *testing.T
	k      *keeper.KeeperTest
	ctx    sdk.Context
	tApp   simappUtil.TestApp
	betFee sdk.Int
}

func TestBetPlacement(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)

	marketUID := uuid.NewString()
	oddsUID1 := uuid.NewString()
	oddsUID2 := uuid.NewString()
	oddsUID3 := uuid.NewString()

	participationAmount := sdk.NewInt(10000)
	betFee := sdk.NewInt(10)
	participationFee := sdk.NewInt(10)

	ts := TestSuite{t, k, ctx, *tApp, betFee}

	market := markettypes.Market{
		UID:     marketUID,
		StartTS: cast.ToUint64(time.Now().Unix()),
		EndTS:   cast.ToUint64(time.Now().Add(5 * time.Minute).Unix()),
		Odds: []*markettypes.Odds{
			{UID: oddsUID1, Meta: "test odds1"},
			{UID: oddsUID2, Meta: "test odds2"},
			{UID: oddsUID3, Meta: "test odds3"},
		},
		Status:  markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
		Creator: simappUtil.TestParamUsers["user1"].Address.String(),
		Meta:    "test market",
		BookUID: marketUID,
	}
	tApp.MarketKeeper.SetMarket(ctx, market)

	err := k.InitiateOrderBook(ctx, marketUID, []string{oddsUID1, oddsUID2, oddsUID3})
	require.NoError(t, err)

	_, err = k.InitiateOrderBookParticipation(ctx, simappUtil.TestParamUsers["user2"].Address, market.BookUID, participationAmount, participationFee)
	require.NoError(t, err)

	_, err = k.InitiateOrderBookParticipation(ctx, simappUtil.TestParamUsers["user3"].Address, market.BookUID, participationAmount, participationFee)
	require.NoError(t, err)

	///// winner bet placement
	//
	//
	winnerBettorAddr := simappUtil.TestParamUsers["user4"].Address
	winnerBal := tApp.BankKeeper.GetBalance(ctx, winnerBettorAddr, params.DefaultBondDenom)
	winnerBetID := uint64(1)
	winnerBet, winnerPayoutProfit, winnerBetFulfillment := ts.placeTestBet(winnerBettorAddr, marketUID, oddsUID1, winnerBetID, betFee)

	// oddsExposures, found := k.GetOrderBookOddsExposure(ctx, market.BookUID, oddsUID1)
	// require.True(t, found)

	// require.Equal(t, 1, len(oddsExposures.FulfillmentQueue))

	winnerBet.BetFulfillment = winnerBetFulfillment
	tApp.BetKeeper.SetBet(ctx, winnerBet, winnerBetID)

	winnerBalAfterPlacement := tApp.BankKeeper.GetBalance(ctx, winnerBettorAddr, params.DefaultBondDenom)
	expWinnerBalanceAfterPlacement := winnerBal.Amount.Sub(winnerBet.Amount).Sub(winnerBet.BetFee)
	require.Equal(t, expWinnerBalanceAfterPlacement.Int64(), winnerBalAfterPlacement.Amount.Int64())

	///// loser bet placement
	//
	//
	loserBettorAddr := simappUtil.TestParamUsers["user5"].Address
	loserBal := tApp.BankKeeper.GetBalance(ctx, loserBettorAddr, params.DefaultBondDenom)
	loserBetID := uint64(2)
	loserBet, _, loserBetFulfillment := ts.placeTestBet(loserBettorAddr, marketUID, oddsUID3, loserBetID, betFee)

	loserBet.BetFulfillment = loserBetFulfillment
	tApp.BetKeeper.SetBet(ctx, loserBet, loserBetID)

	loserBalAfterPlacement := tApp.BankKeeper.GetBalance(ctx, loserBettorAddr, params.DefaultBondDenom)
	expLoserBalanceAfterPlacement := loserBal.Amount.Sub(loserBet.Amount).Sub(loserBet.BetFee)
	require.Equal(t, expLoserBalanceAfterPlacement, loserBalAfterPlacement.Amount)

	// resolve market
	//
	tApp.MarketKeeper.ResolveMarket(ctx, market, &markettypes.MarketResolutionTicketPayload{
		UID:            marketUID,
		ResolutionTS:   market.StartTS + 10,
		WinnerOddsUIDs: []string{oddsUID1, oddsUID2},
		Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
	})

	// settle all of the resolved market
	err = tApp.BetKeeper.BatchMarketSettlements(ctx)
	require.NoError(t, err)

	winnerBettorBalAfterSettlement := tApp.BankKeeper.GetBalance(ctx, winnerBettorAddr, params.DefaultBondDenom)
	expWinnerBalanceAfterSettlement := expWinnerBalanceAfterPlacement.Add(winnerBet.Amount).Add(winnerPayoutProfit.TruncateInt())
	require.Equal(t, expWinnerBalanceAfterSettlement.Int64(), winnerBettorBalAfterSettlement.Amount.Int64())

	loserBettorBalAfterSettlement := tApp.BankKeeper.GetBalance(ctx, loserBettorAddr, params.DefaultBondDenom)
	require.Equal(t, expLoserBalanceAfterPlacement.Int64(), loserBettorBalAfterSettlement.Amount.Int64())

}

func (ts *TestSuite) placeTestBet(bettorAddr sdk.AccAddress, marketUID, oddsUID string, betID uint64, betFee sdk.Int) (types.Bet, sdk.Dec, []*types.BetFulfillment) {
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
	require.NoError(ts.t, err)

	return bet, payoutProfit, betFulfillment
}
