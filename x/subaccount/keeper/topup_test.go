package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMsgServerTopUp_HappyPath(t *testing.T) {

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
