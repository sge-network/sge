package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/bet/types"
)

func TestBetFieldsValidation(t *testing.T) {
	tcs := []struct {
		desc string
		bet  *types.WagerProps
		err  error
	}{
		{
			desc: "space in UID",
			bet: &types.WagerProps{
				UID:    " ",
				Amount: sdkmath.NewInt(int64(10)),
				Ticket: "Ticket",
			},
			err: types.ErrInvalidBetUID,
		},
		{
			desc: "invalid UID",
			bet: &types.WagerProps{
				UID:    "invalidUID",
				Amount: sdkmath.NewInt(int64(10)),
				Ticket: "Ticket",
			},
			err: types.ErrInvalidBetUID,
		},
		{
			desc: "invalid amount",
			bet: &types.WagerProps{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdkmath.NewInt(int64(-1)),
				Ticket: "Ticket",
			},
			err: types.ErrInvalidAmount,
		},
		{
			desc: "empty amount",
			bet: &types.WagerProps{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Ticket: "Ticket",
			},
			err: types.ErrInvalidAmount,
		},
		{
			desc: "space in ticket",
			bet: &types.WagerProps{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdkmath.NewInt(int64(10)),
				Ticket: " ",
			},
			err: types.ErrInvalidTicket,
		},
		{
			desc: "valid message",
			bet: &types.WagerProps{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdkmath.NewInt(int64(10)),
				Ticket: "Ticket",
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := types.WagerValidation(tc.bet)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
