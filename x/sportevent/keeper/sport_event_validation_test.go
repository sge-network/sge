package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/sportevent/types"
	"github.com/stretchr/testify/require"
)

func TestValidateCreationEvent(t *testing.T) {
	k, _, wctx, _ := setupMsgServerAndKeeper(t)
	t1 := time.Now()
	params := k.GetParams(wctx)

	negativeBetAmount := sdk.NewInt(-5)
	lowerBetAmount := params.EventMinBetAmount.Sub(sdk.NewInt(5))

	tests := []struct {
		name string
		msg  types.SportEventAddTicketPayload
		err  error
	}{
		{
			name: "valid request",
			msg: types.SportEventAddTicketPayload{
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
				Status:                 types.SportEventStatus_SPORT_EVENT_STATUS_PENDING,
			},
		},
		{
			name: "same timestamp",
			msg: types.SportEventAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "end timestamp before current timestamp",
			msg: types.SportEventAddTicketPayload{
				Creator: sample.AccAddress(),
				EndTS:   uint64(t1.Add(-time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid uid",
			msg: types.SportEventAddTicketPayload{
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
			msg: types.SportEventAddTicketPayload{
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
			msg: types.SportEventAddTicketPayload{
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
			msg: types.SportEventAddTicketPayload{
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
			msg: types.SportEventAddTicketPayload{
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
			msg: types.SportEventAddTicketPayload{
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
			msg: types.SportEventAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
				MinBetAmount:           params.EventMinBetAmount,
				BetFee:                 params.EventMinBetFee,
				Meta:                   "Winner of x:y",
				SrContributionForHouse: sdk.NewInt(2),
				Status:                 types.SportEventStatus_SPORT_EVENT_STATUS_PENDING,
			},
		},
		{
			name: "large metadata",
			msg: types.SportEventAddTicketPayload{
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
		msg  types.SportEventResolutionTicketPayload
		err  error
	}{
		{
			name: "valid request",
			msg: types.SportEventResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         4,
			},
		},
		{
			name: "invalid resolution ts",
			msg: types.SportEventResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   0,
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         4,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid uid",
			msg: types.SportEventResolutionTicketPayload{
				UID:            "invalid uid",
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         4,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty winner odds",
			msg: types.SportEventResolutionTicketPayload{
				UID:          uuid.NewString(),
				ResolutionTS: uint64(t1.Unix()),
				Status:       4,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid winner odds",
			msg: types.SportEventResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{"invalid winner odds"},
				Status:         4,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "msg status pending",
			msg: types.SportEventResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         0,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "msg invalid enum status",
			msg: types.SportEventResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         5,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "msg invalid enum status, pending",
			msg: types.SportEventResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   uint64(t1.Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         1,
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
	lowerBetAmount := params.EventMinBetAmount.Sub(sdk.NewInt(5))

	sportEvent := types.SportEvent{
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
		msg  types.SportEventUpdateTicketPayload
		err  error
	}{
		{
			name: "valid request",
			msg: types.SportEventUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute * 2).Unix()),
			},
		},
		{
			name: "same timestamp",
			msg: types.SportEventUpdateTicketPayload{
				StartTS: uint64(t1.Add(time.Minute).Unix()),
				EndTS:   uint64(t1.Add(time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "end timestamp before current timestamp",
			msg: types.SportEventUpdateTicketPayload{
				EndTS: uint64(t1.Add(-time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid min amount, negative",
			msg: types.SportEventUpdateTicketPayload{
				UID:          uuid.NewString(),
				StartTS:      uint64(t1.Add(time.Minute).Unix()),
				EndTS:        uint64(t1.Add(time.Minute * 2).Unix()),
				MinBetAmount: negativeBetAmount,
				BetFee:       params.EventMinBetFee,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid min amount, less than required",
			msg: types.SportEventUpdateTicketPayload{
				UID:          uuid.NewString(),
				StartTS:      uint64(t1.Add(time.Minute).Unix()),
				EndTS:        uint64(t1.Add(time.Minute * 2).Unix()),
				MinBetAmount: lowerBetAmount,
				BetFee:       params.EventMinBetFee,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := k.ValidateEventUpdate(wctx, tt.msg, sportEvent)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
