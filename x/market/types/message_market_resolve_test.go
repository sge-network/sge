package types_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/market/types"
)

func TestMsgResolveValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgResolve
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgResolve{
				Creator: "invalid_address",
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: types.MsgResolve{
				Creator: sample.AccAddress(),
				Ticket:  "Ticket",
			},
		},
		{
			name: "no ticket",
			msg: types.MsgResolve{
				Creator: sample.AccAddress(),
			},
			err: sdkerrtypes.ErrInvalidRequest,
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

func TestNewResolve(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		expected := &types.MsgResolve{
			Creator: uuid.NewString(),
			Ticket:  "Ticket",
		}
		res := types.NewMsgResolve(
			expected.Creator,
			expected.Ticket,
		)
		require.Equal(t, expected, res)
	})
}
