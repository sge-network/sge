package types

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateSubAccountRequest_Validate(t *testing.T) {
	sender := sample.AccAddress()
	owner := sample.AccAddress()

	someTime := time.Now()
	tests := []struct {
		name string
		msg  MsgCreateSubAccountRequest
		want error
	}{
		{
			name: "invalid sender",
			msg: MsgCreateSubAccountRequest{
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
			msg: MsgCreateSubAccountRequest{
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
			msg: MsgCreateSubAccountRequest{
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
			got := tt.msg.Validate()
			require.ErrorContains(t, got, tt.want.Error())
		})
	}
}

func TestLockedBalanceValidate(t *testing.T) {
	tests := []struct {
		name string
		lb   LockedBalance
		want error
	}{
		{
			name: "unlock time zero",
			lb: LockedBalance{
				UnlockTime: time.Time{},
				Amount:     sdk.Int{},
			},
			want: fmt.Errorf("unlock time is zero"),
		},
		{
			name: "negative amount",
			lb: LockedBalance{
				UnlockTime: time.Now(),
				Amount:     sdk.NewInt(-1),
			},
			want: fmt.Errorf("amount is negative"),
		},
		{
			name: "nil amount",
			lb: LockedBalance{
				UnlockTime: time.Now(),
				Amount:     sdk.Int{},
			},
			want: fmt.Errorf("amount is nil"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lb.Validate()
			require.Equal(t, tt.want, got)
		})
	}
}
