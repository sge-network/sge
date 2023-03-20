package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/require"
)

func TestValidateCreationEvent(t *testing.T) {
	k, _, wctx, _ := setupMsgServerAndKeeper(t)
	t1 := time.Now()
	params := k.GetParams(wctx)

	negativeBetAmount := sdk.NewInt(-5)
	lowerBetAmount := params.MinBetAmount.Sub(sdk.NewInt(5))

	tests := []struct {
		name string
		msg  types.MarketAddTicketPayload
		err  error
	}{
		{
			name: "valid request",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
				Meta:                   "Winner of x:y",
				SrContributionForHouse: sdk.NewInt(2),
				Status:                 types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
		},
		{
			name: "same timestamp",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "end timestamp before current timestamp",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				EndTS:   uint64(t1.Add(-time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid uid",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     "invalid uuid",
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "few odds than required",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid odds id",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: "invalid id", Meta: "invalid odds"},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "duplicate odds id",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: "8779cf93-925c-4818-bc81-13c359e0deb8", Meta: "Odds 1"},
					{UID: "8779cf93-925c-4818-bc81-13c359e0deb8", Meta: "invalid odds"},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid min amount, negative",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
				MinBetAmount: negativeBetAmount,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid min amount, less than required",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
				MinBetAmount: lowerBetAmount,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid request, with bet constraint",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
				MinBetAmount:           params.MinBetAmount,
				BetFee:                 params.MinBetFee,
				Meta:                   "Winner of x:y",
				SrContributionForHouse: sdk.NewInt(2),
				Status:                 types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
		},
		{
			name: "large metadata",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
				Meta: `Winner of x:y is the final winner of the game,
				it is obvious the winner is not the champion yet but if it happens,
				the winning users will reward 1M dollars each plus a furnished villa in the Beverley hills as a gift.
				attention! this detail will not be stored in the chain because it's definitely a scam.`,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := k.ValidateEventAdd(wctx, tt.msg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidateResolveEvent(t *testing.T) {
	k, _, _, _ := setupMsgServerAndKeeper(t)
	t1 := time.Now()

	tests := []struct {
		name string
		msg  types.MarketResolutionTicketPayload
		err  error
	}{
		{
			name: "valid request",
			msg: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
		},
		{
			name: "invalid resolution ts",
			msg: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   0,
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid uid",
			msg: types.MarketResolutionTicketPayload{
				UID:            "invalid uid",
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty winner odds",
			msg: types.MarketResolutionTicketPayload{
				UID:          uuid.NewString(),
				ResolutionTS: uint64(t1.Unix()),
				Status:       types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid winner odds",
			msg: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{"invalid winner odds"},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "msg status inactive",
			msg: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_INACTIVE,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "msg invalid enum status",
			msg: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         6,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "msg invalid enum status, active",
			msg: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := k.ValidateEventResolution(tt.msg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdateEventValidation(t *testing.T) {
	k, _, wctx, _ := setupMsgServerAndKeeper(t)
	params := k.GetParams(wctx)

	t1 := time.Now()

	negativeBetAmount := sdk.NewInt(-5)
	lowerBetAmount := params.MinBetAmount.Sub(sdk.NewInt(5))

	market := types.Market{
		Creator: sample.AccAddress(),
		StartTS: uint64(t1.Add(time.Minute).Unix()),
		EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
		UID:     uuid.NewString(),
		Odds: []*types.Odds{
			{UID: uuid.NewString(), Meta: "Odds 1"},
			{UID: uuid.NewString(), Meta: "Odds 2"},
		},
		Meta: "Winner of x:y",
	}

	tests := []struct {
		name string
		msg  types.MarketUpdateTicketPayload
		err  error
	}{
		{
			name: "valid request active",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
		},
		{
			name: "valid request inactive",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_INACTIVE,
			},
		},
		{
			name: "invalid status, declared",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid status, canceled",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_CANCELED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid status, aborted",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_ABORTED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid status, unpecified",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "same timestamp",
			msg: types.MarketUpdateTicketPayload{
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "end timestamp before current timestamp",
			msg: types.MarketUpdateTicketPayload{
				EndTS: uint64(t1.Add(-time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid min amount, negative",
			msg: types.MarketUpdateTicketPayload{
				UID:          uuid.NewString(),
				StartTS:      uint64(t1.Add(time.Minute).Unix()),
				EndTS:        uint64(t1.Add(time.Minute * 2).Unix()),
				MinBetAmount: negativeBetAmount,
				BetFee:       params.MinBetFee,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid min amount, less than required",
			msg: types.MarketUpdateTicketPayload{
				UID:          uuid.NewString(),
				StartTS:      uint64(t1.Add(time.Minute).Unix()),
				EndTS:        uint64(t1.Add(time.Minute * 2).Unix()),
				MinBetAmount: lowerBetAmount,
				BetFee:       params.MinBetFee,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := k.ValidateEventUpdate(wctx, tt.msg, market)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
