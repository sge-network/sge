package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
	bettypes "github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

func TestBetSettlement(t *testing.T) {
	ts := newTestBetSuite(t)
	bets, winner1PayoutProfit, winner2PayoutProfit := ts.placeBetsAndTest()
	ts.settleBetsAndTest(bets, winner1PayoutProfit, winner2PayoutProfit)
}

func (ts *testBetSuite) settleBetsAndTest(
	bets []bettypes.Bet,
	winner1PayoutProfit, winner2PayoutProfit sdkmath.LegacyDec,
) {
	winner1BalAfterWager := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		sdk.MustAccAddressFromBech32(bets[0].Creator),
		params.DefaultBondDenom,
	).Amount
	winner2BalAfterWager := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		sdk.MustAccAddressFromBech32(bets[1].Creator),
		params.DefaultBondDenom,
	).Amount
	loserBalanceAfterWager := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		sdk.MustAccAddressFromBech32(bets[2].Creator),
		params.DefaultBondDenom,
	).Amount

	// resolve market
	//
	ts.tApp.MarketKeeper.Resolve(ts.ctx, ts.market, &markettypes.MarketResolutionTicketPayload{
		UID:            ts.market.UID,
		ResolutionTS:   ts.market.StartTS + 10,
		WinnerOddsUIDs: []string{ts.market.Odds[0].UID, ts.market.Odds[1].UID},
		Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
	})

	// settle all the resolved market
	err := ts.tApp.BetKeeper.BatchMarketSettlements(ts.ctx)
	require.NoError(ts.t, err)

	winner1BettorBalAfterSettlement := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		sdk.MustAccAddressFromBech32(bets[0].Creator),
		params.DefaultBondDenom,
	)
	expWinner1BalanceAfterSettlement := winner1BalAfterWager.Add(bets[0].Amount).
		Add(winner1PayoutProfit.TruncateInt())
	require.Equal(
		ts.t,
		expWinner1BalanceAfterSettlement.Int64(),
		winner1BettorBalAfterSettlement.Amount.Int64(),
	)

	winner2BettorBalAfterSettlement := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		sdk.MustAccAddressFromBech32(bets[1].Creator),
		params.DefaultBondDenom,
	)
	expWinner2BalanceAfterSettlement := winner2BalAfterWager.Add(bets[1].Amount).
		Add(winner2PayoutProfit.TruncateInt())
	require.Equal(
		ts.t,
		expWinner2BalanceAfterSettlement.Int64(),
		winner2BettorBalAfterSettlement.Amount.Int64(),
	)

	loserBettorBalAfterSettlement := ts.tApp.BankKeeper.GetBalance(
		ts.ctx,
		sdk.MustAccAddressFromBech32(bets[2].Creator),
		params.DefaultBondDenom,
	)
	require.Equal(
		ts.t,
		loserBalanceAfterWager.Int64(),
		loserBettorBalAfterSettlement.Amount.Int64(),
	)
}
