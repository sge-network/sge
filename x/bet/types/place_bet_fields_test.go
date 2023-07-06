package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/stretchr/testify/require"
)

func TestBetFieldsValidation(t *testing.T) {
	tcs := []struct {
		desc string
		bet  *types.PlaceBetFields
		err  error
	}{
		{
			desc: "space in UID",
			bet: &types.PlaceBetFields{
				UID:    " ",
				Amount: sdk.NewInt(int64(10)),
				Ticket: "Ticket",
			},
			err: types.ErrInvalidBetUID,
		},
		{
			desc: "invalid UID",
			bet: &types.PlaceBetFields{
				UID:    "invalidUID",
				Amount: sdk.NewInt(int64(10)),
				Ticket: "Ticket",
			},
			err: types.ErrInvalidBetUID,
		},
		{
			desc: "invalid amount",
			bet: &types.PlaceBetFields{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdk.NewInt(int64(-1)),
				Ticket: "Ticket",
			},
			err: types.ErrInvalidAmount,
		},
		{
			desc: "empty amount",
			bet: &types.PlaceBetFields{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Ticket: "Ticket",
			},
			err: types.ErrInvalidAmount,
		},
		{
			desc: "space in ticket",
			bet: &types.PlaceBetFields{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdk.NewInt(int64(10)),
				Ticket: " ",
			},
			err: types.ErrInvalidTicket,
		},
		{
			desc: "valid message",
			bet: &types.PlaceBetFields{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdk.NewInt(int64(10)),
				Ticket: "Ticket",
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := types.BetFieldsValidation(tc.bet)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
