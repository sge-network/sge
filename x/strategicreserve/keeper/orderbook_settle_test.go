package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
	"github.com/stretchr/testify/require"
)

func TestOrderBookSettlement(t *testing.T) {
	ts := newTestBetSuite(t)
	participant1BalanceBeforeDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.deposits[0].DepositorAddress),
		params.DefaultBondDenom).Amount

	participant2BalanceBeforeDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.deposits[1].DepositorAddress),
		params.DefaultBondDenom).Amount

	participant3BalanceBeforeDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.deposits[2].DepositorAddress),
		params.DefaultBondDenom).Amount

	bets, winner1PayoutProfit, winner2PayoutProfit := ts.placeBetsAndTest()
	ts.settleBetsAndTest(bets, winner1PayoutProfit, winner2PayoutProfit)

	participant1BalanceAfterDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.deposits[0].DepositorAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		ts.deposits[0].Amount,
		participant1BalanceBeforeDeposit.Sub(participant1BalanceAfterDeposit))

	participant2BalanceAfterDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.deposits[1].DepositorAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		ts.deposits[1].Amount,
		participant2BalanceBeforeDeposit.Sub(participant2BalanceAfterDeposit))

	participant3BalanceAfterDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.deposits[2].DepositorAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		ts.deposits[2].Amount,
		participant3BalanceBeforeDeposit.Sub(participant3BalanceAfterDeposit))

	err := ts.k.BatchOrderBookSettlements(ts.ctx)
	require.NoError(t, err)

	participant1BalanceAfterSettlement := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.deposits[0].DepositorAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		participant1BalanceBeforeDeposit.
			// subtract first winner payoutprofit
			Sub(winner1PayoutProfit.TruncateInt()).
			// subtract participation fee
			Sub(ts.deposits[0].Fee).
			// add loser bet amount
			Add(bets[2].Amount),
		participant1BalanceAfterSettlement)

	participant2BalanceAfterSettlement := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.deposits[1].DepositorAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		participant2BalanceBeforeDeposit.
			// subtract second winner payoutprofit
			Sub(winner2PayoutProfit.TruncateInt()).
			// subtract participation fee
			Sub(ts.deposits[1].Fee),
		participant2BalanceAfterSettlement)

	participant3BalanceAfterSettlement := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.deposits[2].DepositorAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		participant3BalanceBeforeDeposit.Int64(),
		participant3BalanceAfterSettlement.Int64())
}
