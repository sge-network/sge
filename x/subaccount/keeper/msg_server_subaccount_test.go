package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
)

func TestMsgServer_Create(t *testing.T) {
	account := sample.NativeAccAddress()
	creatorAddr := sample.NativeAccAddress()

	app, _, msgServer, ctx := setupMsgServerAndApp(t)

	err := testutil.FundAccount(app.BankKeeper, ctx, creatorAddr, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdkmath.NewInt(100000000))))
	require.NoError(t, err)

	// Check that the account has been created
	require.False(t, app.AccountKeeper.HasAccount(ctx, types.NewAddressFromSubaccount(1)))

	someTime := uint64(time.Now().Add(10 * time.Minute).Unix())
	msg := &types.MsgCreate{
		Creator: creatorAddr.String(),
		Owner:   account.String(),
		LockedBalances: []types.LockedBalance{
			{
				UnlockTS: someTime,
				Amount:   sdkmath.NewInt(123),
			},
		},
	}

	_, err = msgServer.Create(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	// Check that the account has been created
	require.True(t, app.AccountKeeper.HasAccount(ctx, types.NewAddressFromSubaccount(1)))

	// Check that the account has the correct balance
	balance := app.BankKeeper.GetBalance(ctx, types.NewAddressFromSubaccount(1), params.DefaultBondDenom)
	require.Equal(t, sdkmath.NewInt(123), balance.Amount)

	// Check that we can get the account by owner
	owner, exists := app.SubaccountKeeper.GetSubaccountOwner(ctx, types.NewAddressFromSubaccount(1))
	require.True(t, exists)
	require.Equal(t, account, owner)

	// check that balance unlocks are set correctly
	lockedBalances, _ := app.SubaccountKeeper.GetBalances(ctx, types.NewAddressFromSubaccount(1), types.LockedBalanceStatus_LOCKED_BALANCE_STATUS_LOCKED)
	require.Len(t, lockedBalances, 1)
	require.True(t, someTime == lockedBalances[0].UnlockTS)
	require.Equal(t, sdkmath.NewInt(123), lockedBalances[0].Amount)

	// get the balance of the account
	subaccountBalance, exists := app.SubaccountKeeper.GetAccountSummary(ctx, types.NewAddressFromSubaccount(1))
	require.True(t, exists)
	require.Equal(t, sdk.ZeroInt(), subaccountBalance.SpentAmount)
	require.Equal(t, sdk.ZeroInt(), subaccountBalance.LostAmount)
	require.Equal(t, sdk.ZeroInt(), subaccountBalance.WithdrawnAmount)
	require.Equal(t, sdkmath.NewInt(123), subaccountBalance.DepositedAmount)
}

func TestMsgServer_CreateSubaccount_Errors(t *testing.T) {
	beforeTime := uint64(time.Now().Add(-10 * time.Minute).Unix())
	afterTime := uint64(time.Now().Add(10 * time.Minute).Unix())
	account := sample.NativeAccAddress()
	creatorAddr := sample.NativeAccAddress()

	tests := []struct {
		name        string
		msg         types.MsgCreate
		prepare     func(ctx sdk.Context, keeper *keeper.Keeper)
		expectedErr string
	}{
		{
			name: "unlock time is expired",
			msg: types.MsgCreate{
				Creator: creatorAddr.String(),
				Owner:   account.String(),
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: beforeTime,
						Amount:   sdkmath.NewInt(123),
					},
				},
			},
			prepare:     func(ctx sdk.Context, k *keeper.Keeper) {},
			expectedErr: types.ErrUnlockTokenTimeExpired.Error(),
		},
		{
			name: "account has already sub account",
			msg: types.MsgCreate{
				Creator: creatorAddr.String(),
				Owner:   account.String(),
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: afterTime,
						Amount:   sdkmath.NewInt(123),
					},
				},
			},
			prepare: func(ctx sdk.Context, k *keeper.Keeper) {
				k.SetSubaccountOwner(ctx, types.NewAddressFromSubaccount(1), account)
			},
			expectedErr: types.ErrSubaccountAlreadyExist.Error(),
		},
		{
			name: "invalid request",
			msg: types.MsgCreate{
				Creator: creatorAddr.String(),
				Owner:   account.String(),
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: afterTime,
						Amount:   sdkmath.Int{},
					},
				},
			},
			prepare:     func(ctx sdk.Context, k *keeper.Keeper) {},
			expectedErr: "invalid request",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, k, msgServer, ctx := setupMsgServerAndApp(t)

			tt.prepare(ctx, k)

			_, err := msgServer.Create(sdk.WrapSDKContext(ctx), &tt.msg)
			require.ErrorContains(t, err, tt.expectedErr)
		})
	}
}
