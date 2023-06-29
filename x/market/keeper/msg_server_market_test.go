package keeper_test

import (
	"testing"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgServerAddMarket(t *testing.T) {
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
				msg: types.NewMsgAddMarket(sample.AccAddress(), ""),
			},
			want:    nil,
			wantErr: types.ErrInTicketVerification,
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

		response, err := msgk.AddMarket(
			wctx,
			types.NewMsgAddMarket(sample.AccAddress(), validEmptyTicket),
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

		response, err := msgk.AddMarket(
			wctx,
			types.NewMsgAddMarket(sample.AccAddress(), validEmptyTicket),
		)
		assert.ErrorIs(t, err, types.ErrMarketAlreadyExist)
		assert.Nil(t, response)
	})
}

func TestMsgServerUpdateMarket(t *testing.T) {
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
				msg: types.NewMsgUpdateMarket(sample.AccAddress(), ""),
			},
			want:    nil,
			wantErr: types.ErrInTicketVerification,
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

		response, err := msgk.UpdateMarket(
			wctx,
			types.NewMsgUpdateMarket(sample.AccAddress(), validEmptyTicket),
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

		response, err := msgk.UpdateMarket(
			wctx,
			types.NewMsgUpdateMarket(sample.AccAddress(), validEmptyTicket),
		)
		assert.ErrorIs(t, err, types.ErrMarketCanNotBeAltered)
		assert.Nil(t, response)
	})
}

func TestValidateCreationMarket(t *testing.T) {
	k, _, wctx, _ := setupMsgServerAndKeeper(t)
	tests := []struct {
		name string
		msg  types.MarketAddTicketPayload
		err  error
	}{
		{
			name: "valid request",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
				Meta:   "Winner of x:y",
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
		},
		{
			name: "same timestamp",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "end timestamp before current timestamp",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				EndTS:   uint64(time.Now().Add(-time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid uid",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
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
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
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
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
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
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: "8779cf93-925c-4818-bc81-13c359e0deb8", Meta: "Odds 1"},
					{UID: "8779cf93-925c-4818-bc81-13c359e0deb8", Meta: "invalid odds"},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid request, with bet constraint",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
				Meta:   "Winner of x:y",
				Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
		},
		{
			name: "large metadata",
			msg: types.MarketAddTicketPayload{
				Creator: sample.AccAddress(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				UID:     uuid.NewString(),
				Odds: []*types.Odds{
					{UID: uuid.NewString(), Meta: "Odds 1"},
					{UID: uuid.NewString(), Meta: "Odds 2"},
				},
				Meta: simappUtil.RandomString(types.MaxAllowedCharactersForMeta + 1),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := k.ValidateMarketAdd(wctx, tt.msg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestUpdateMarketValidation(t *testing.T) {
	k, _, wctx, _ := setupMsgServerAndKeeper(t)

	market := types.Market{
		Creator: sample.AccAddress(),
		StartTS: uint64(time.Now().Add(time.Minute).Unix()),
		EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
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
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_ACTIVE,
			},
		},
		{
			name: "valid request inactive",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_INACTIVE,
			},
		},
		{
			name: "invalid status, declared",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid status, canceled",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_CANCELED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid status, aborted",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_ABORTED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid status, unpecified",
			msg: types.MarketUpdateTicketPayload{
				UID:     uuid.NewString(),
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute * 2).Unix()),
				Status:  types.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "same timestamp",
			msg: types.MarketUpdateTicketPayload{
				StartTS: uint64(time.Now().Add(time.Minute).Unix()),
				EndTS:   uint64(time.Now().Add(time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "end timestamp before current timestamp",
			msg: types.MarketUpdateTicketPayload{
				EndTS: uint64(time.Now().Add(-time.Minute).Unix()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := k.ValidateMarketUpdate(wctx, tt.msg, market)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
