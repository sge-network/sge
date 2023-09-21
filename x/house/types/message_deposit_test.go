package types_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/house/types"
)

func TestMsgDepositValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgDeposit
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgDeposit{
				Creator: "invalid_address",
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			name: "valid",
			msg: types.MsgDeposit{
				Creator:   sample.AccAddress(),
				MarketUID: uuid.NewString(),
				Amount:    sdkmath.NewInt(100),
				Ticket:    "Ticket",
			},
		},
		{
			name: "invalid market UID",
			msg: types.MsgDeposit{
				Creator:   sample.AccAddress(),
				MarketUID: "Invalid UID",
			},
			err: types.ErrInvalidMarketUID,
		},
		{
			name: "invalid amount",
			msg: types.MsgDeposit{
				Creator:   sample.AccAddress(),
				MarketUID: uuid.NewString(),
				Amount:    sdk.ZeroInt(),
				Ticket:    "Ticket",
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

func TestNewDeposit(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		expected := types.Deposit{
			Creator:               uuid.NewString(),
			DepositorAddress:      uuid.NewString(),
			MarketUID:             uuid.NewString(),
			Amount:                sdkmath.NewInt(100),
			ParticipationIndex:    0,
			WithdrawalCount:       0,
			TotalWithdrawalAmount: sdk.ZeroInt(),
		}
		res := types.NewDeposit(
			expected.Creator,
			expected.DepositorAddress,
			expected.MarketUID,
			expected.Amount,
			expected.TotalWithdrawalAmount,
			expected.WithdrawalCount,
		)
		require.Equal(t, expected, res)
	})
}
