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

func TestMsgServer_CreateSubAccount(t *testing.T) {

}

func TestMsgServer_CreateSubAccount_Errors(t *testing.T) {
	beforeTime := time.Now().Add(-10 * time.Minute)
	afterTime := time.Now().Add(10 * time.Minute)
	account := sample.AccAddress()

	tests := []struct {
		name        string
		msg         types.MsgCreateSubAccountRequest
		prepare     func(ctx sdk.Context, keeper keeper.Keeper)
		expectedErr error
	}{
		{
			name: "unlock time is expired",
			msg: types.MsgCreateSubAccountRequest{
				Sender:          "someSender",
				SubAccountOwner: account.String(),
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: &beforeTime,
						Amount:     sdk.Int{},
					},
				},
			},
			prepare:     func(ctx sdk.Context, k keeper.Keeper) {},
			expectedErr: types.ErrUnlockTokenTimeExpired,
		},
		{
			name: "account has already sub account",
			msg: types.MsgCreateSubAccountRequest{
				Sender:          "someSender",
				SubAccountOwner: account.String(),
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: &afterTime,
						Amount:     sdk.Int{},
					},
				},
			},
			prepare: func(ctx sdk.Context, k keeper.Keeper) {
				k.SetSubAccountOwner(ctx, 1, account)
			},
			expectedErr: types.ErrSubaccountAlreadyExist,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, k, msgServer, ctx := setupMsgServerAndApp(t)

			tt.prepare(ctx, k)

			_, err := msgServer.CreateSubAccount(sdk.WrapSDKContext(ctx), &tt.msg)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}
