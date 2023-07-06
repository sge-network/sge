package types_test

import (
	"fmt"
	"github.com/sge-network/sge/x/subaccount/types"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateSubAccount_Validate(t *testing.T) {
	sender := sample.NativeAccAddress()
	owner := sample.NativeAccAddress()

	someTime := time.Now()
	tests := []struct {
		name string
		msg  types.MsgCreateSubAccount
		want error
	}{
		{
			name: "invalid sender",
			msg: types.MsgCreateSubAccount{
				Sender:          "someInvalidAddress",
				SubAccountOwner: owner.String(),
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: someTime,
						Amount:     sdk.NewInt(123),
					},
				},
			},
			want: errors.ErrInvalidAddress,
		},
		{
			name: "invalid sub account owner",
			msg: types.MsgCreateSubAccount{
				Sender:          sender.String(),
				SubAccountOwner: "someInvalidAddress",
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: someTime,
						Amount:     sdk.NewInt(123),
					},
				},
			},
			want: errors.ErrInvalidAddress,
		},
		{
			name: "unlock time zero",
			msg: types.MsgCreateSubAccount{
				Sender:          sender.String(),
				SubAccountOwner: owner.String(),
				LockedBalances: []*types.LockedBalance{
					{
						UnlockTime: time.Time{},
						Amount:     sdk.NewInt(123),
					},
				},
			},
			want: fmt.Errorf("invalid locked balance: unlock time is zero"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msg.ValidateBasic()
			require.ErrorContains(t, got, tt.want.Error())
		})
	}
}
