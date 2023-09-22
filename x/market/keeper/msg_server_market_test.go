package keeper_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/market/types"
)

func TestMsgServerAdd(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgAdd
	}

	u1 := uuid.NewString()
	k.SetMarket(ctx, types.Market{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.MsgAddResponse
		wantErr error
	}{
		{
			name: "test the empty or invalid format ticket",
			args: args{
				msg: types.NewMsgAdd(sample.AccAddress(), ""),
			},
			want:    nil,
			wantErr: types.ErrInTicketVerification,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := msgk.Add(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func TestMsgServerAddMarketResponse(t *testing.T) {
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

		response, err := msgk.Add(
			wctx,
			types.NewMsgAdd(sample.AccAddress(), validEmptyTicket),
		)
		assert.ErrorIs(t, err, types.ErrInTicketPayloadValidation)
		assert.Nil(t, response)
	})

	t.Run("pre existing uid", func(t *testing.T) {
		validEmptyTicketClaims := jwt.MapClaims{
			"uid":      u1,
			"start_ts": uint64(time.Now().Add(time.Minute).Unix()),
			"end_ts":   uint64(time.Now().Add(time.Minute * 5).Unix()),
			"odds": []types.Odds{
				{UID: uuid.NewString(), Meta: "odds 1"},
				{UID: uuid.NewString(), Meta: "odds 2"},
			},
			"exp":    9999999999,
			"iat":    1111111111,
			"meta":   "Winner of x:y",
			"status": types.MarketStatus_MARKET_STATUS_ACTIVE,
		}
		validEmptyTicket, err := createJwtTicket(validEmptyTicketClaims)
		require.NoError(t, err)

		response, err := msgk.Add(
			wctx,
			types.NewMsgAdd(sample.AccAddress(), validEmptyTicket),
		)
		assert.ErrorIs(t, err, types.ErrMarketAlreadyExist)
		assert.Nil(t, response)
	})
}

func TestMsgServerUpdate(t *testing.T) {
	k, msgk, ctx, wctx := setupMsgServerAndKeeper(t)
	type args struct {
		msg *types.MsgUpdate
	}

	u1 := uuid.NewString()
	k.SetMarket(ctx, types.Market{
		UID:     u1,
		Creator: sample.AccAddress(),
	})

	tests := []struct {
		name    string
		args    args
		want    *types.MsgUpdateResponse
		wantErr error
	}{
		{
			name: "test the empty or invalid format ticket",
			args: args{
				msg: types.NewMsgUpdate(sample.AccAddress(), ""),
			},
			want:    nil,
			wantErr: types.ErrInTicketVerification,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := msgk.Update(wctx, tt.args.msg)
			require.ErrorIs(t, err, tt.wantErr)
			require.EqualValues(t, got, tt.want)
		})
	}
}

func TestMsgServerUpdateMarketResponse(t *testing.T) {
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

		response, err := msgk.Update(
			wctx,
			types.NewMsgUpdate(sample.AccAddress(), validEmptyTicket),
		)
		assert.ErrorIs(t, err, types.ErrMarketNotFound)
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

		response, err := msgk.Update(
			wctx,
			types.NewMsgUpdate(sample.AccAddress(), validEmptyTicket),
		)
		assert.ErrorIs(t, err, types.ErrMarketCanNotBeAltered)
		assert.Nil(t, response)
	})
}
