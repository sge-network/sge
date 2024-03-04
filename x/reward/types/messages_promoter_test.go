package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/reward/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSetPromoterConf_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgSetPromoterConf
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgSetPromoterConf{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid campaign uid",
			msg: types.MsgSetPromoterConf{
				Creator: sample.AccAddress(),
				Uid:     "bad uid",
				Ticket:  "ticket",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid ticket",
			msg: types.MsgSetPromoterConf{
				Creator: sample.AccAddress(),
				Uid:     uuid.NewString(),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: types.MsgSetPromoterConf{
				Creator: sample.AccAddress(),
				Uid:     uuid.NewString(),
				Ticket:  "ticket",
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
