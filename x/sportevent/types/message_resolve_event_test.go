package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgResolveEvent_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgResolveEvent
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgResolveEvent{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid ticket",
			msg: MsgResolveEvent{
				Creator: sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid request",
			msg: MsgResolveEvent{
				Creator: sample.AccAddress(),
				Ticket:  "contains data",
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
