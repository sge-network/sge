package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/testutil/sample"

	"github.com/sge-network/sge/x/reward/types"
)

func TestMsgCreateCampaign_ValidateBasic(t *testing.T) {
	poolBalance := sdkmath.NewInt(100)
	tests := []struct {
		name string
		msg  types.MsgCreateCampaign
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgCreateCampaign{
				Creator:    "invalid_address",
				Uid:        uuid.NewString(),
				Ticket:     "ticket",
				TotalFunds: poolBalance,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid uid",
			msg: types.MsgCreateCampaign{
				Creator:    sample.AccAddress(),
				Uid:        "invalid uid",
				Ticket:     "ticket",
				TotalFunds: poolBalance,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid ticket",
			msg: types.MsgCreateCampaign{
				Creator:    sample.AccAddress(),
				Uid:        uuid.NewString(),
				Ticket:     "",
				TotalFunds: poolBalance,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "negative pool amount",
			msg: types.MsgCreateCampaign{
				Creator: sample.AccAddress(),
				Uid:     uuid.NewString(),
				Ticket:  "ticket",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: types.MsgCreateCampaign{
				Creator:    sample.AccAddress(),
				Uid:        uuid.NewString(),
				Ticket:     "ticket",
				TotalFunds: poolBalance,
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

func TestMsgUpdateCampaign_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateCampaign
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgUpdateCampaign{
				Creator: "invalid_address",
				Uid:     uuid.NewString(),
				Ticket:  "ticket",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid uid",
			msg: types.MsgUpdateCampaign{
				Creator: sample.AccAddress(),
				Uid:     "invalid uid",
				Ticket:  "ticket",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid ticket",
			msg: types.MsgUpdateCampaign{
				Creator: sample.AccAddress(),
				Uid:     uuid.NewString(),
				Ticket:  "",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: types.MsgUpdateCampaign{
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
