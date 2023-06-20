package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	sgetypes "github.com/sge-network/sge/types"
	"github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestMsgServerWithdraw(t *testing.T) {
	tApp, k, msgk, ctx, wctx := setupMsgServerAndApp(t)
	creator := simappUtil.TestParamUsers["user1"]
	depositor := simappUtil.TestParamUsers["user2"]
	// var err error

	marketItem := markettypes.Market{
		UID:     testMarketUID,
		Creator: creator.Address.String(),
		StartTS: cast.ToUint64(time.Now().Unix()),
		EndTS:   cast.ToUint64(ctx.BlockTime().Unix()) + 1000,
		Odds:    testMarketOdds,
		Status:  markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
		BetConstraints: &markettypes.MarketBetConstraints{
			MinAmount: sdk.NewInt(1),
			BetFee:    sdk.NewInt(1),
		},
	}

	tApp.MarketKeeper.SetMarket(ctx, marketItem)

	var oddsUIDs []string
	for _, v := range marketItem.Odds {
		oddsUIDs = append(oddsUIDs, v.UID)
	}
	err := tApp.OrderbookKeeper.InitiateOrderBook(ctx, marketItem.UID, oddsUIDs)
	require.NoError(t, err)

	err = tApp.AuthzKeeper.SaveGrant(ctx,
		creator.Address,
		depositor.Address,
		types.NewDepositAuthorization(sdk.NewInt(1000)),
		time.Now().Add(5*time.Minute),
	)
	require.NoError(t, err)

	depositKyc := &sgetypes.KycDataPayload{
		Approved: true,
		ID:       depositor.Address.String(),
	}
	ticketClaim := jwt.MapClaims{
		"exp":               time.Now().Add(time.Minute * 5).Unix(),
		"iat":               time.Now().Unix(),
		"depositor_address": depositor.Address.String(),
		"kyc_data":          depositKyc,
	}
	ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
	require.Nil(t, err)

	inputDeposit := &types.MsgDeposit{
		Creator:   creator.Address.String(),
		MarketUID: testMarketUID,
		Amount:    sdk.NewInt(1000),
		Ticket:    ticket,
	}

	deposit, err := msgk.Deposit(wctx, inputDeposit)
	require.NoError(t, err)

	t.Run("no ticket", func(t *testing.T) {
		inputWithdraw := &types.MsgWithdraw{
			Creator: creator.Address.String(),
			Amount:  sdk.NewInt(1),
		}

		_, err := msgk.Withdraw(wctx, inputWithdraw)
		require.ErrorIs(t, types.ErrInTicketVerification, err)
	})

	t.Run("no authorization found", func(t *testing.T) {
		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       depositor.Address.String(),
		}
		ticketClaim := jwt.MapClaims{
			"exp":               time.Now().Add(time.Minute * 5).Unix(),
			"iat":               time.Now().Unix(),
			"kyc_data":          testKyc,
			"depositor_address": depositor.Address.String(),
		}
		ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
		require.Nil(t, err)

		inputWithdraw := &types.MsgWithdraw{
			Creator:            creator.Address.String(),
			MarketUID:          testMarketUID,
			ParticipationIndex: deposit.ParticipationIndex,
			Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
			Amount:             sdk.NewInt(1000),
			Ticket:             ticket,
		}

		_, err = msgk.Withdraw(wctx, inputWithdraw)
		require.ErrorIs(t, types.ErrAuthorizationNotFound, err)
	})

	t.Run("success", func(t *testing.T) {
		grantAmount := sdk.NewInt(1000)
		err := tApp.AuthzKeeper.SaveGrant(ctx,
			creator.Address,
			depositor.Address,
			types.NewWithdrawAuthorization(grantAmount),
			time.Now().Add(5*time.Minute),
		)
		require.NoError(t, err)

		authzBefore, _ := tApp.AuthzKeeper.GetCleanAuthorization(
			ctx,
			creator.Address,
			depositor.Address,
			sdk.MsgTypeURL(&types.MsgWithdraw{}),
		)
		authzBeforeW := authzBefore.(*types.WithdrawAuthorization)
		require.Equal(t, grantAmount, authzBeforeW.WithdrawLimit)

		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       depositor.Address.String(),
		}
		ticketClaim := jwt.MapClaims{
			"exp":               time.Now().Add(time.Minute * 5).Unix(),
			"iat":               time.Now().Unix(),
			"depositor_address": depositor.Address.String(),
			"kyc_data":          testKyc,
		}
		ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
		require.Nil(t, err)

		inputWithdraw := &types.MsgWithdraw{
			Creator:            creator.Address.String(),
			MarketUID:          testMarketUID,
			Amount:             sdk.NewInt(1000),
			ParticipationIndex: deposit.ParticipationIndex,
			Mode:               types.WithdrawalMode_WITHDRAWAL_MODE_FULL,
			Ticket:             ticket,
		}

		_, err = msgk.Withdraw(wctx, inputWithdraw)
		require.NoError(t, err)
		rst, err := k.GetAllWithdrawals(ctx)
		require.NoError(t, err)
		require.Equal(t, inputWithdraw.ParticipationIndex, rst[0].ParticipationIndex)

		authzAfter, _ := tApp.AuthzKeeper.GetCleanAuthorization(ctx,
			creator.Address,
			depositor.Address,
			sdk.MsgTypeURL(&types.MsgWithdraw{}),
		)
		authzAfterW := authzAfter.(*types.WithdrawAuthorization)
		expectedAuthzGrant := grantAmount.Sub(sdk.NewInt(rst[0].Amount.Int64()))
		require.Equal(t, expectedAuthzGrant, authzAfterW.WithdrawLimit)
	})
}
