package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/require"
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
			err: sdkerrors.ErrInvalidAddress,
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
			err: sdkerrors.ErrInvalidRequest,
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

func TestNewAddMarket(t *testing.T) {
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

func TestMsgUpdateMarketValidateBasic(t *testing.T) {
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
			err: sdkerrors.ErrInvalidAddress,
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
			err: sdkerrors.ErrInvalidRequest,
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

func TestNewUpdateMarket(t *testing.T) {
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
