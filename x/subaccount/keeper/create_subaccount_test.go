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

func TestMsgServer_CreateSubAccount(t *testing.T) {
	account := sample.AccAddress()
	sender := sample.AccAddress()

	app, _, msgServer, ctx := setupMsgServerAndApp(t)

	err := simapp.FundAccount(app.BankKeeper, ctx, sender, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(100000000))))
	require.NoError(t, err)

	// Check that the account has been created
	require.False(t, app.AccountKeeper.HasAccount(ctx, types.NewAddressFromSubaccount(1)))

	someTime := time.Now().Add(10 * time.Minute)
	msg := &types.MsgCreateSubAccountRequest{
		Sender:          sender.String(),
		SubAccountOwner: account.String(),
		LockedBalances: []*types.LockedBalance{
			{
				UnlockTime: someTime,
				Amount:     sdk.NewInt(123),
			},
		},
	}

	_, err = msgServer.CreateSubAccount(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	// Check that the account has been created
	require.True(t, app.AccountKeeper.HasAccount(ctx, types.NewAddressFromSubaccount(1)))

	// Check that the account has the correct balance
	balance := app.BankKeeper.GetBalance(ctx, types.NewAddressFromSubaccount(1), params.DefaultBondDenom)
	require.Equal(t, sdk.NewInt(123), balance.Amount)

	// Check that we can get the account by owner
	owner := app.SubAccountKeeper.GetSubAccountOwner(ctx, 1)
	require.Equal(t, account, owner)

	// check that balance unlocks are set correctly
	lockedBalances := app.SubAccountKeeper.GetLockedBalances(ctx, types.NewAddressFromSubaccount(1))
	require.Len(t, lockedBalances, 1)
	require.True(t, someTime.Equal(lockedBalances[0].UnlockTime))
	require.Equal(t, sdk.NewInt(123), lockedBalances[0].Amount)
}

func TestMsgServer_CreateSubAccount_Errors(t *testing.T) {
	beforeTime := time.Now().Add(-10 * time.Minute)
	afterTime := time.Now().Add(10 * time.Minute)
	account := sample.AccAddress()
	sender := sample.AccAddress()

	tests := []struct {
		name        string
		msg         types.MsgCreateSubAccountRequest
		prepare     func(ctx sdk.Context, keeper keeper.Keeper)
		expectedErr string
	}{
		{
			name: "unlock time is expired",
			msg: types.MsgCreateSubAccountRequest{
				Sender:          sender.String(),
				SubAccountOwner: account.String(),
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
			name: "account has already sub account",
			msg: types.MsgCreateSubAccountRequest{
				Sender:          sender.String(),
				SubAccountOwner: account.String(),
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: afterTime,
						Amount:     sdk.NewInt(123),
					},
				},
			},
			prepare: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetSubAccountOwner(ctx, 1, account)
			},
			expectedErr: types.ErrSubaccountAlreadyExist.Error(),
		},
		{
			name: "invalid request",
			msg: types.MsgCreateSubAccountRequest{
				Sender:          sender.String(),
				SubAccountOwner: account.String(),
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: afterTime,
						Amount:     sdk.Int{},
					},
				},
			},
			prepare:     func(ctx sdk.Context, k keeper.Keeper) {},
			expectedErr: "invalid request",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, k, msgServer, ctx := setupMsgServerAndApp(t)

			tt.prepare(ctx, k)

			_, err := msgServer.CreateSubAccount(sdk.WrapSDKContext(ctx), &tt.msg)
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
