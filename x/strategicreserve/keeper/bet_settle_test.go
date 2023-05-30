package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/require"
)

func TestBetSettlement(t *testing.T) {
	ts := NewTestSuite(t)

	bets, winner1PayoutProfit, winner2PayoutProfit := ts.placeBetsAndTest()

	winner1BalAfterPlacement := ts.tApp.BankKeeper.GetBalance(ts.ctx, sdk.MustAccAddressFromBech32(bets[0].Creator), params.DefaultBondDenom).Amount
	winner2BalAfterPlacement := ts.tApp.BankKeeper.GetBalance(ts.ctx, sdk.MustAccAddressFromBech32(bets[1].Creator), params.DefaultBondDenom).Amount
	loserBalanceAfterPlacement := ts.tApp.BankKeeper.GetBalance(ts.ctx, sdk.MustAccAddressFromBech32(bets[2].Creator), params.DefaultBondDenom).Amount

	// resolve market
	//
	ts.tApp.MarketKeeper.ResolveMarket(ts.ctx, ts.market, &markettypes.MarketResolutionTicketPayload{
		UID:            ts.market.UID,
		ResolutionTS:   ts.market.StartTS + 10,
		WinnerOddsUIDs: []string{ts.market.Odds[0].UID, ts.market.Odds[1].UID},
		Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
	})

	// settle all of the resolved market
	err := ts.tApp.BetKeeper.BatchMarketSettlements(ts.ctx)
	require.NoError(t, err)

	winner1BettorBalAfterSettlement := ts.tApp.BankKeeper.GetBalance(ts.ctx, sdk.MustAccAddressFromBech32(bets[0].Creator), params.DefaultBondDenom)
	expWinner1BalanceAfterSettlement := winner1BalAfterPlacement.Add(bets[0].Amount).Add(winner1PayoutProfit.TruncateInt())
	require.Equal(t, expWinner1BalanceAfterSettlement.Int64(), winner1BettorBalAfterSettlement.Amount.Int64())

	winner2BettorBalAfterSettlement := ts.tApp.BankKeeper.GetBalance(ts.ctx, sdk.MustAccAddressFromBech32(bets[1].Creator), params.DefaultBondDenom)
	expWinner2BalanceAfterSettlement := winner2BalAfterPlacement.Add(bets[1].Amount).Add(winner2PayoutProfit.TruncateInt())
	require.Equal(t, expWinner2BalanceAfterSettlement.Int64(), winner2BettorBalAfterSettlement.Amount.Int64())

	loserBettorBalAfterSettlement := ts.tApp.BankKeeper.GetBalance(ts.ctx, sdk.MustAccAddressFromBech32(bets[2].Creator), params.DefaultBondDenom)
	require.Equal(t, loserBalanceAfterPlacement.Int64(), loserBettorBalAfterSettlement.Amount.Int64())
}
