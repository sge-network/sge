package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/house/types"
	"github.com/stretchr/testify/require"
)

func TestMsgWithdrawValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgWithdraw
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgWithdraw{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: types.MsgWithdraw{
				Creator:            sample.AccAddress(),
				MarketUID:          uuid.NewString(),
				Amount:             sdk.NewInt(100),
				Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
				ParticipationIndex: 1,
				Ticket:             "Ticket",
			},
		},
		{
			name: "invalid withdrawal mode",
			msg: types.MsgWithdraw{
				Creator:   sample.AccAddress(),
				MarketUID: uuid.NewString(),
				Mode:      types.WithdrawalMode_WITHDRAWAL_MODE_UNSPECIFIED,
			},
			err: types.ErrInvalidWithdrawMode,
		},
		{
			name: "invalid participation index",
			msg: types.MsgWithdraw{
				Creator:            sample.AccAddress(),
				MarketUID:          uuid.NewString(),
				Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
				ParticipationIndex: 0,
				Ticket:             "Ticket",
			},
			err: types.ErrInvalidIndex,
		},
		{
			name: "invalid market UID",
			msg: types.MsgWithdraw{
				Creator:   sample.AccAddress(),
				MarketUID: "Invalid UID",
				Mode:      types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
			},
			err: types.ErrInvalidMarketUID,
		},
		{
			name: "invalid amount",
			msg: types.MsgWithdraw{
				Creator:            sample.AccAddress(),
				MarketUID:          uuid.NewString(),
				Amount:             sdk.NewInt(-1),
				Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL,
				ParticipationIndex: 1,
				Ticket:             "Ticket",
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

func TestNewWithdraw(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		expected := types.Withdrawal{
			ID:                 1,
			Creator:            uuid.NewString(),
			MarketUID:          uuid.NewString(),
			Amount:             sdk.NewInt(100),
			ParticipationIndex: 0,
			Address:            sample.AccAddress(),
			Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
		}
		res := types.NewWithdrawal(
			expected.ID,
			expected.Creator,
			expected.Address,
			expected.MarketUID,
			0,
			expected.Amount,
			expected.Mode,
		)
		require.Equal(t, expected, res)
	})
}
