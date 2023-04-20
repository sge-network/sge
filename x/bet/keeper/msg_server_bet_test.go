package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"

	simappUtil "github.com/sge-network/sge/testutil/simapp"
	sgetypes "github.com/sge-network/sge/types"
	"github.com/sge-network/sge/x/bet/types"

	markettypes "github.com/sge-network/sge/x/market/types"
)

func TestBetMsgServerPlaceBet(t *testing.T) {
	tApp, k, msgk, ctx, wctx := setupMsgServerAndApp(t)
	creator := simappUtil.TestParamUsers["user1"]
	var err error

	t.Run("Redundant UID", func(t *testing.T) {
		betItem := types.Bet{UID: "betUID"}
		inputMsg := &types.MsgPlaceBet{
			Creator: creator.Address.String(),
			Bet: &types.PlaceBetFields{
				UID: betItem.UID,
			},
		}

		k.SetBet(ctx, betItem, 1)
		_, err := msgk.PlaceBet(wctx, inputMsg)
		require.ErrorIs(t, types.ErrDuplicateUID, err)
	})

	t.Run("Error in verifying ticket", func(t *testing.T) {
		inputBet := &types.MsgPlaceBet{
			Creator: creator.Address.String(),
			Bet: &types.PlaceBetFields{
				UID:    "betUID_1",
				Amount: sdk.NewInt(500),
				Ticket: "wrongTicket",
			},
		}

		_, err = msgk.PlaceBet(wctx, inputBet)
		require.ErrorIs(t, types.ErrInTicketVerification, err)
	})

	t.Run("Error in ticket fields validation", func(t *testing.T) {
		selectedBetOdds := *testSelectedBetOdds

		selectedBetOdds.MarketUID = ""
		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       creator.Address.String(),
		}
		placeBetClaim := jwt.MapClaims{
			"exp":           9999999999,
			"iat":           1111111111,
			"selected_odds": selectedBetOdds,
			"kyc_data":      testKyc,
		}
		placeBetTicket, err := createJwtTicket(placeBetClaim)
		require.Nil(t, err)

		inputBet := &types.MsgPlaceBet{
			Creator: creator.Address.String(),

			Bet: &types.PlaceBetFields{
				UID:    "betUID_1",
				Amount: sdk.NewInt(500),
				Ticket: placeBetTicket,
			},
		}

		_, err = msgk.PlaceBet(wctx, inputBet)
		require.ErrorIs(t, types.ErrInTicketValidation, err)
	})

	t.Run("No matching market", func(t *testing.T) {
		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       creator.Address.String(),
		}
		placeBetClaim := jwt.MapClaims{
			"exp":           9999999999,
			"iat":           1111111111,
			"selected_odds": testSelectedBetOdds,
			"kyc_data":      testKyc,
			"odds_type":     types.OddsType_ODDS_TYPE_DECIMAL,
		}
		placeBetTicket, err := createJwtTicket(placeBetClaim)
		require.Nil(t, err)

		inputBet := &types.MsgPlaceBet{
			Creator: creator.Address.String(),

			Bet: &types.PlaceBetFields{
				UID:    "betUID_1",
				Amount: sdk.NewInt(500),
				Ticket: placeBetTicket,
			},
		}

		_, err = msgk.PlaceBet(wctx, inputBet)
		require.ErrorIs(t, types.ErrInBetPlacement, err)
	})

	t.Run("Success", func(t *testing.T) {
		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       creator.Address.String(),
		}
		placeBetClaim := jwt.MapClaims{
			"exp":           9999999999,
			"iat":           1111111111,
			"selected_odds": testSelectedBetOdds,
			"kyc_data":      testKyc,
			"odds_type":     types.OddsType_ODDS_TYPE_DECIMAL,
		}
		placeBetTicket, err := createJwtTicket(placeBetClaim)
		require.Nil(t, err)

		inputBet := &types.MsgPlaceBet{
			Creator: creator.Address.String(),
			Bet: &types.PlaceBetFields{
				UID:    "BetUID_2",
				Amount: sdk.NewInt(500),
				Ticket: placeBetTicket,
			},
		}

		marketItem := markettypes.Market{
			UID:     testMarketUID,
			Creator: testCreator,
			StartTS: 1111111111,
			EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
			Odds:    testMarketOdds,
			Status:  markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
			BetConstraints: &markettypes.MarketBetConstraints{
				MinAmount: sdk.NewInt(1),
				BetFee:    sdk.NewInt(1),
			},
			SrContributionForHouse: sdk.NewInt(50000),
		}

		tApp.MarketKeeper.SetMarket(ctx, marketItem)

		var oddsUIDs []string
		for _, v := range marketItem.Odds {
			oddsUIDs = append(oddsUIDs, v.UID)
		}
		err = tApp.StrategicReserveKeeper.InitiateOrderBook(ctx, marketItem.UID, marketItem.SrContributionForHouse, oddsUIDs)
		require.NoError(t, err)

		_, err = msgk.PlaceBet(wctx, inputBet)
		require.NoError(t, err)
		rst, found := k.GetBet(ctx,
			creator.Address.String(),
			1,
		)
		require.True(t, found)
		require.Equal(t, inputBet.Creator, rst.Creator)

		uid2ID, found := k.GetBetID(ctx, inputBet.Bet.UID)
		require.True(t, found)
		require.Equal(t, types.UID2ID{UID: inputBet.Bet.UID, ID: 1}, uid2ID)

		stats := k.GetBetStats(ctx)
		require.Equal(t, types.BetStats{Count: 1}, stats)

		inputBet.Bet.UID = "BetUID_3"
		_, err = msgk.PlaceBet(wctx, inputBet)
		require.NoError(t, err)
		rst, found = k.GetBet(ctx,
			creator.Address.String(),
			2,
		)
		require.True(t, found)
		require.Equal(t, inputBet.Creator, rst.Creator)

		uid2ID, found = k.GetBetID(ctx, inputBet.Bet.UID)
		require.True(t, found)
		require.Equal(t, types.UID2ID{UID: inputBet.Bet.UID, ID: 2}, uid2ID)

		stats = k.GetBetStats(ctx)
		require.Equal(t, types.BetStats{Count: 2}, stats)
	})
}
