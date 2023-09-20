package types_test

import (
	"testing"
	"time"

	"github.com/sge-network/sge/x/subaccount/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateValidateBasic(t *testing.T) {
	creatorAddr := sample.NativeAccAddress()
	owner := sample.NativeAccAddress()

	someTime := uint64(time.Now().Unix())
	tests := []struct {
		name string
		msg  types.MsgCreate
		want error
	}{
		{
			name: "invalid creator",
			msg: types.MsgCreate{
				Creator:         "someInvalidAddress",
				SubAccountOwner: owner.String(),
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: someTime,
						Amount:   sdk.NewInt(123),
					},
				},
			},
			want: errors.ErrInvalidAddress,
		},
		{
			name: "invalid sub account owner",
			msg: types.MsgCreate{
				Creator:         creatorAddr.String(),
				SubAccountOwner: "someInvalidAddress",
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: someTime,
						Amount:   sdk.NewInt(123),
					},
				},
			},
			want: errors.ErrInvalidAddress,
		},
		{
			name: "unlock time zero",
			msg: types.MsgCreate{
				Creator:         creatorAddr.String(),
				SubAccountOwner: owner.String(),
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: 0,
						Amount:   sdk.NewInt(123),
					},
				},
			},
			want: types.ErrInvalidLockedBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msg.ValidateBasic()
			require.ErrorContains(t, got, tt.want.Error())
		})
	}
}
