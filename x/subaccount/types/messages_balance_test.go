package types_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/types"
)

func TestMsgTopUpValidateBasic(t *testing.T) {
	creatorAddr := sample.NativeAccAddress()
	owner := sample.NativeAccAddress()

	someTime := uint64(time.Now().Unix())
	tests := []struct {
		name string
		msg  types.MsgTopUp
		want error
	}{
		{
			name: "invalid creator",
			msg: types.MsgTopUp{
				Creator: "someInvalidAddress",
				Address: owner.String(),
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: someTime,
						Amount:   sdkmath.NewInt(123),
					},
				},
			},
			want: errors.ErrInvalidAddress,
		},
		{
			name: "invalid sub account owner",
			msg: types.MsgTopUp{
				Creator: creatorAddr.String(),
				Address: "someInvalidAddress",
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: someTime,
						Amount:   sdkmath.NewInt(123),
					},
				},
			},
			want: errors.ErrInvalidAddress,
		},
		{
			name: "unlock time zero",
			msg: types.MsgTopUp{
				Creator: creatorAddr.String(),
				Address: owner.String(),
				LockedBalances: []types.LockedBalance{
					{
						UnlockTS: 0,
						Amount:   sdkmath.NewInt(123),
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
