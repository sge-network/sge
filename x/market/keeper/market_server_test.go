package keeper_test

import (
	"testing"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgServerAddEvent(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgAddMarket
	}

	u1 := uuid.NewString()
	k.SetMarket(ctx, types.Market{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.MsgAddMarketResponse
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
			got, err := msgk.AddMarket(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func TestMsgServerResolveEvent(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgResolveMarket
	}

	u1 := uuid.NewString()
	k.SetMarket(ctx, types.Market{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.MsgResolveMarketResponse
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
			got, err := msgk.ResolveMarket(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func TestMsgServerAddEventResponse(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	u1 := uuid.NewString()
	k.SetMarket(ctx, types.Market{
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

		response, err := msgk.AddMarket(wctx, types.NewMsgAddEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("pre existing uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":                       u1,
			"start_ts":                  uint64(time.Now().Add(time.Minute).Unix()),
			"end_ts":                    uint64(time.Now().Add(time.Minute * 5).Unix()),
			"odds":                      []types.Odds{{UID: uuid.NewString(), Meta: "odds 1"}, {UID: uuid.NewString(), Meta: "odds 2"}},
			"exp":                       9999999999,
			"iat":                       1111111111,
			"meta":                      "Winner of x:y",
			"sr_contribution_for_house": "2",
			"status":                    types.MarketStatus_MARKET_STATUS_ACTIVE,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.AddMarket(wctx, types.NewMsgAddEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventAlreadyExist)
		assert.Nil(t, response)
	})
}

func TestMsgServerResolveEventResponse(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

	u1, u2, u3 := uuid.NewString(), uuid.NewString(), uuid.NewString()
	k.SetMarket(ctx, types.Market{
		UID:     u1,
		Creator: sample.AccAddress(),
		Status:  types.MarketStatus_MARKET_STATUS_ACTIVE,
	})
	k.SetMarket(ctx, types.Market{
		UID:     u3,
		Creator: sample.AccAddress(),
		Status:  types.MarketStatus_MARKET_STATUS_CANCELED,
	})

	t.Run("invalid uid for the market", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":           "invalid uid",
			"status":        types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			"resolution_ts": uint64(time.Now().Unix()),
			"exp":           9999999999,
			"iat":           1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("non existing market", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u2,
			"status":           types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventNotFound)
		assert.Nil(t, response)
	})

	t.Run("non active market resolution", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u3,
			"status":           types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventIsNotActiveOrInactive)
		assert.Nil(t, response)
	})

	t.Run("invalid winner odds uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u1,
			"status":           types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{"invalidWId"},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("lots of winner odds uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u1,
			"status":           types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{uuid.NewString(), uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("canceled or aborted with winner id", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u1,
			"status":           types.MarketStatus_MARKET_STATUS_CANCELED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		assert.Nil(t, response)
	})

	t.Run("invalid winner odds uid, not contained in the parent", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":              u1,
			"status":           types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			"resolution_ts":    uint64(time.Now().Unix()),
			"winner_odds_uids": []string{uuid.NewString()},
			"exp":              9999999999,
			"iat":              1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrInvalidWinnerOdds)
		assert.Nil(t, response)
	})
}

func TestMsgServerUpdateEvent(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgUpdateMarket
	}

	u1 := uuid.NewString()
	k.SetMarket(ctx, types.Market{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.MsgUpdateMarketResponse
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
			got, err := msgk.UpdateMarket(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func TestMsgServerUpdateEventResponse(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)

	u1, u2, u3 := uuid.NewString(), uuid.NewString(), uuid.NewString()
	k.SetMarket(ctx, types.Market{
		UID:     u1,
		Creator: sample.AccAddress(),
	})
	k.SetMarket(ctx, types.Market{
		UID:     u2,
		Creator: sample.AccAddress(),
		Status:  types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
	})

	t.Run("invalid Market id", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"UID": u3,
			"exp": 9999999999,
			"iat": 1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.UpdateMarket(wctx, types.NewMsgUpdateEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrEventNotFound)
		assert.Nil(t, response)
	})
	t.Run("updating an declared market", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":      u2,
			"start_ts": uint64(time.Now().UnixNano()),
			"end_ts":   uint64(time.Now().Add(time.Hour).UnixNano()),
			"exp":      9999999999,
			"iat":      1111111111,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.UpdateMarket(wctx, types.NewMsgUpdateEvent(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrCanNotBeAltered)
		assert.Nil(t, response)
	})
}
