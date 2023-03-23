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

func TestMsgServerResolveMarket(t *testing.T) {
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
				msg: types.NewMsgResolveMarket(sample.AccAddress(), ""),
			},
			want:    nil,
			wantErr: types.ErrInTicketVerification,
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

func TestMsgServerResolveMarketResponse(t *testing.T) {
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

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveMarket(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrInTicketPayloadValidation)
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

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveMarket(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrMarketNotFound)
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

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveMarket(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrMarketResolutionNotAllowed)
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

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveMarket(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrInTicketPayloadValidation)
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

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveMarket(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrInTicketPayloadValidation)
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

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveMarket(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrInTicketPayloadValidation)
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

		response, err := msgk.ResolveMarket(wctx, types.NewMsgResolveMarket(sample.AccAddress(), validEmptyTicket))
		assert.ErrorIs(t, err, types.ErrInvalidWinnerOdds)
		assert.Nil(t, response)
	})
}

func TestValidateResolveMarket(t *testing.T) {
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
			err := k.ValidateMarketResolution(tt.msg)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
