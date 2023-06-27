package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/market/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestAddMarketTicketPayloadValidation(t *testing.T) {
	param := types.DefaultParams()
	_, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(time.Now())

	sampleUID := uuid.NewString()

	tests := []struct {
		name    string
		payload types.MarketAddTicketPayload
		err     error
	}{
		{
			name: "valid",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "odds 1"},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status:       types.MarketStatus_MARKET_STATUS_ACTIVE,
				MinBetAmount: sdk.NewInt(1000000),
				BetFee:       sdk.NewInt(10),
				Meta:         "sample market",
			},
		},
		{
			name: "invalid end time",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid status",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid uuid",
			payload: types.MarketAddTicketPayload{
				UID:     "invalid uuid",
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid odds count",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "odds 1"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty meta",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "odds 1"},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
				Meta:   "",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "large meta",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "odds 1"},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
				Meta: `Winner of x:y is the final winner of the game,
				it is obvious the winner is not the champion yet but if it happens,
				the winning users will reward 1M dollars each plus a furnished villa in the Beverley hills as a gift.
				attention! this detail will not be stored in the chain because it's definitely a scam.`,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid odds meta",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: ""},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status:       types.MarketStatus_MARKET_STATUS_ACTIVE,
				MinBetAmount: sdk.NewInt(1000000),
				BetFee:       sdk.NewInt(10),
				Meta:         "sample market",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid odds long meta",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: `Winner odds of x:y is the final winner of the game,
					it is obvious the winner is not the champion yet but if it happens,
					the winning users will reward 1M dollars each plus a furnished villa in the Beverley hills as a gift.
					attention! this detail will not be stored in the chain because it's definitely a scam.`},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status:       types.MarketStatus_MARKET_STATUS_ACTIVE,
				MinBetAmount: sdk.NewInt(1000000),
				BetFee:       sdk.NewInt(10),
				Meta:         "sample market",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid odds uuid",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: "invalid uuid", Meta: "odds 1"},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status:       types.MarketStatus_MARKET_STATUS_ACTIVE,
				MinBetAmount: sdk.NewInt(1000000),
				BetFee:       sdk.NewInt(10),
				Meta:         "sample market",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "duplicate odds uuid",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: sampleUID, Meta: "odds 1"},
					{UID: sampleUID, Meta: "odds 2"},
				},
				Status:       types.MarketStatus_MARKET_STATUS_ACTIVE,
				MinBetAmount: sdk.NewInt(1000000),
				BetFee:       sdk.NewInt(10),
				Meta:         "sample market",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid min bet amount",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				Creator: sample.AccAddress(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "odds 1"},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status:       types.MarketStatus_MARKET_STATUS_ACTIVE,
				MinBetAmount: sdk.NewInt(100000),
				BetFee:       sdk.NewInt(10),
				Meta:         "sample market",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate(ctx, &param)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdateMarketTicketPayloadValidation(t *testing.T) {
	param := types.DefaultParams()
	_, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(time.Now())

	tests := []struct {
		name    string
		payload types.MarketUpdateTicketPayload
		err     error
	}{
		{
			name: "valid",
			payload: types.MarketUpdateTicketPayload{
				UID:          uuid.NewString(),
				StartTS:      cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:        cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Status:       types.MarketStatus_MARKET_STATUS_ACTIVE,
				MinBetAmount: sdk.NewInt(1000000),
				BetFee:       sdk.NewInt(10),
			},
		},
		{
			name: "invalid end time",
			payload: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid status",
			payload: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid min bet amount",
			payload: types.MarketUpdateTicketPayload{
				UID:          uuid.NewString(),
				StartTS:      cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:        cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Status:       types.MarketStatus_MARKET_STATUS_ACTIVE,
				MinBetAmount: sdk.NewInt(100000),
				BetFee:       sdk.NewInt(10),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate(ctx, &param)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestResolveMarketTicketPayloadValidation(t *testing.T) {
	_, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	ctx = ctx.WithBlockTime(time.Now())

	tests := []struct {
		name    string
		payload types.MarketResolutionTicketPayload
		err     error
	}{
		{
			name: "valid",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
		},
		{
			name: "invaid status",
			payload: types.MarketResolutionTicketPayload{
				UID:          uuid.NewString(),
				ResolutionTS: cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				Status:       types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "winner set when result not declared",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "excessive winner odds",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{uuid.NewString(), uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "no resolution time set",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid uuid",
			payload: types.MarketResolutionTicketPayload{
				UID:            "invalid uuid",
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "result declared no winner odds",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid odds uuid",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{"invalid uuid"},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
