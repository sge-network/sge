package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/app/params"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/stretchr/testify/require"
)

func TestFund(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)

	senderAddr := simappUtil.TestParamUsers["user1"].Address
	initialBalance := tApp.BankKeeper.GetBalance(ctx, senderAddr, params.DefaultBondDenom)
	successAmount := sdk.NewInt(1000)

	for _, tc := range []struct {
		desc   string
		amount sdk.Int

		err error
	}{
		{
			desc:   "not enough balance",
			amount: sdk.NewInt(100000000000000),
			err:    types.ErrInsufficientAccountBalance,
		},
		{
			desc:   "success",
			amount: successAmount,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := k.Fund(
				types.OrderBookLiquidityFunder{},
				ctx,
				senderAddr,
				tc.amount,
			)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				balance := tApp.BankKeeper.GetBalance(ctx, senderAddr, params.DefaultBondDenom)
				require.Equal(t, initialBalance.Sub(balance).Amount, successAmount)
			}
		})
	}
}

func TestReFund(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)

	successAmount := sdk.NewInt(1000)
	err := k.Fund(
		types.OrderBookLiquidityFunder{},
		ctx,
		simappUtil.TestParamUsers["user2"].Address,
		successAmount,
	)
	require.NoError(t, err)

	receiverAddr := simappUtil.TestParamUsers["user1"].Address
	initialBalance := tApp.BankKeeper.GetBalance(ctx, receiverAddr, params.DefaultBondDenom)

	for _, tc := range []struct {
		desc   string
		amount sdk.Int

		err error
	}{
		{
			desc:   "not enough balance",
			amount: successAmount.Add(sdk.NewInt(1)),
			err:    types.ErrInsufficientBalanceInModuleAccount,
		},
		{
			desc:   "success",
			amount: successAmount,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := k.ReFund(
				types.OrderBookLiquidityFunder{},
				ctx,
				receiverAddr,
				tc.amount,
			)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				balance := tApp.BankKeeper.GetBalance(ctx, receiverAddr, params.DefaultBondDenom)
				require.Equal(t, balance.Sub(initialBalance).Amount, successAmount)
			}
		})
	}
}
