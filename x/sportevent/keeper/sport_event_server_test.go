package keeper_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/sportevent/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		want    *types.MsgSportResponse
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
		{
			name: "test the valid case",
			args: args{
				msg: types.NewMsgAddEvent(sample.AccAddress(), func() string {
					validEmptyTicketClaims := jwt.MapClaims{
						"events": []interface{}{},
						"exp":    9999999999,
						"iat":    1111111111,
					}
					validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
					require.NoError(t, err)
					return validEmptyTicket
				}()),
			},
			want:    &types.MsgSportResponse{FailedEvents: []*types.FailedEvent{}, SuccessEvents: []string{}},
			wantErr: nil,
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
		want    *types.MsgSportResponse
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
		{
			name: "test the valid case",
			args: args{
				msg: types.NewMsgResolveEvent(sample.AccAddress(), func() string {
					validEmptyTicketClaims := jwt.MapClaims{
						"events": []interface{}{},
						"exp":    9999999999,
						"iat":    1111111111,
					}
					validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
					require.NoError(t, err)
					return validEmptyTicket
				}()),
			},
			want:    &types.MsgSportResponse{FailedEvents: []*types.FailedEvent{}, SuccessEvents: []string{}},
			wantErr: nil,
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

	u1, u2 := uuid.NewString(), uuid.NewString()
	k.SetSportEvent(ctx, types.SportEvent{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	t.Run("Error in ticket fields validation", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"events": []types.SportEvent{{
				UID: "invalid uid",
			}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.AddEvent(wctx, types.NewMsgAddEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, "invalid uid")
	})

	t.Run("pre existing uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"events": []types.SportEvent{{
				UID:      u1,
				StartTS:  uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:    uint64(time.Now().Add(time.Minute * 5).Unix()),
				OddsUIDs: []string{uuid.NewString(), uuid.NewString()},
			}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.AddEvent(wctx, types.NewMsgAddEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, u1)
	})
	t.Run("duplicate uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"events": []types.SportEvent{
				{
					UID:      u2,
					StartTS:  uint64(time.Now().Add(time.Minute).Unix()),
					EndTS:    uint64(time.Now().Add(time.Minute * 5).Unix()),
					OddsUIDs: []string{uuid.NewString(), uuid.NewString()},
				},
				{
					UID:      u2,
					StartTS:  uint64(time.Now().Add(time.Minute).Unix()),
					EndTS:    uint64(time.Now().Add(time.Minute * 5).Unix()),
					OddsUIDs: []string{uuid.NewString(), uuid.NewString()},
				}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.AddEvent(wctx, types.NewMsgAddEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, u2)
		assert.Len(t, response.SuccessEvents, 1)
		assert.EqualValues(t, response.SuccessEvents[0], u2)
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
			"events": []types.SportEvent{{
				UID: "invalid uid",
			}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, "invalid uid")
	})

	t.Run("non existing event", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"events": []types.SportEvent{{
				UID:          u2,
				Status:       types.SportEventStatus_STATUS_RESULT_DECLARED,
				ResolutionTS: uint64(time.Now().UnixNano()),
				WinnerOddsUIDs: map[string][]byte{
					uuid.NewString(): nil,
				},
			}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, u2)
	})

	t.Run("non pending event resolution", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"events": []types.SportEvent{{
				UID:          u3,
				Status:       types.SportEventStatus_STATUS_RESULT_DECLARED,
				ResolutionTS: uint64(time.Now().UnixNano()),
				WinnerOddsUIDs: map[string][]byte{
					uuid.NewString(): nil,
				},
			}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, u3)
	})

	t.Run("invalid winner odds uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"events": []types.SportEvent{{
				UID:          u1,
				Status:       types.SportEventStatus_STATUS_RESULT_DECLARED,
				ResolutionTS: uint64(time.Now().UnixNano()),
				WinnerOddsUIDs: map[string][]byte{
					"invalidWId": nil,
				},
			}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, u1)
	})

	t.Run("invalid winner odds uid, not contained in the parent", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"events": []types.SportEvent{{
				UID:          u1,
				Status:       types.SportEventStatus_STATUS_RESULT_DECLARED,
				ResolutionTS: uint64(time.Now().UnixNano()),
				WinnerOddsUIDs: map[string][]byte{
					uuid.NewString(): nil,
				},
			}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveEvent(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, u1)
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
		want    *types.MsgSportResponse
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
		{
			name: "test the valid case",
			args: args{
				msg: types.NewMsgUpdateEvent(sample.AccAddress(), func() string {
					validEmptyTicketClaims := jwt.MapClaims{
						"events": []interface{}{},
						"exp":    9999999999,
						"iat":    1111111111,
					}
					validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
					require.NoError(t, err)
					return validEmptyTicket
				}()),
			},
			want:    &types.MsgSportResponse{FailedEvents: []*types.FailedEvent{}, SuccessEvents: []string{}},
			wantErr: nil,
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
			"events": []types.SportEvent{{
				UID: u3,
			}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.UpdateEvent(wctx, types.NewMsgUpdateEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, u3)
	})
	t.Run("updating an declared event", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"events": []types.SportEvent{{
				UID:     u2,
				StartTS: uint64(time.Now().UnixNano()),
				EndTS:   uint64(time.Now().Add(time.Hour).UnixNano()),
			}},
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.UpdateEvent(wctx, types.NewMsgUpdateEvent(sample.AccAddress(), validEmptyTicket))
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.FailedEvents, 1)
		assert.EqualValues(t, response.FailedEvents[0].ID, u2)
	})

}
