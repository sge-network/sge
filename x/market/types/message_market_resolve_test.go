package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/require"
)

func TestMsgResolveMarketValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgResolveMarket
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgResolveMarket{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: types.MsgResolveMarket{
				Creator: sample.AccAddress(),
				Ticket:  "Ticket",
			},
		},
		{
			name: "no ticket",
			msg: types.MsgResolveMarket{
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

func TestNewResolveMarket(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		expected := &types.MsgResolveMarket{
			Creator: uuid.NewString(),
			Ticket:  "Ticket",
		}
		res := types.NewMsgResolveMarket(
			expected.Creator,
			expected.Ticket,
		)
		require.Equal(t, expected, res)
	})
}
