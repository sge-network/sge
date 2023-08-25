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

func TestMsgServerDeposit(t *testing.T) {
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
	}

	tApp.MarketKeeper.SetMarket(ctx, marketItem)

	var oddsUIDs []string
	for _, v := range marketItem.Odds {
		oddsUIDs = append(oddsUIDs, v.UID)
	}
	err := tApp.OrderbookKeeper.InitiateOrderBook(ctx, marketItem.UID, oddsUIDs)
	require.NoError(t, err)

	t.Run("min deposit", func(t *testing.T) {
		inputDeposit := &types.MsgDeposit{
			Creator: creator.Address.String(),
			Amount:  sdk.NewInt(1),
		}

		_, err := msgk.Deposit(wctx, inputDeposit)
		require.ErrorIs(t, types.ErrDepositTooSmall, err)
	})

	t.Run("no ticket", func(t *testing.T) {
		inputDeposit := &types.MsgDeposit{
			Creator: creator.Address.String(),
			Amount:  sdk.NewInt(1000),
		}

		_, err := msgk.Deposit(wctx, inputDeposit)
		require.ErrorIs(t, types.ErrInTicketVerification, err)
	})

	t.Run("no authorization found", func(t *testing.T) {
		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       creator.Address.String(),
		}
		ticketClaim := jwt.MapClaims{
			"exp":               time.Now().Add(time.Minute * 5).Unix(),
			"iat":               time.Now().Unix(),
			"kyc_data":          testKyc,
			"depositor_address": depositor.Address.String(),
		}
		ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
		require.Nil(t, err)

		inputDeposit := &types.MsgDeposit{
			Creator:   creator.Address.String(),
			MarketUID: testMarketUID,
			Amount:    sdk.NewInt(1000),
			Ticket:    ticket,
		}

		_, err = msgk.Deposit(wctx, inputDeposit)
		require.ErrorIs(t, types.ErrAuthorizationNotFound, err)
	})

	t.Run("success without authorization", func(t *testing.T) {
		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       depositor.Address.String(),
		}
		ticketClaim := jwt.MapClaims{
			"exp":      time.Now().Add(time.Minute * 5).Unix(),
			"iat":      time.Now().Unix(),
			"kyc_data": testKyc,
		}
		ticket, err := simappUtil.CreateJwtTicket(ticketClaim)
		require.Nil(t, err)

		inputDeposit := &types.MsgDeposit{
			Creator:   depositor.Address.String(),
			MarketUID: testMarketUID,
			Amount:    sdk.NewInt(1000),
			Ticket:    ticket,
		}

		depResp, err := msgk.Deposit(wctx, inputDeposit)
		require.NoError(t, err)
		rst, found := k.GetDeposit(ctx,
			depositor.Address.String(),
			testMarketUID,
			depResp.ParticipationIndex,
		)
		require.True(t, found)
		require.Equal(t, inputDeposit.Creator, rst.Creator)
	})

	t.Run("success with authorization", func(t *testing.T) {
		grantAmount := sdk.NewInt(1000)

		expTime := time.Now().Add(5 * time.Minute)
		err := tApp.AuthzKeeper.SaveGrant(ctx,
			creator.Address,
			depositor.Address,
			types.NewDepositAuthorization(grantAmount),
			&expTime,
		)
		require.NoError(t, err)

		authzBefore, _ := tApp.AuthzKeeper.GetAuthorization(
			ctx,
			creator.Address,
			depositor.Address,
			sdk.MsgTypeURL(&types.MsgDeposit{}),
		)
		authzBeforeW, ok := authzBefore.(*types.DepositAuthorization)
		require.True(t, ok)
		require.Equal(t, grantAmount, authzBeforeW.SpendLimit)

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

		inputDeposit := &types.MsgDeposit{
			Creator:   creator.Address.String(),
			MarketUID: testMarketUID,
			Amount:    sdk.NewInt(1000),
			Ticket:    ticket,
		}

		depResp, err := msgk.Deposit(wctx, inputDeposit)
		require.NoError(t, err)
		rst, found := k.GetDeposit(ctx,
			depositor.Address.String(),
			testMarketUID,
			depResp.ParticipationIndex,
		)
		require.True(t, found)
		require.Equal(t, inputDeposit.Creator, rst.Creator)

		participation, found := tApp.OrderbookKeeper.GetOrderBookParticipation(
			ctx,
			testMarketUID,
			depResp.ParticipationIndex,
		)
		require.True(t, found)
		require.Equal(t, inputDeposit.Amount, participation.Liquidity.Add(participation.Fee))

		authzAfter, _ := tApp.AuthzKeeper.GetAuthorization(ctx,
			creator.Address,
			depositor.Address,
			sdk.MsgTypeURL(&types.MsgDeposit{}),
		)
		require.Nil(t, authzAfter)
	})
}
