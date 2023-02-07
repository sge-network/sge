package keeper_test

import (
	"testing"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/sportevent/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgServerAddEvent(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgAddSportEventRequest
	}

	u1 := uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.MsgAddSportEventResponse
		wantErr error
	}{
		{
			name: "test the empty or invalid format ticket",
			args: args{
				msg: types.NewMsgAddEvent(sample.AccAddress(), ""),
			},
			want:    nil,
			wantErr: types.ErrInVerification,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := msgk.AddSportEvent(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func TestMsgServerResolveEvent(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgResolveSportEventRequest
	}

	u1 := uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.MsgResolveSportEventResponse
		wantErr error
	}{
		{
			name: "test the empty or invalid format ticket",
			args: args{
				msg: types.NewMsgResolveEvent(sample.AccAddress(), ""),
			},
			want:    nil,
			wantErr: types.ErrInVerification,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := msgk.ResolveSportEvent(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func TestMsgServerAddEventResponse(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	u1 := uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	t.Run("Error in ticket fields validation", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":      "invalid uid",
			"start_ts": uint64(time.Now().Add(time.Minute).Unix()),
			"end_ts":   uint64(time.Now().Add(time.Minute * 5).Unix()),
			"exp":      9999999999,
			"iat":      1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.AddSportEvent(wctx, types.NewMsgAddEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("pre existing uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":      u1,
			"start_ts": uint64(time.Now().Add(time.Minute).Unix()),
			"end_ts":   uint64(time.Now().Add(time.Minute * 5).Unix()),
			"odds":     []types.Odds{{UID: uuid.NewString(), Meta: "odds 1"}, {UID: uuid.NewString(), Meta: "odds 2"}},
			"exp":      9999999999,
			"iat":      1111111111,
			"meta":     "Winner of x:y",
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.AddSportEvent(wctx, types.NewMsgAddEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventAlreadyExist)
		assert.Nil(t, response)
	})
}

func TestMsgServerResolveEventResponse(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

	u1, u2, u3 := uuid.NewString(), uuid.NewString(), uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u3,
		Creator: sample.AccAddress(),
		Status:  types.SportEventStatus_SPORT_EVENT_STATUS_CANCELLED,
	})

	t.Run("Error in ticket fields validation", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":           "invalid uid",
			"status":        types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED,
			"resolution_ts": uint64(time.Now().Unix()),
			"exp":           9999999999,
			"iat":           1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveSportEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("non existing event", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u2,
			"status":           types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveSportEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventNotFound)
		assert.Nil(t, response)
	})

	t.Run("non pending event resolution", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u3,
			"status":           types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveSportEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventIsNotPending)
		assert.Nil(t, response)
	})

	t.Run("invalid winner odds uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u1,
			"status":           types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{"invalidWId"},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveSportEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("invalid winner odds uid, not contained in the parent", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u1,
			"status":           types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveSportEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrInvalidWinnerOdd)
		assert.Nil(t, response)
	})
}

func TestMsgServerUpdateEvent(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgUpdateSportEventRequest
	}

	u1 := uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.MsgUpdateSportEventResponse
		wantErr error
	}{
		{
			name: "test the empty or invalid format ticket",
			args: args{
				msg: types.NewMsgUpdateEvent(sample.AccAddress(), ""),
			},
			want:    nil,
			wantErr: types.ErrInVerification,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := msgk.UpdateSportEvent(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func TestMsgServerUpdateEventResponse(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

	u1, u2, u3 := uuid.NewString(), uuid.NewString(), uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u2,
		Creator: sample.AccAddress(),
		Status:  types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED,
	})

	t.Run("invalid SportEvent id", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"UID": u3,
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.UpdateSportEvent(wctx, types.NewMsgUpdateEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventNotFound)
		assert.Nil(t, response)
	})
	t.Run("updating an declared event", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":      u2,
			"start_ts": uint64(time.Now().UnixNano()),
			"end_ts":   uint64(time.Now().Add(time.Hour).UnixNano()),
			"exp":      9999999999,
			"iat":      1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.UpdateSportEvent(wctx, types.NewMsgUpdateEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrCanNotBeAltered)
		assert.Nil(t, response)
	})
}
