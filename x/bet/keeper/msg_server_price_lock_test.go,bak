package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/types"
)

func TestMsgServerFundPriceLock(t *testing.T) {
	tApp, _, msgk, ctx, wctx := setupMsgServerAndApp(t)
	creator := sample.AccAddress()
	funder := simapp.TestParamUsers["user1"]

	topUpAmount := sdkmath.NewInt(1000)

	t.Run("insufficient funder balance", func(t *testing.T) {
		inputMsg := &types.MsgPriceLockPoolTopUp{
			Creator: creator,
			Funder:  creator,
			Amount:  topUpAmount,
		}

		_, err := msgk.PriceLockPoolTopUp(wctx, inputMsg)
		require.ErrorIs(t, types.ErrInsufficientBalanceInPriceLockFunder, err)
	})

	t.Run("authorization not found", func(t *testing.T) {
		inputMsg := &types.MsgPriceLockPoolTopUp{
			Creator: creator,
			Funder:  funder.Address.String(),
			Amount:  topUpAmount,
		}

		_, err := msgk.PriceLockPoolTopUp(wctx, inputMsg)
		require.ErrorIs(t, types.ErrAuthorizationNotFound, err)
	})

	t.Run("Success", func(t *testing.T) {
		inputMsg := &types.MsgPriceLockPoolTopUp{
			Creator: funder.Address.String(),
			Funder:  funder.Address.String(),
			Amount:  topUpAmount,
		}

		_, err := msgk.PriceLockPoolTopUp(wctx, inputMsg)
		require.NoError(t, err)
		balance := tApp.BankKeeper.GetBalance(ctx, tApp.AccountKeeper.GetModuleAddress(types.PriceLockFunder{}.GetModuleAcc()), params.DefaultBondDenom)
		require.Equal(t, topUpAmount, balance.Amount)
	})
}
