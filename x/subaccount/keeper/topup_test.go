package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMsgServerTopUp_HappyPath(t *testing.T) {
	afterTime := time.Now().Add(10 * time.Minute)
	sender := sample.NativeAccAddress()
	subaccount := sample.AccAddress()

	app, k, msgServer, ctx := setupMsgServerAndApp(t)

	// Fund sender
	err := simapp.FundAccount(app.BankKeeper, ctx, sender, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(100000000))))
	require.NoError(t, err)

	// Create subaccount
	msg := &types.MsgCreateSubAccount{
		Sender:          sender.String(),
		SubAccountOwner: subaccount,
		LockedBalances:  []*types.LockedBalance{},
	}
	_, err = msgServer.CreateSubAccount(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	balance := k.GetBalance(ctx, 1)
	require.Equal(t, sdk.NewInt(0), balance.DepositedAmount)
	balances := k.GetLockedBalances(ctx, 1)
	require.Len(t, balances, 0)

	msgTopUp := &types.MsgTopUp{
		Sender:     sender.String(),
		SubAccount: subaccount,
		LockedBalances: []*types.LockedBalance{
			{
				UnlockTime: afterTime,
				Amount:     sdk.NewInt(123),
			},
		},
	}
	_, err = msgServer.TopUp(sdk.WrapSDKContext(ctx), msgTopUp)
	require.NoError(t, err)

	// Check balance
	balance = k.GetBalance(ctx, 1)
	require.Equal(t, sdk.NewInt(123), balance.DepositedAmount)
	balances = k.GetLockedBalances(ctx, 1)
	require.Len(t, balances, 1)
	require.True(t, afterTime.Equal(balances[0].UnlockTime))
	require.Equal(t, sdk.NewInt(123), balances[0].Amount)
}

func TestNewMsgServerTopUp_Errors(t *testing.T) {
	beforeTime := time.Now().Add(-10 * time.Minute)
	afterTime := time.Now().Add(10 * time.Minute)

	sender := sample.AccAddress()
	subaccount := sample.AccAddress()

	tests := []struct {
		name        string
		msg         types.MsgTopUp
		prepare     func(ctx sdk.Context, keeper keeper.Keeper)
		expectedErr string
	}{
		{
			name: "unlock time is expired",
			msg: types.MsgTopUp{
				Sender:     sender,
				SubAccount: subaccount,
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: beforeTime,
						Amount:     sdk.NewInt(123),
					},
				},
			},
			prepare:     func(ctx sdk.Context, k keeper.Keeper) {},
			expectedErr: types.ErrUnlockTokenTimeExpired.Error(),
		},
		{
			name: "sub account does not exist",
			msg: types.MsgTopUp{
				Sender:     sender,
				SubAccount: subaccount,
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: afterTime,
						Amount:     sdk.NewInt(123),
					},
				},
			},
			prepare:     func(ctx sdk.Context, k keeper.Keeper) {},
			expectedErr: types.ErrSubaccountDoesNotExist.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, k, msgServer, ctx := setupMsgServerAndApp(t)

			tt.prepare(ctx, k)

			_, err := msgServer.TopUp(sdk.WrapSDKContext(ctx), &tt.msg)
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
