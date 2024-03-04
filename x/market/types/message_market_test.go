package types_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/market/types"
)

func TestMsgAddValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgAdd
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgAdd{
				Creator: "invalid_address",
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: types.MsgAdd{
				Creator: sample.AccAddress(),
				Ticket:  "Ticket",
			},
		},
		{
			name: "no ticket",
			msg: types.MsgAdd{
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

func TestNewAdd(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		expected := &types.MsgAdd{
			Creator: uuid.NewString(),
			Ticket:  "Ticket",
		}
		res := types.NewMsgAdd(
			expected.Creator,
			expected.Ticket,
		)
		require.Equal(t, expected, res)
	})
}

func TestMsgUpdateValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdate
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgUpdate{
				Creator: "invalid_address",
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: types.MsgUpdate{
				Creator: sample.AccAddress(),
				Ticket:  "Ticket",
			},
		},
		{
			name: "no ticket",
			msg: types.MsgUpdate{
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

func TestNewUpdate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		expected := &types.MsgUpdate{
			Creator: uuid.NewString(),
			Ticket:  "Ticket",
		}
		res := types.NewMsgUpdate(
			expected.Creator,
			expected.Ticket,
		)
		require.Equal(t, expected, res)
	})
}
