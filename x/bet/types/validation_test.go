package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestBetFieldsValidation(t *testing.T) {
	tcs := []struct {
		desc string
		bet  *BetPlaceFields
		err  error
	}{
		{
			desc: "space in UID",
			bet: &BetPlaceFields{
				UID:    " ",
				Amount: sdk.NewInt(int64(10)),
				Ticket: "Ticket",
			},
			err: ErrInvalidBetUID,
		},
		{
			desc: "invalid UID",
			bet: &BetPlaceFields{
				UID:    "invalidUID",
				Amount: sdk.NewInt(int64(10)),
				Ticket: "Ticket",
			},
			err: ErrInvalidBetUID,
		},
		{
			desc: "invalid amount",
			bet: &BetPlaceFields{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdk.NewInt(int64(-1)),
				Ticket: "Ticket",
			},
			err: ErrInvalidAmount,
		},
		{
			desc: "empty amount",
			bet: &BetPlaceFields{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Ticket: "Ticket",
			},
			err: ErrInvalidAmount,
		},
		{
			desc: "space in ticket",
			bet: &BetPlaceFields{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdk.NewInt(int64(10)),
				Ticket: " ",
			},
			err: ErrInvalidTicket,
		},
		{
			desc: "valid message",
			bet: &BetPlaceFields{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdk.NewInt(int64(10)),
				Ticket: "Ticket",
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := BetFieldsValidation(tc.bet)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestTicketFieldsValidation(t *testing.T) {
	tcs := []struct {
		desc    string
		betOdds *BetOdds
		err     error
	}{
		{
			desc: "space in odds UID",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           " ",
				Value:         "10",
			},
			err: ErrInvalidOddsUID,
		},
		{
			desc: "space in sport event UID",
			betOdds: &BetOdds{
				SportEventUID: " ",
				UID:           "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:         "10",
			},
			err: ErrInvalidSportEventUID,
		},
		{
			desc: "empty odds value",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           "6e31c60f-2025-48ce-ae79-1dc110f16355",
			},
			err: ErrInvalidOddsValue,
		},
		{
			desc: "invalid odds value",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:         "invalidOdds",
			},
			err: ErrInConvertingOddsToDec,
		},
		{
			desc: "odds value less than 1",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:         "0",
			},
			err: ErrInvalidOddsValue,
		},
		{
			desc: "valid message",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:         "10",
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := TicketFieldsValidation(tc.betOdds)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
