package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/subaccount/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
)

func TestQueryServer(t *testing.T) {
	app, k, msgServer, ctx := setupMsgServerAndApp(t)
	queryServer := keeper.NewQueryServer(*k)

	subAccOwner := simapp.TestParamUsers["user1"].Address
	subAccFunder := simapp.TestParamUsers["user1"].Address
	// setup
	wantParams := types.DefaultParams()
	k.SetParams(ctx, wantParams)

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

	t.Run("Params", func(t *testing.T) {
		gotParams, err := queryServer.Params(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
		require.NoError(t, err)
		require.Equal(t, wantParams, gotParams.Params)
	})

	t.Run("Subaccount", func(t *testing.T) {
		gotSubaccount, err := queryServer.Subaccount(sdk.WrapSDKContext(ctx), &types.QuerySubaccountRequest{
			Address: subAccOwner.String(),
		})
		require.NoError(t, err)
		require.Equal(t, subAccAddr.String(), gotSubaccount.Address)
		require.Equal(t, types.AccountSummary{
			DepositedAmount: subAccFunds,
			SpentAmount:     sdk.ZeroInt(),
			WithdrawnAmount: sdk.ZeroInt(),
			LostAmount:      sdk.ZeroInt(),
		}, gotSubaccount.Balance)
		require.Len(t, gotSubaccount.LockedBalance, 1)
		require.Equal(t, subAccFunds, gotSubaccount.LockedBalance[0].Amount)
	})
}
