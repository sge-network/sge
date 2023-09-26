package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/reward/types"
	"github.com/stretchr/testify/require"
)

func TestMsgApplyReward_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgApplyReward
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgApplyReward{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid uid",
			msg: types.MsgApplyReward{
				Creator:     sample.AccAddress(),
				CampaignUid: "bad uid",
				Ticket:      "ticket",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid ticket",
			msg: types.MsgApplyReward{
				Creator:     sample.AccAddress(),
				CampaignUid: uuid.NewString(),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: types.MsgApplyReward{
				Creator:     sample.AccAddress(),
				CampaignUid: uuid.NewString(),
				Ticket:      "ticket",
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
