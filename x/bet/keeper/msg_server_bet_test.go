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

func TestBetMsgServerWager(t *testing.T) {
	tApp, k, msgk, ctx, wctx := setupMsgServerAndApp(t)
	creator := simappUtil.TestParamUsers["user1"]
	var err error

	t.Run("Redundant UID", func(t *testing.T) {
		betItem := types.Bet{UID: "betUID"}
		inputMsg := &types.MsgWager{
			Creator: creator.Address.String(),
			Props: &types.WagerProps{
				UID: betItem.UID,
			},
		}

		k.SetBet(ctx, betItem, 1)
		_, err := msgk.Wager(wctx, inputMsg)
		require.ErrorIs(t, types.ErrDuplicateUID, err)
	})

	t.Run("Error in verifying ticket", func(t *testing.T) {
		inputBet := &types.MsgWager{
			Creator: creator.Address.String(),
			Props: &types.WagerProps{
				UID:    "betUID_1",
				Amount: sdk.NewInt(500),
				Ticket: "wrongTicket",
			},
		}

		_, err = msgk.Wager(wctx, inputBet)
		require.ErrorIs(t, types.ErrInTicketVerification, err)
	})

	t.Run("Error in ticket fields validation", func(t *testing.T) {
		selectedBetOdds := *testSelectedBetOdds

		selectedBetOdds.MarketUID = ""
		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       creator.Address.String(),
		}
		wagerClaim := jwt.MapClaims{
			"exp":           9999999999,
			"iat":           1111111111,
			"selected_odds": selectedBetOdds,
			"kyc_data":      testKyc,
			"all_odds":      testBetOdds,
		}
		wagerTicket, err := createJwtTicket(wagerClaim)
		require.Nil(t, err)

		inputBet := &types.MsgWager{
			Creator: creator.Address.String(),

			Props: &types.WagerProps{
				UID:    "betUID_1",
				Amount: sdk.NewInt(500),
				Ticket: wagerTicket,
			},
		}

		_, err = msgk.Wager(wctx, inputBet)
		require.ErrorIs(t, types.ErrInTicketValidation, err)
	})

	t.Run("No matching market", func(t *testing.T) {
		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       creator.Address.String(),
		}
		wagerClaim := jwt.MapClaims{
			"exp":           9999999999,
			"iat":           1111111111,
			"selected_odds": testSelectedBetOdds,
			"kyc_data":      testKyc,
			"odds_type":     types.OddsType_ODDS_TYPE_DECIMAL,
			"all_odds":      testBetOdds,
		}
		wagerTicket, err := createJwtTicket(wagerClaim)
		require.Nil(t, err)

		inputBet := &types.MsgWager{
			Creator: creator.Address.String(),

			Props: &types.WagerProps{
				UID:    "betUID_1",
				Amount: sdk.NewInt(500),
				Ticket: wagerTicket,
			},
		}

		_, err = msgk.Wager(wctx, inputBet)
		require.ErrorIs(t, types.ErrInWager, err)
	})

	t.Run("Success", func(t *testing.T) {
		testKyc := &sgetypes.KycDataPayload{
			Approved: true,
			ID:       creator.Address.String(),
		}
		wagerClaim := jwt.MapClaims{
			"exp":           9999999999,
			"iat":           1111111111,
			"selected_odds": testSelectedBetOdds,
			"kyc_data":      testKyc,
			"odds_type":     types.OddsType_ODDS_TYPE_DECIMAL,
			"all_odds":      testBetOdds,
		}
		wagerTicket, err := createJwtTicket(wagerClaim)
		require.Nil(t, err)

		inputBet := &types.MsgWager{
			Creator: creator.Address.String(),
			Props: &types.WagerProps{
				UID:    "BetUID_2",
				Amount: sdk.NewInt(1000000),
				Ticket: wagerTicket,
			},
		}

		marketItem := markettypes.Market{
			UID:     testMarketUID,
			Creator: testCreator,
			StartTS: 1111111111,
			EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
			Odds:    testMarketOdds,
			Status:  markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
		}

		tApp.MarketKeeper.SetMarket(ctx, marketItem)

		var oddsUIDs []string
		for _, v := range marketItem.Odds {
			oddsUIDs = append(oddsUIDs, v.UID)
		}
		err = tApp.OrderbookKeeper.InitiateOrderBook(ctx, marketItem.UID, oddsUIDs)
		require.NoError(t, err)

		_, err = tApp.OrderbookKeeper.InitiateOrderBookParticipation(
			ctx,
			simappUtil.TestParamUsers["user1"].Address,
			marketItem.UID,
			sdk.NewInt(100000000),
			sdk.NewInt(1),
		)
		require.NoError(t, err)

		_, err = msgk.Wager(wctx, inputBet)
		require.NoError(t, err)
		rst, found := k.GetBet(ctx,
			creator.Address.String(),
			1,
		)
		require.True(t, found)
		require.Equal(t, inputBet.Creator, rst.Creator)

		uid2ID, found := k.GetBetID(ctx, inputBet.Props.UID)
		require.True(t, found)
		require.Equal(t, types.UID2ID{UID: inputBet.Props.UID, ID: 1}, uid2ID)

		stats := k.GetBetStats(ctx)
		require.Equal(t, types.BetStats{Count: 1}, stats)

		inputBet.Props.UID = "BetUID_3"
		_, err = msgk.Wager(wctx, inputBet)
		require.NoError(t, err)
		rst, found = k.GetBet(ctx,
			creator.Address.String(),
			2,
		)
		require.True(t, found)
		require.Equal(t, inputBet.Creator, rst.Creator)

		uid2ID, found = k.GetBetID(ctx, inputBet.Props.UID)
		require.True(t, found)
		require.Equal(t, types.UID2ID{UID: inputBet.Props.UID, ID: 2}, uid2ID)

		stats = k.GetBetStats(ctx)
		require.Equal(t, types.BetStats{Count: 2}, stats)
	})
}
