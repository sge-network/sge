package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
)

func TestMsgServer_WithdrawUnlockedBalances(t *testing.T) {
	creatorAddr := sample.NativeAccAddress()
	subaccountOwner := sample.NativeAccAddress()
	lockedTime := time.Now().Add(time.Hour * 24 * 365).UTC()
	lockedTime2 := time.Now().Add(time.Hour * 24 * 365 * 2).UTC()

	app, _, msgServer, ctx := setupMsgServerAndApp(t)

	t.Log("funder account")
	err := testutil.FundAccount(app.BankKeeper, ctx, creatorAddr, sdk.NewCoins(sdk.NewInt64Coin("usge", 1000)))
	require.NoError(t, err)

	t.Log("Create sub account")
	_, err = msgServer.Create(sdk.WrapSDKContext(ctx), &types.MsgCreate{
		Creator:         creatorAddr.String(),
		SubAccountOwner: subaccountOwner.String(),
		LockedBalances: []types.LockedBalance{
			{
				Amount:   sdk.NewInt(100),
				UnlockTS: uint64(lockedTime.Unix()),
			},
			{
				Amount:   sdk.NewInt(200),
				UnlockTS: uint64(lockedTime2.Unix()),
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
		Creator: subaccountOwner.String(),
	})
	require.ErrorContains(t, err, types.ErrNothingToWithdraw.Error())

	t.Log("balance of subaccount owner should be zero")
	balance = app.BankKeeper.GetBalance(ctx, subaccountOwner, "usge")
	require.True(t, balance.IsZero())

	t.Log("expire first locked balance")
	ctx = ctx.WithBlockTime(lockedTime.Add(1 * time.Second))
	t.Log("Withdraw unlocked balances, with 1 expires")
	_, err = msgServer.WithdrawUnlockedBalances(sdk.WrapSDKContext(ctx), &types.MsgWithdrawUnlockedBalances{
		Creator: subaccountOwner.String(),
	})
	require.NoError(t, err)

	t.Log("balance of subaccount owner should be the same as first locked balance")
	balance = app.BankKeeper.GetBalance(ctx, subaccountOwner, "usge")
	require.True(t, balance.Amount.Equal(sdk.NewInt(100)), balance.Amount.String())

	t.Log("expire second locked balance, also force money to be spent")
	// we force some money to be spent on the subaccount to correctly test
	// that if the amount is unlocked but spent, it will not be withdrawable.
	subAccountAddress := types.NewAddressFromSubaccount(1)
	subaccountBalance, exists := app.SubaccountKeeper.GetBalance(ctx, subAccountAddress)
	require.True(t, exists)
	require.NoError(t, subaccountBalance.Spend(sdk.NewInt(100)))
	app.SubaccountKeeper.SetBalance(ctx, subAccountAddr, subaccountBalance)

	ctx = ctx.WithBlockTime(lockedTime2.Add(1 * time.Second))
	t.Log("Withdraw unlocked balances, with 2 expires")
	_, err = msgServer.WithdrawUnlockedBalances(sdk.WrapSDKContext(ctx), &types.MsgWithdrawUnlockedBalances{
		Creator: subaccountOwner.String(),
	})
	require.NoError(t, err)

	t.Log("balance of subaccount owner should be the same as both expired locked balances minus spent money")
	balance = app.BankKeeper.GetBalance(ctx, subaccountOwner, "usge")
	require.Equal(t, sdk.NewInt(200), balance.Amount)

	t.Log("check bank balance of sub account address")
	balance = app.BankKeeper.GetBalance(ctx, subAccountAddr, "usge")
	require.Equal(t, sdk.NewInt(100), balance.Amount)

	t.Log("after unspending the money of the subaccount, the owner will be able to get the money back when withdrawing")
	subaccountBalance, exists = app.SubaccountKeeper.GetBalance(ctx, subAccountAddress)
	require.True(t, exists)
	require.NoError(t, subaccountBalance.Unspend(sdk.NewInt(100)))
	app.SubaccountKeeper.SetBalance(ctx, subAccountAddr, subaccountBalance)
	_, err = msgServer.WithdrawUnlockedBalances(sdk.WrapSDKContext(ctx), &types.MsgWithdrawUnlockedBalances{
		Creator: subaccountOwner.String(),
	})
	require.NoError(t, err)

	// check balances
	balance = app.BankKeeper.GetBalance(ctx, subAccountAddr, "usge")
	require.Equal(t, sdk.NewInt(0), balance.Amount)
	subaccountBalance, exists = app.SubaccountKeeper.GetBalance(ctx, subAccountAddress)
	require.True(t, exists)
	require.Equal(t, sdk.NewInt(300), subaccountBalance.WithdrawmAmount)

	// check that the owner has received the last money
	balance = app.BankKeeper.GetBalance(ctx, subaccountOwner, "usge")
	require.Equal(t, sdk.NewInt(300), balance.Amount)

	// check that the owner can't withdraw again
	_, err = msgServer.WithdrawUnlockedBalances(sdk.WrapSDKContext(ctx), &types.MsgWithdrawUnlockedBalances{
		Creator: subaccountOwner.String(),
	})
	require.ErrorContains(t, err, types.ErrNothingToWithdraw.Error())
}

func TestMsgServer_WithdrawUnlockedBalances_Errors(t *testing.T) {
	creatorAddr := sample.AccAddress()
	tests := []struct {
		name        string
		msg         types.MsgWithdrawUnlockedBalances
		prepare     func(ctx sdk.Context, keeper keeper.Keeper)
		expectedErr string
	}{
		{
			name: "sub account does not exist",
			msg: types.MsgWithdrawUnlockedBalances{
				Creator: creatorAddr,
			},
			prepare:     func(ctx sdk.Context, keeper keeper.Keeper) {},
			expectedErr: types.ErrSubaccountDoesNotExist.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, k, msgServer, ctx := setupMsgServerAndApp(t)

			tt.prepare(ctx, *k)

			_, err := msgServer.WithdrawUnlockedBalances(sdk.WrapSDKContext(ctx), &tt.msg)
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestMsgServerTopUp_HappyPath(t *testing.T) {
	afterTime := uint64(time.Now().Add(10 * time.Minute).Unix())
	creatirAddr := sample.NativeAccAddress()
	subaccount := sample.AccAddress()

	app, k, msgServer, ctx := setupMsgServerAndApp(t)

	// Funder
	err := testutil.FundAccount(app.BankKeeper, ctx, creatirAddr, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(100000000))))
	require.NoError(t, err)

	// Create subaccount
	msg := &types.MsgCreate{
		Creator:         creatirAddr.String(),
		SubAccountOwner: subaccount,
		LockedBalances:  []types.LockedBalance{},
	}
	_, err = msgServer.Create(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	subAccountAddr := types.NewAddressFromSubaccount(1)

	balance, exists := k.GetBalance(ctx, subAccountAddr)
	require.True(t, exists)
	require.Equal(t, sdk.NewInt(0), balance.DepositedAmount)
	balances := k.GetLockedBalances(ctx, subAccountAddr)
	require.Len(t, balances, 0)

	msgTopUp := &types.MsgTopUp{
		Creator:    creatirAddr.String(),
		SubAccount: subaccount,
		LockedBalances: []types.LockedBalance{
			{
				UnlockTS: afterTime,
				Amount:   sdk.NewInt(123),
			},
		},
	}
	_, err = msgServer.TopUp(sdk.WrapSDKContext(ctx), msgTopUp)
	require.NoError(t, err)

	// Check balance
	balance, exists = k.GetBalance(ctx, subAccountAddr)
	require.True(t, exists)
	require.Equal(t, sdk.NewInt(123), balance.DepositedAmount)
	balances = k.GetLockedBalances(ctx, subAccountAddr)
	require.Len(t, balances, 1)
	require.True(t, afterTime == balances[0].UnlockTS)
	require.Equal(t, sdk.NewInt(123), balances[0].Amount)
}

func TestNewMsgServerTopUp_Errors(t *testing.T) {
	beforeTime := uint64(time.Now().Add(-10 * time.Minute).Unix())
	afterTime := uint64(time.Now().Add(10 * time.Minute).Unix())

	creatorAddr := sample.AccAddress()
	subaccount := sample.AccAddress()

	tests := []struct {
		name        string
		msg         types.MsgTopUp
		prepare     func(ctx sdk.Context, msgServer types.MsgServer)
		expectedErr string
	}{
		{
			name: "unlock time is expired",
			msg: types.MsgTopUp{
				Creator:    creatorAddr,
				SubAccount: subaccount,
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: beforeTime,
						Amount:   sdk.NewInt(123),
					},
				},
			},
			prepare:     func(ctx sdk.Context, msgServer types.MsgServer) {},
			expectedErr: types.ErrUnlockTokenTimeExpired.Error(),
		},
		{
			name: "sub account does not exist",
			msg: types.MsgTopUp{
				Creator:    creatorAddr,
				SubAccount: subaccount,
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: afterTime,
						Amount:   sdk.NewInt(123),
					},
				},
			},
			prepare:     func(ctx sdk.Context, msgServer types.MsgServer) {},
			expectedErr: types.ErrSubaccountDoesNotExist.Error(),
		},
		{
			name: "creator has not enough balance",
			msg: types.MsgTopUp{
				Creator:    creatorAddr,
				SubAccount: subaccount,
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: afterTime,
						Amount:   sdk.NewInt(123),
					},
				},
			},
			prepare: func(ctx sdk.Context, msgServer types.MsgServer) {
				// Create subaccount
				msg := &types.MsgCreate{
					Creator:         creatorAddr,
					SubAccountOwner: subaccount,
					LockedBalances:  []types.LockedBalance{},
				}
				_, err := msgServer.Create(sdk.WrapSDKContext(ctx), msg)
				require.NoError(t, err)
			},
			expectedErr: "0usge is smaller than 123usge: insufficient funds",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, _, msgServer, ctx := setupMsgServerAndApp(t)

			tt.prepare(ctx, msgServer)

			_, err := msgServer.TopUp(sdk.WrapSDKContext(ctx), &tt.msg)
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
