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

func TestTicketFieldsValidation(t *testing.T) {
	tcs := []struct {
		desc     string
		betOdds  *types.BetOdds
		kyc      types.KycDataPayload
		oddsType types.OddsType
		err      error
	}{
		{
			desc: "space in odds UID",
			betOdds: &types.BetOdds{
				SportEventUID:     "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               " ",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: types.KycDataPayload{
				Ignore:   true,
				Approved: true,
				ID:       testAddress,
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err:      types.ErrInvalidOddsUID,
		},
		{
			desc: "space in sport-event UID",
			betOdds: &types.BetOdds{
				SportEventUID:     " ",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: types.KycDataPayload{
				Approved: true,
				ID:       testAddress,
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err:      types.ErrInvalidSportEventUID,
		},
		{
			desc: "empty odds value",
			betOdds: &types.BetOdds{
				SportEventUID:     "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: types.KycDataPayload{
				Approved: true,
				ID:       testAddress,
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err:      types.ErrEmptyOddsValue,
		},
		{
			desc: "no kyc",
			betOdds: &types.BetOdds{
				SportEventUID:     "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
      oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err: types.ErrUserKycFailed,
		},
		{
			desc: "no kyc ID field",
			betOdds: &types.BetOdds{
				SportEventUID:     "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: types.KycDataPayload{
				Approved: true,
				ID:       "",
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err:      types.ErrUserKycFailed,
		},
		{
			desc: "valid message",
			betOdds: &types.BetOdds{
				SportEventUID:     "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: types.KycDataPayload{
				Approved: true,
				ID:       testAddress,
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			p := types.BetPlacementTicketPayload{
				SelectedOdds: tc.betOdds,
				KycData:      tc.kyc,
				OddsType:     tc.oddsType,
			}
			err := p.Validate(testAddress)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
