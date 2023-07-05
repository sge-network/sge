package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMsgCreateSubAccount_Validate(t *testing.T) {
	sender := sample.NativeAccAddress()
	owner := sample.NativeAccAddress()

	someTime := time.Now()
	tests := []struct {
		name string
		msg  MsgCreateSubAccount
		want error
	}{
		{
			name: "invalid sender",
			msg: MsgCreateSubAccount{
				Sender:          "someInvalidAddress",
				SubAccountOwner: owner.String(),
				LockedBalances: []*LockedBalance{
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
			msg: MsgCreateSubAccount{
				Sender:          sender.String(),
				SubAccountOwner: "someInvalidAddress",
				LockedBalances: []*LockedBalance{
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
			msg: MsgCreateSubAccount{
				Sender:          sender.String(),
				SubAccountOwner: owner.String(),
				LockedBalances: []*LockedBalance{
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
