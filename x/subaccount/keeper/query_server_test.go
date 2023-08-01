package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
	"github.com/stretchr/testify/require"
)

func TestQueryServer(t *testing.T) {
	app, k, msgServer, ctx := setupMsgServerAndApp(t)
	queryServer := keeper.NewQueryServer(k)

	// setup
	wantParams := types.DefaultParams()
	k.SetParams(ctx, wantParams)

	// do subaccount creation
	require.NoError(
		t,
		simapp.FundAccount(
			app.BankKeeper,
			ctx,
			subAccFunder,
			sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).LockedBalanceDenom, subAccFunds)),
		),
	)

	_, err := msgServer.CreateSubAccount(sdk.WrapSDKContext(ctx), &types.MsgCreateSubAccount{
		Sender:          subAccFunder.String(),
		SubAccountOwner: subAccOwner.String(),
		LockedBalances: []*types.LockedBalance{
			{
				UnlockTime: time.Now().Add(24 * time.Hour),
				Amount:     subAccFunds,
			},
		},
	})
	require.NoError(t, err)

	t.Run("Params", func(t *testing.T) {
		gotParams, err := queryServer.Params(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
		require.NoError(t, err)
		require.Equal(t, wantParams, *gotParams.Params)
	})

	t.Run("Subaccount", func(t *testing.T) {
		gotSubaccount, err := queryServer.Subaccount(sdk.WrapSDKContext(ctx), &types.QuerySubaccountRequest{
			SubaccountOwner: subAccOwner.String(),
		})
		require.NoError(t, err)
		require.Equal(t, subAccAddr.String(), gotSubaccount.SubaccountAddress)
		require.Equal(t, types.Balance{
			DepositedAmount: subAccFunds,
			SpentAmount:     sdk.ZeroInt(),
			WithdrawmAmount: sdk.ZeroInt(),
			LostAmount:      sdk.ZeroInt(),
		}, gotSubaccount.Balance)
		require.Len(t, gotSubaccount.LockedBalance, 1)
		require.Equal(t, subAccFunds, gotSubaccount.LockedBalance[0].Amount)
	})
}
