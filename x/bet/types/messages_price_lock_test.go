package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/bet/types"
)

func TestMsgPriceLockTopUpFundsValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgPriceLockPoolTopUp
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgPriceLockPoolTopUp{
				Creator: "invalid_address",
				Funder:  sample.AccAddress(),
				Amount:  sdkmath.OneInt(),
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			name: "invalid funder",
			msg: types.MsgPriceLockPoolTopUp{
				Creator: sample.AccAddress(),
				Funder:  "invalid_address",
				Amount:  sdkmath.OneInt(),
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			name: "invalid amount",
			msg: types.MsgPriceLockPoolTopUp{
				Creator: sample.AccAddress(),
				Funder:  sample.AccAddress(),
				Amount:  sdkmath.ZeroInt(),
			},
			err: types.ErrInvalidPriceLockTopUpAmount,
		},
		{
			name: "valid bet message",
			msg: types.MsgPriceLockPoolTopUp{
				Creator: sample.AccAddress(),
				Funder:  sample.AccAddress(),
				Amount:  sdkmath.OneInt(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
