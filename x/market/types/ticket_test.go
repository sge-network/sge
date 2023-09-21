package types_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"

	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/market/types"
)

func TestAddMarketTicketPayloadValidation(t *testing.T) {
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
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "odds 1"},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
				Meta:   "sample market",
			},
		},
		{
			name: "invalid end time",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Unix()),
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid status",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid uuid",
			payload: types.MarketAddTicketPayload{
				UID:     "invalid uuid",
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid odds count",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "odds 1"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "empty meta",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "odds 1"},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
				Meta:   "",
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "large meta",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "odds 1"},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
				Meta:   simappUtil.RandomString(types.MaxAllowedCharactersForMeta + 1),
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid odds meta",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: ""},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
				Meta:   "sample market",
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid odds long meta",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: simappUtil.RandomString(types.MaxAllowedCharactersForMeta + 1)},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
				Meta:   "sample market",
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid odds uuid",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: "invalid uuid", Meta: "odds 1"},
					{UID: uuid.NewString(), Meta: "odds 2"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
				Meta:   "sample market",
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "duplicate odds uuid",
			payload: types.MarketAddTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Odds: []*types.Odds{
					{UID: sampleUID, Meta: "odds 1"},
					{UID: sampleUID, Meta: "odds 2"},
				},
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
				Meta:   "sample market",
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate(ctx)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdateMarketTicketPayloadValidation(t *testing.T) {
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
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
		},
		{
			name: "invalid end time",
			payload: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Unix()),
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid status",
			payload: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: cast.ToUint64(ctx.BlockTime().Unix()),
				EndTS:   cast.ToUint64(ctx.BlockTime().Add(5 * time.Minute).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate(ctx)
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
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "winner set when result not declared",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "excessive winner odds",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{uuid.NewString(), uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "no resolution time set",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid uuid",
			payload: types.MarketResolutionTicketPayload{
				UID:            "invalid uuid",
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{uuid.NewString()},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "result declared no winner odds",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrtypes.ErrInvalidRequest,
		},
		{
			name: "invalid odds uuid",
			payload: types.MarketResolutionTicketPayload{
				UID:            uuid.NewString(),
				ResolutionTS:   cast.ToUint64(ctx.BlockTime().Add(10 * time.Minute).Unix()),
				WinnerOddsUIDs: []string{"invalid uuid"},
				Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrtypes.ErrInvalidRequest,
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
