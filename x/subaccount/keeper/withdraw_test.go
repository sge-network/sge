package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMsgServer_WithdrawUnlockedBalances(t *testing.T) {
	sender := sample.NativeAccAddress()
	subaccountOwner := sample.NativeAccAddress()
	lockedTime := time.Now().Add(time.Hour * 24 * 365)
	lockedTime2 := time.Now().Add(time.Hour * 24 * 365 * 2)

	app, _, msgServer, ctx := setupMsgServerAndApp(t)

	t.Log("fund sender account")
	err := simapp.FundAccount(app.BankKeeper, ctx, sender, sdk.NewCoins(sdk.NewInt64Coin("usge", 1000)))

	t.Log("Create sub account")
	_, err = msgServer.CreateSubAccount(sdk.WrapSDKContext(ctx), &types.MsgCreateSubAccount{
		Sender:          sender.String(),
		SubAccountOwner: subaccountOwner.String(),
		LockedBalances: []*types.LockedBalance{
			{
				Amount:     sdk.NewInt(100),
				UnlockTime: lockedTime,
			},
			{
				Amount:     sdk.NewInt(200),
				UnlockTime: lockedTime2,
			},
		},
	})
	require.NoError(t, err)

	t.Log("check balance of sub account")
	subAccountAddr := types.NewAddressFromSubaccount(1)
	balance := app.BankKeeper.GetBalance(ctx, subAccountAddr, "usge")
	require.Equal(t, sdk.NewInt(300), balance.Amount)

	t.Log("check balance of subaccount owner")
	balance = app.BankKeeper.GetBalance(ctx, subaccountOwner, "usge")
	require.Equal(t, sdk.NewInt(0), balance.Amount)

	t.Log("Withdraw unlocked balances, with 0 expires")
	_, err = msgServer.WithdrawUnlockedBalances(sdk.WrapSDKContext(ctx), &types.MsgWithdrawUnlockedBalances{
		Sender: subaccountOwner.String(),
	})
	require.NoError(t, err)

	t.Log("balance of subaccount owner should be zero")
	balance = app.BankKeeper.GetBalance(ctx, subaccountOwner, "usge")
	require.True(t, balance.IsZero())

	t.Log("expire first locked balance")
	ctx = ctx.WithBlockTime(lockedTime.Add(1 * time.Second))
	t.Log("Withdraw unlocked balances, with 1 expires")
	_, err = msgServer.WithdrawUnlockedBalances(sdk.WrapSDKContext(ctx), &types.MsgWithdrawUnlockedBalances{
		Sender: subaccountOwner.String(),
	})
	require.NoError(t, err)

	t.Log("balance of subaccount owner should be the same as first locked balance")
	balance = app.BankKeeper.GetBalance(ctx, subaccountOwner, "usge")
	require.True(t, balance.Amount.Equal(sdk.NewInt(100)), balance.Amount.String())

	t.Log("expire second locked balance")
	ctx = ctx.WithBlockTime(lockedTime2.Add(1 * time.Second))
	t.Log("Withdraw unlocked balances, with 2 expires")
	_, err = msgServer.WithdrawUnlockedBalances(sdk.WrapSDKContext(ctx), &types.MsgWithdrawUnlockedBalances{
		Sender: subaccountOwner.String(),
	})
	require.NoError(t, err)

	t.Log("balance of subaccount owner should be the same as both expired locked balances")
	balance = app.BankKeeper.GetBalance(ctx, subaccountOwner, "usge")
	require.True(t, balance.Amount.Equal(sdk.NewInt(300)), balance.Amount.String())

	t.Log("check bank balance of sub account address")
	balance = app.BankKeeper.GetBalance(ctx, subAccountAddr, "usge")
	require.True(t, balance.IsZero())

	t.Log("check balance of subaccount")
	subaccountBalance := app.SubaccountKeeper.GetBalance(ctx, 1)
	require.True(t, subaccountBalance.WithdrawmAmount.Equal(sdk.NewInt(300)))
	require.True(t, subaccountBalance.DepositedAmount.Equal(sdk.NewInt(300)))
	require.True(t, subaccountBalance.SpentAmount.IsZero())
	require.True(t, subaccountBalance.LostAmount.IsZero())
}

func TestMsgServer_WithdrawUnlockedBalances_Errors(t *testing.T) {
	sender := sample.AccAddress()
	tests := []struct {
		name        string
		msg         types.MsgWithdrawUnlockedBalances
		prepare     func(ctx sdk.Context, keeper keeper.Keeper)
		expectedErr string
	}{
		{
			name: "sub account does not exist",
			msg: types.MsgWithdrawUnlockedBalances{
				Sender: sender,
			},
			prepare:     func(ctx sdk.Context, keeper keeper.Keeper) {},
			expectedErr: types.ErrSubaccountDoesNotExist.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, k, msgServer, ctx := setupMsgServerAndApp(t)

			tt.prepare(ctx, k)

			_, err := msgServer.WithdrawUnlockedBalances(sdk.WrapSDKContext(ctx), &tt.msg)
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
