package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
	"github.com/stretchr/testify/require"
)

func TestOrderBookSettlement(t *testing.T) {
	ts := NewTestBetSuite(t)
	participant1BalanceBeforeDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.participations[0].ParticipantAddress),
		params.DefaultBondDenom).Amount

	participant2BalanceBeforeDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.participations[1].ParticipantAddress),
		params.DefaultBondDenom).Amount

	participant3BalanceBeforeDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.participations[2].ParticipantAddress),
		params.DefaultBondDenom).Amount

	bets, winner1PayoutProfit, winner2PayoutProfit := ts.placeBetsAndTest()
	ts.settleBetsAndTest(bets, winner1PayoutProfit, winner2PayoutProfit)

	participant1BalanceAfterDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.participations[0].ParticipantAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		ts.participationFee.Add(ts.participations[0].Liquidity),
		participant1BalanceBeforeDeposit.Sub(participant1BalanceAfterDeposit))

	participant2BalanceAfterDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.participations[1].ParticipantAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		ts.participationFee.Add(ts.participations[1].Liquidity),
		participant2BalanceBeforeDeposit.Sub(participant2BalanceAfterDeposit))

	participant3BalanceAfterDeposit := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.participations[2].ParticipantAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		ts.participationFee.Add(ts.participations[2].Liquidity),
		participant3BalanceBeforeDeposit.Sub(participant3BalanceAfterDeposit))

	err := ts.k.BatchOrderBookSettlements(ts.ctx)
	require.NoError(t, err)

	participant1BalanceAfterSettlement := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.participations[0].ParticipantAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		participant2BalanceBeforeDeposit.
			// subtract first winner payoutprofit
			Sub(winner1PayoutProfit.TruncateInt()).
			// subtract participation fee
			Sub(ts.participationFee).
			// add loser bet amount
			Add(bets[2].Amount),
		participant1BalanceAfterSettlement)

	participant2BalanceAfterSettlement := ts.tApp.BankKeeper.GetBalance(
		ts.ctx, sdk.MustAccAddressFromBech32(ts.participations[1].ParticipantAddress),
		params.DefaultBondDenom).Amount
	require.Equal(t,
		participant2BalanceBeforeDeposit.
			// subtract second winner payoutprofit
			Sub(winner2PayoutProfit.TruncateInt()).
			// subtract participation fee
			Sub(ts.participationFee),
		participant2BalanceAfterSettlement)
}
