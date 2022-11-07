package keeper_test

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/sportevent/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_msgServer_AddEvent(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgAddEvent
	}

	u1 := uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.SportResponse
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
			got, err := msgk.AddEvent(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func Test_msgServer_ResolveEvent(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgResolveEvent
	}

	u1 := uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.SportResponse
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
			got, err := msgk.ResolveEvent(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func Test_msgServer_AddEventResponse(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

	u1, _ := uuid.NewString(), uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	t.Run("Error in ticket fields validation", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid": "invalid uid",
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		_, err = msgk.AddEvent(wctx, types.NewMsgAddEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
	})

	t.Run("pre existing uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":       u1,
			"start_ts":  uint64(time.Now().Add(time.Minute).Unix()),
			"end_ts":    uint64(time.Now().Add(time.Minute * 5).Unix()),
			"odds_uids": []string{uuid.NewString(), uuid.NewString()},
			"exp":       9999999999,
			"iat":       1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.AddEvent(wctx, types.NewMsgAddEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventAlreadyExist)
		assert.Nil(t, response)
	})
}

func Test_msgServer_ResolveEventResponse(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

	u1, u2, u3 := uuid.NewString(), uuid.NewString(), uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u3,
		Creator: sample.AccAddress(),
		Status:  types.SportEventStatus_STATUS_CANCELLED,
	})

	t.Run("Error in ticket fields validation", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid": "invalid uid",
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("non existing event", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u2,
			"status":           types.SportEventStatus_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().UnixNano()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventNotFound)
		assert.Nil(t, response)
	})

	t.Run("non pending event resolution", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u3,
			"status":           types.SportEventStatus_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().UnixNano()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventIsNotPending)
		assert.Nil(t, response)
	})

	t.Run("invalid winner odds uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u1,
			"status":           types.SportEventStatus_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().UnixNano()),
			"winner_odds_uids": []string{"invalid"},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("invalid winner odds uid, not contained in the parent", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u1,
			"status":           types.SportEventStatus_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().UnixNano()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrInvalidWinnerOdd)
		assert.Nil(t, response)
	})
}

func Test_msgServer_UpdateEvent(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgUpdateEvent
	}

	u1 := uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.SportResponse
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
			got, err := msgk.UpdateEvent(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func Test_msgServer_UpdateEventResponse(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

	u1, u2, u3 := uuid.NewString(), uuid.NewString(), uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u2,
		Creator: sample.AccAddress(),
		Status:  types.SportEventStatus_STATUS_RESULT_DECLARED,
	})

	t.Run("invalid SportEvent id", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid": u3,
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.UpdateEvent(wctx, types.NewMsgUpdateEvent(sample.AccAddress(), validEmptyTicket))
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

		response, err := msgk.UpdateEvent(wctx, types.NewMsgUpdateEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrCanNotBeAltered)
		assert.Nil(t, response)
	})

}
