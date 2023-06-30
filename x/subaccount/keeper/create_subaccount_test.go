package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMsgServer_CreateSubAccount(t *testing.T) {

}

func TestMsgServer_CreateSubAccount_Errors(t *testing.T) {
	beforeTime := time.Now().Add(-10 * time.Minute)
	tests := []struct {
		name        string
		msg         types.MsgCreateSubAccountRequest
		expectedErr error
	}{
		{
			name: "unlock time is expired",
			msg: types.MsgCreateSubAccountRequest{
				Sender:          "someSender",
				SubAccountOwner: "someAccountOwner",
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: &beforeTime,
						Amount:     sdk.Int{},
					},
				},
			},
			expectedErr: types.ErrUnlockTokenTimeExpired,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, _, msgServer, ctx := setupMsgServerAndApp(t)

			_, err := msgServer.CreateSubAccount(sdk.WrapSDKContext(ctx), &tt.msg)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}
