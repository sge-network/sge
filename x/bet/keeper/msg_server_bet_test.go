package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/app/params"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/types"

	sporteventtypes "github.com/sge-network/sge/x/sportevent/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

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

		k.SetBet(ctx, betItem)
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
		require.ErrorIs(t, types.ErrInVerification, err)
	})

	t.Run("Error in ticket fields validation", func(t *testing.T) {
		selectedBetOdds := *testSelectedBetOdds

		selectedBetOdds.SportEventUID = ""
		testKyc := &types.KycDataPayload{
			KycRequired: true,
			KycApproved: true,
			KycId:       creator.Address.String(),
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
		require.Equal(t, types.ErrInvalidSportEventUID, err)
	})

	t.Run("No matching sportEvent", func(t *testing.T) {
		testKyc := &types.KycDataPayload{
			KycRequired: true,
			KycApproved: true,
			KycId:       creator.Address.String(),
		}
		placeBetClaim := jwt.MapClaims{
			"exp":           9999999999,
			"iat":           1111111111,
			"selected_odds": testSelectedBetOdds,
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
		require.Equal(t, types.ErrNoMatchingSportEvent, err)
	})

	t.Run("Success", func(t *testing.T) {
		testKyc := &types.KycDataPayload{
			KycRequired: true,
			KycApproved: true,
			KycId:       creator.Address.String(),
		}
		placeBetClaim := jwt.MapClaims{
			"exp":           9999999999,
			"iat":           1111111111,
			"selected_odds": testSelectedBetOdds,
			"kyc_data":      testKyc,
		}
		placeBetTicket, err := createJwtTicket(placeBetClaim)
		require.Nil(t, err)

		inputBet := &types.MsgPlaceBet{

			Creator: creator.Address.String(),
			Bet: &types.PlaceBetFields{
				UID:      "BetUID_2",
				Amount:   sdk.NewInt(500),
				OddsType: types.OddsType_ODD_TYPE_DECIMAL,
				Ticket:   placeBetTicket,
			},
		}

		sportEventItem := sporteventtypes.SportEvent{
			UID:     testSportEventUID,
			Creator: testCreator,
			StartTS: 1111111111,
			EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
			Odds:    testEventOdds,
			Status:  sporteventtypes.SportEventStatus_STATUS_PENDING,
			Active:  true,
			BetConstraints: &sporteventtypes.EventBetConstraints{
				MinAmount: sdk.NewInt(1),
				BetFee:    sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
			},
		}

		tApp.SporteventKeeper.SetSportEvent(ctx, sportEventItem)
		_, err = msgk.PlaceBet(wctx, inputBet)
		require.NoError(t, err)
		rst, found := k.GetBet(ctx,
			inputBet.Bet.UID,
		)
		require.True(t, found)
		require.Equal(t, inputBet.Creator, rst.Creator)
	})
}

func TestBetMsgServerSettleBet(t *testing.T) {
	tApp, k, msgk, ctx, wctx := setupMsgServerAndApp(t)
	creator := simappUtil.TestParamUsers["user1"]

	tcs := []struct {
		desc       string
		betUID     string
		sportEvent *sporteventtypes.SportEvent
		bet        *types.Bet
		err        error
	}{
		{
			desc:   "invalid betUID",
			betUID: "invalidUID",
			err:    types.ErrInvalidBetUID,
		},
		{
			desc: "success",
			bet: &types.Bet{
				SportEventUID: testSportEventUID,
				OddsValue:     "10",
				OddsType:      types.OddsType_ODD_TYPE_DECIMAL,
				Amount:        sdk.NewInt(500),
				Creator:       creator.Address.String(),
				OddsUID:       testOddsUID1,
				Ticket:        "Ticket",

				Result: types.Bet_RESULT_WON,
			},
			sportEvent: &sporteventtypes.SportEvent{
				UID:     testSportEventUID,
				Creator: creator.Address.String(),
				StartTS: 1111111111,
				EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:    testEventOdds,
				Status:  sporteventtypes.SportEventStatus_STATUS_RESULT_DECLARED,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			betUID := uuid.NewString()

			if tc.bet != nil {
				resetSportEvent := sporteventtypes.SportEvent{
					UID:     testSportEventUID,
					Creator: creator.Address.String(),
					StartTS: 1111111111,
					EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
					Odds:    testEventOdds,
					Status:  sporteventtypes.SportEventStatus_STATUS_PENDING,
					Active:  true,
					BetConstraints: &sporteventtypes.EventBetConstraints{
						MinAmount: sdk.NewInt(1),
						BetFee:    sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
					},
				}
				tApp.SporteventKeeper.SetSportEvent(ctx, resetSportEvent)
				tc.bet.UID = betUID
				placeTestBet(ctx, t, tApp, betUID)
				k.SetBet(ctx, *tc.bet)
			}
			if tc.sportEvent != nil {
				tApp.SporteventKeeper.SetSportEvent(ctx, *tc.sportEvent)
			}
			if tc.betUID != "" {
				betUID = tc.betUID
			}
			inputMsg := &types.MsgSettleBet{
				Creator: creator.Address.String(),
				BetUID:  betUID,
			}
			expectedResp := &types.MsgSettleBetResponse{}
			res, err := msgk.SettleBet(wctx, inputMsg)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				require.Equal(t, &types.MsgSettleBetResponse{Error: tc.err.Error(), BetUID: betUID}, res)
				return
			}
			require.NoError(t, err)
			require.Equal(t, expectedResp, res)
		})
	}
}
