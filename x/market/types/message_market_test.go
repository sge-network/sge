package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/require"
)

func TestMsgAddMarketValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgAddMarket
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgAddMarket{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: types.MsgAddMarket{
				Creator: sample.AccAddress(),
				Ticket:  "Ticket",
			},
		},
		{
			name: "no ticket",
			msg: types.MsgAddMarket{
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
		expected := &types.MsgAddMarket{
			Creator: uuid.NewString(),
			Ticket:  "Ticket",
		}
		res := types.NewMsgAddMarket(
			expected.Creator,
			expected.Ticket,
		)
		require.Equal(t, expected, res)
	})
}

func TestMsgUpdateMarketValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateMarket
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgUpdateMarket{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: types.MsgUpdateMarket{
				Creator: sample.AccAddress(),
				Ticket:  "Ticket",
			},
		},
		{
			name: "no ticket",
			msg: types.MsgUpdateMarket{
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
		expected := &types.MsgUpdateMarket{
			Creator: uuid.NewString(),
			Ticket:  "Ticket",
		}
		res := types.NewMsgUpdateMarket(
			expected.Creator,
			expected.Ticket,
		)
		require.Equal(t, expected, res)
	})
}
