package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/testutil/simapp"
	bettypes "github.com/sge-network/sge/x/bet/types"
	housetypes "github.com/sge-network/sge/x/house/types"
	"github.com/sge-network/sge/x/strategicreserve"
	"github.com/sge-network/sge/x/strategicreserve/types"
	"github.com/stretchr/testify/require"
)

var (
	houseFeeSpend = sdk.NewCoins(sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(1000)))
	betFeeSpend   = sdk.NewCoins(sdk.NewCoin(params.BaseCoinUnit, sdk.NewInt(100)))
)

func testProposal() *types.DataFeeCollectorFeedProposal {
	return types.NewDataFeeCollectorFeedProposal(
		"Test",
		"description",
		houseFeeSpend.AmountOf(params.BaseCoinUnit),
		betFeeSpend.AmountOf(params.BaseCoinUnit),
	)
}

func TestProposalHandlerPassed(t *testing.T) {
	tApp, _, ctx := setupKeeperAndApp(t)

	// add coins to the module account
	houseFeeMacc := tApp.AccountKeeper.GetModuleAccount(ctx, housetypes.HouseFeeCollector)
	require.NoError(t, simapp.FundModuleAccount(tApp, ctx, houseFeeMacc.GetName(), houseFeeSpend))
	tApp.AccountKeeper.SetModuleAccount(ctx, houseFeeMacc)

	betFeeMacc := tApp.AccountKeeper.GetModuleAccount(ctx, bettypes.BetFeeCollector)
	require.NoError(t, simapp.FundModuleAccount(tApp, ctx, betFeeMacc.GetName(), betFeeSpend))
	tApp.AccountKeeper.SetModuleAccount(ctx, betFeeMacc)

	dataFeeMacc := tApp.AccountKeeper.GetModuleAccount(ctx, types.DataFeeCollector)
	require.True(t, tApp.BankKeeper.GetBalance(ctx, dataFeeMacc.GetAddress(), params.BaseCoinUnit).IsZero())

	tp := testProposal()
	hdlr := strategicreserve.NewDataFeeCollectorFeedProposalHandler(tApp.StrategicReserveKeeper)
	require.NoError(t, hdlr(ctx, tp))

	balance := tApp.BankKeeper.GetBalance(ctx, dataFeeMacc.GetAddress(), params.BaseCoinUnit)
	require.Equal(t,
		balance.Amount,
		houseFeeSpend.AmountOf(params.BaseCoinUnit).Add(betFeeSpend.AmountOf(params.BaseCoinUnit)),
	)
}

func TestProposalHandlerFailed(t *testing.T) {
	tApp, _, ctx := setupKeeperAndApp(t)

	// add coins to the module account
	houseFeeMacc := tApp.AccountKeeper.GetModuleAccount(ctx, housetypes.HouseFeeCollector)
	require.NoError(t, simapp.FundModuleAccount(tApp, ctx, houseFeeMacc.GetName(), houseFeeSpend))
	tApp.AccountKeeper.SetModuleAccount(ctx, houseFeeMacc)

	dataFeeMacc := tApp.AccountKeeper.GetModuleAccount(ctx, types.DataFeeCollector)
	require.True(t, tApp.BankKeeper.GetBalance(ctx, dataFeeMacc.GetAddress(), params.BaseCoinUnit).IsZero())

	tp := testProposal()
	hdlr := strategicreserve.NewDataFeeCollectorFeedProposalHandler(tApp.StrategicReserveKeeper)
	require.ErrorIs(t, hdlr(ctx, tp), types.ErrInsufficientBalanceInModuleAccount)
}
