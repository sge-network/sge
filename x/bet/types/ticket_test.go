package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sgetypes "github.com/sge-network/sge/types"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/stretchr/testify/require"
)

func TestTicketFieldsValidation(t *testing.T) {
	tcs := []struct {
		desc     string
		betOdds  *types.BetOdds
		kyc      sgetypes.KycDataPayload
		oddsType types.OddsType
		err      error
	}{
		{
			desc: "space in odds UID",
			betOdds: &types.BetOdds{
				MarketUID:         "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               " ",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: sgetypes.KycDataPayload{
				Ignore:   true,
				Approved: true,
				ID:       testAddress,
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err:      types.ErrInvalidOddsUID,
		},
		{
			desc: "space in market UID",
			betOdds: &types.BetOdds{
				MarketUID:         " ",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: sgetypes.KycDataPayload{
				Approved: true,
				ID:       testAddress,
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err:      types.ErrInvalidMarketUID,
		},
		{
			desc: "empty odds value",
			betOdds: &types.BetOdds{
				MarketUID:         "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: sgetypes.KycDataPayload{
				Approved: true,
				ID:       testAddress,
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err:      types.ErrEmptyOddsValue,
		},
		{
			desc: "no kyc",
			betOdds: &types.BetOdds{
				MarketUID:         "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err:      types.ErrUserKycFailed,
		},
		{
			desc: "no kyc ID field",
			betOdds: &types.BetOdds{
				MarketUID:         "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: sgetypes.KycDataPayload{
				Approved: true,
				ID:       "",
			},
			oddsType: types.OddsType_ODDS_TYPE_DECIMAL,
			err:      types.ErrUserKycFailed,
		},
		{
			desc: "valid message",
			betOdds: &types.BetOdds{
				MarketUID:         "6e31c60f-2025-48ce-ae79-1dc110f16355",
				UID:               "6e31c60f-2025-48ce-ae79-1dc110f16355",
				Value:             "10",
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			kyc: sgetypes.KycDataPayload{
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
