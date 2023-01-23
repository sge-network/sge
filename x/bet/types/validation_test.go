package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestBetFieldsValidation(t *testing.T) {
	tcs := []struct {
		desc string
		bet  *PlaceBetFields
		err  error
	}{
		{
			desc: "space in UID",
			bet: &PlaceBetFields{
				UID:      " ",
				Amount:   sdk.NewInt(int64(10)),
				Ticket:   "Ticket",
				OddsType: 1,
			},
			err: ErrInvalidBetUID,
		},
		{
			desc: "invalid UID",
			bet: &PlaceBetFields{
				UID:      "invalidUID",
				Amount:   sdk.NewInt(int64(10)),
				Ticket:   "Ticket",
				OddsType: 1,
			},
			err: ErrInvalidBetUID,
		},
		{
			desc: "invalid amount",
			bet: &PlaceBetFields{
				UID:      "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount:   sdk.NewInt(int64(-1)),
				Ticket:   "Ticket",
				OddsType: 1,
			},
			err: ErrInvalidAmount,
		},
		{
			desc: "empty amount",
			bet: &PlaceBetFields{
				UID:      "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Ticket:   "Ticket",
				OddsType: 1,
			},
			err: ErrInvalidAmount,
		},
		{
			desc: "space in ticket",
			bet: &PlaceBetFields{
				UID:      "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount:   sdk.NewInt(int64(10)),
				Ticket:   " ",
				OddsType: 1,
			},
			err: ErrInvalidTicket,
		},
		{
			desc: "space in ticket",
			bet: &PlaceBetFields{
				UID:    "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount: sdk.NewInt(int64(10)),
				Ticket: " ",
			},
			err: ErrInvalidOddsType,
		},
		{
			desc: "valid message",
			bet: &PlaceBetFields{
				UID:      "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Amount:   sdk.NewInt(int64(10)),
				Ticket:   "Ticket",
				OddsType: 1,
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
		kyc     *KycDataPayload
		err     error
	}{
		{
			desc: "space in odds UID",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           " ",
				Value:         "10",
			},
			kyc: &KycDataPayload{
				KycRequired: true,
				KycApproved: true,
				KycId:       "creator",
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
			kyc: &KycDataPayload{
				KycRequired: true,
				KycApproved: true,
				KycId:       "creator",
			},
			err: ErrInvalidSportEventUID,
		},
		{
			desc: "empty odds value",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:         "",
			},
			kyc: &KycDataPayload{
				KycRequired: true,
				KycApproved: true,
				KycId:       "creator",
			},
			err: ErrEmptyOddsValue,
		},
		{
			desc: "no kyc",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:         "10",
			},
			err: ErrNoKycField,
		},
		{
			desc: "no kyc ID field",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:         "10",
			},
			kyc: &KycDataPayload{
				KycRequired: true,
				KycApproved: true,
				KycId:       "",
			},
			err: ErrNoKycIdField,
		},
		{
			desc: "valid message",
			betOdds: &BetOdds{
				SportEventUID: "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:           "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:         "10",
			},
			kyc: &KycDataPayload{
				KycRequired: true,
				KycApproved: true,
				KycId:       "creator",
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {

			err := TicketFieldsValidation(&BetPlacementTicketPayload{
				SelectedOdds: tc.betOdds,
				KycData:      tc.kyc,
			})
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
