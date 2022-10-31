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
	dvmtypes "github.com/sge-network/sge/x/dvm/types"

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
			Bet: &types.BetPlaceFields{
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
			Bet: &types.BetPlaceFields{
				UID:    "betUID_1",
				Amount: sdk.NewInt(500),
				Ticket: "wrongTicket",
			},
		}

		_, err = msgk.PlaceBet(wctx, inputBet)
		require.ErrorIs(t, types.ErrInVerification, err)
	})

	t.Run("Error in ticket fields validation", func(t *testing.T) {
		activeOdds := make([]*types.BetOdds, len(testActiveBetOdds))
		for i, v := range testActiveBetOdds {
			shallow := *v
			activeOdds[i] = &shallow
		}

		activeOdds[0].SportEventUID = ""
		placeBetClaim := jwt.MapClaims{
			"value":    sdk.NewDec(10),
			"odds_uid": testOddsUID1,
			"exp":      9999999999,
			"iat":      1111111111,
			"odds":     activeOdds,
		}
		placeBetTicket, err := createJwtTicket(placeBetClaim)
		require.Nil(t, err)

		inputBet := &types.MsgPlaceBet{

			Creator: creator.Address.String(),

			Bet: &types.BetPlaceFields{
				UID:    "betUID_1",
				Amount: sdk.NewInt(500),
				Ticket: placeBetTicket,
			},
		}

		_, err = msgk.PlaceBet(wctx, inputBet)
		require.Equal(t, types.ErrInvalidSportEventUID, err)
	})

	t.Run("No matching sportEvent", func(t *testing.T) {
		placeBetClaim := jwt.MapClaims{
			"value":    sdk.NewDec(10),
			"odds_uid": testOddsUID1,
			"exp":      9999999999,
			"iat":      1111111111,
			"odds":     testActiveBetOdds,
		}
		placeBetTicket, err := createJwtTicket(placeBetClaim)
		require.Nil(t, err)

		inputBet := &types.MsgPlaceBet{
			Creator: creator.Address.String(),

			Bet: &types.BetPlaceFields{
				UID:    "betUID_1",
				Amount: sdk.NewInt(500),
				Ticket: placeBetTicket,
			},
		}

		_, err = msgk.PlaceBet(wctx, inputBet)
		require.Equal(t, types.ErrNoMatchingSportEvent, err)
	})

	t.Run("Success", func(t *testing.T) {
		placeBetClaim := jwt.MapClaims{
			"value":    sdk.NewDec(10),
			"odds_uid": testOddsUID1,
			"exp":      9999999999,
			"iat":      1111111111,
			"odds":     testActiveBetOdds,
		}
		placeBetTicket, err := createJwtTicket(placeBetClaim)
		require.Nil(t, err)

		inputBet := &types.MsgPlaceBet{

			Creator: creator.Address.String(),
			Bet: &types.BetPlaceFields{
				UID:    "BetUID_2",
				Amount: sdk.NewInt(500),
				Ticket: placeBetTicket,
			},
		}

		totalOddsStat := make(map[string]*sporteventtypes.TotalOddsStats)
		totalOddsStat[testEventOddsUIDs[0]] = &sporteventtypes.TotalOddsStats{
			ExtraPayout: sdk.NewInt(0),
			BetAmount:   sdk.NewInt(0),
		}
		totalOddsStat[testEventOddsUIDs[1]] = &sporteventtypes.TotalOddsStats{
			ExtraPayout: sdk.NewInt(0),
			BetAmount:   sdk.NewInt(0),
		}

		sportEventItem := sporteventtypes.SportEvent{
			UID:      testSportEventUID,
			Creator:  testCreator,
			StartTS:  1111111111,
			EndTS:    9999999999,
			OddsUIDs: testEventOddsUIDs,
			Status:   sporteventtypes.SportEventStatus_STATUS_PENDING,
			Active:   true,
			BetConstraints: &sporteventtypes.EventBetConstraints{
				MaxBetCap:      sdk.NewInt(10000000000000),
				MinAmount:      sdk.NewInt(1),
				BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
				MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
				MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
				MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
				TotalOddsStats: totalOddsStat,
				TotalStats: &sporteventtypes.TotalStats{
					HouseLoss: sdk.NewInt(0),
					BetAmount: sdk.NewInt(0),
				},
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

func TestBetMsgServerPlaceBetSlip(t *testing.T) {
	tApp, k, msgk, ctx, wctx := setupMsgServerAndApp(t)
	creator := simappUtil.TestParamUsers["user1"]
	var err error

	bets := []*types.Bet{
		{
			UID: "invalid", // error: invalid UID
		},
		{
			UID:    "7e31c60f-2025-48ce-ae79-1dc110f16356",
			Amount: sdk.NewInt(int64(10)),
			Ticket: "invalidTicket", // err in verifying ticket
		},
		{
			UID:           "7e31c60f-2025-48ce-ae79-1dc110f16356",
			SportEventUID: "invalidUID", // error: ErrInvalidSportEventUID
			OddsUID:       testOddsUID1,
			OddsValue:     sdk.NewDec(int64(10)),
			Amount:        sdk.NewInt(int64(10)),
		},
		{
			UID:           "5e31c60f-2025-48ce-ae79-1dc110f16356",
			SportEventUID: "5e31c60f-2025-48ce-ae79-1dc110f16357", //error: no matching sport event
			OddsUID:       testOddsUID1,
			OddsValue:     sdk.NewDec(int64(10)),
			Amount:        sdk.NewInt(int64(10)),
		},
		{
			UID:           "6e31c60f-2025-48ce-ae79-1dc110f16356",
			SportEventUID: testSportEventUID, // no error
			OddsUID:       testOddsUID1,
			OddsValue:     sdk.NewDec(int64(10)),
			Amount:        sdk.NewInt(int64(10)),
		},
	}
	inputBets := &types.MsgPlaceBetSlip{
		Creator: creator.Address.String(),
		Bets:    []*types.BetPlaceFields{},
	}
	for _, bet := range bets {
		placeBetTicket := bet.Ticket
		if placeBetTicket == "" {
			activeOdds := make([]*types.BetOdds, len(testActiveBetOdds))
			for i, v := range testActiveBetOdds {
				shallow := *v
				activeOdds[i] = &shallow
				if bet.SportEventUID != "" {
					activeOdds[i].SportEventUID = bet.SportEventUID
				}
			}

			placeBetClaim := jwt.MapClaims{
				"value":    bet.OddsValue,
				"odds_uid": bet.OddsUID,
				"exp":      9999999999,
				"iat":      1111111111,
				"odds":     activeOdds,
			}
			placeBetTicket, err = createJwtTicket(placeBetClaim)
			require.Nil(t, err)
		}
		inputBets.Bets = append(inputBets.Bets, &types.BetPlaceFields{
			UID:    bet.UID,
			Amount: bet.Amount,
			Ticket: placeBetTicket,
		})

	}

	totalOddsStat := make(map[string]*sporteventtypes.TotalOddsStats)
	totalOddsStat[testOddsUID1] = &sporteventtypes.TotalOddsStats{
		ExtraPayout: sdk.NewInt(0),
		BetAmount:   sdk.NewInt(0),
	}
	totalOddsStat[testOddsUID2] = &sporteventtypes.TotalOddsStats{
		ExtraPayout: sdk.NewInt(0),
		BetAmount:   sdk.NewInt(0),
	}

	sportEventItem := sporteventtypes.SportEvent{
		UID:      bets[4].SportEventUID,
		EndTS:    99999999999,
		OddsUIDs: testEventOddsUIDs,
		Active:   true,
		BetConstraints: &sporteventtypes.EventBetConstraints{
			MaxBetCap:      sdk.NewInt(10000000000000),
			MinAmount:      sdk.NewInt(1),
			BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
			MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
			MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
			MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
			TotalOddsStats: totalOddsStat,
			TotalStats: &sporteventtypes.TotalStats{
				HouseLoss: sdk.NewInt(0),
				BetAmount: sdk.NewInt(0),
			},
		},
	}
	tApp.SporteventKeeper.SetSportEvent(ctx, sportEventItem)
	k.SetBet(ctx, *bets[0])
	expected := &types.MsgPlaceBetSlipResponse{
		SuccessfulBetUIDsList: []string{inputBets.Bets[4].UID},
		FailedBetUIDsErrorMap: map[string]string{
			inputBets.Bets[0].UID: types.ErrDuplicateUID.Error(),
			inputBets.Bets[1].UID: dvmtypes.ErrInvalidTicketFormat.Error(),
			inputBets.Bets[2].UID: types.ErrInvalidSportEventUID.Error(),
			inputBets.Bets[3].UID: types.ErrNoMatchingSportEvent.Error(),
		},
	}
	resp, err := msgk.PlaceBetSlip(wctx, inputBets)
	require.NoError(t, err)
	require.Equal(t, expected, resp)
	rst, found := k.GetBet(ctx,
		inputBets.Bets[4].UID,
	)
	require.True(t, found)
	require.Equal(t, inputBets.Creator, rst.Creator)
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
				OddsValue:     sdk.NewDec(10),
				Amount:        sdk.NewInt(500),
				Creator:       creator.Address.String(),
				OddsUID:       testOddsUID1,
				Ticket:        "Ticket",

				Result: types.Bet_RESULT_WON,
			},
			sportEvent: &sporteventtypes.SportEvent{
				UID:      testSportEventUID,
				Creator:  creator.Address.String(),
				StartTS:  1111111111,
				EndTS:    9999999999,
				OddsUIDs: testEventOddsUIDs,
				Status:   sporteventtypes.SportEventStatus_STATUS_RESULT_DECLARED,
			},
		},
	}

	totalOddsStat := make(map[string]*sporteventtypes.TotalOddsStats)
	totalOddsStat[testEventOddsUIDs[0]] = &sporteventtypes.TotalOddsStats{
		ExtraPayout: sdk.NewInt(0),
		BetAmount:   sdk.NewInt(0),
	}
	totalOddsStat[testEventOddsUIDs[1]] = &sporteventtypes.TotalOddsStats{
		ExtraPayout: sdk.NewInt(0),
		BetAmount:   sdk.NewInt(0),
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			betUID := uuid.NewString()

			if tc.bet != nil {
				resetSportEvent := sporteventtypes.SportEvent{
					UID:      testSportEventUID,
					Creator:  creator.Address.String(),
					StartTS:  1111111111,
					EndTS:    9999999999,
					OddsUIDs: testEventOddsUIDs,
					Status:   sporteventtypes.SportEventStatus_STATUS_PENDING,
					Active:   true,
					BetConstraints: &sporteventtypes.EventBetConstraints{
						MaxBetCap:      sdk.NewInt(10000000000000),
						MinAmount:      sdk.NewInt(1),
						BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
						MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
						MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
						MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
						TotalOddsStats: totalOddsStat,
						TotalStats: &sporteventtypes.TotalStats{
							HouseLoss: sdk.NewInt(0),
							BetAmount: sdk.NewInt(0),
						},
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
				require.Nil(t, res)
				return
			}
			require.NoError(t, err)
			require.Equal(t, expectedResp, res)
		})
	}
}

func TestBetMsgServerSettleBetBulk(t *testing.T) {
	tApp, _, msgk, ctx, wctx := setupMsgServerAndApp(t)
	creator := simappUtil.TestParamUsers["user1"]
	var err error

	inputMsg := &types.MsgSettleBetBulk{
		Creator: creator.Address.String(),
		BetUIDs: []string{
			"InvalidUID",
			uuid.NewString(),
		},
	}

	totalOddsStat := make(map[string]*sporteventtypes.TotalOddsStats)
	totalOddsStat[testEventOddsUIDs[0]] = &sporteventtypes.TotalOddsStats{
		ExtraPayout: sdk.NewInt(0),
		BetAmount:   sdk.NewInt(0),
	}
	totalOddsStat[testEventOddsUIDs[1]] = &sporteventtypes.TotalOddsStats{
		ExtraPayout: sdk.NewInt(0),
		BetAmount:   sdk.NewInt(0),
	}

	resetSportEvent := sporteventtypes.SportEvent{
		UID:      testSportEventUID,
		Creator:  testCreator,
		StartTS:  1111111111,
		EndTS:    9999999999,
		OddsUIDs: testEventOddsUIDs,
		Status:   sporteventtypes.SportEventStatus_STATUS_PENDING,
		Active:   true,
		BetConstraints: &sporteventtypes.EventBetConstraints{
			MaxBetCap:      sdk.NewInt(10000000000000),
			MinAmount:      sdk.NewInt(1),
			BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
			MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
			MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
			MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
			TotalOddsStats: totalOddsStat,
			TotalStats: &sporteventtypes.TotalStats{
				HouseLoss: sdk.NewInt(0),
				BetAmount: sdk.NewInt(0),
			},
		},
	}

	updateSportEvent := &sporteventtypes.SportEvent{
		UID:      testSportEventUID,
		Creator:  testCreator,
		StartTS:  1111111111,
		EndTS:    9999999999,
		OddsUIDs: testEventOddsUIDs,
		Active:   true,
		BetConstraints: &sporteventtypes.EventBetConstraints{
			MaxBetCap:      sdk.NewInt(10000000000000),
			MinAmount:      sdk.NewInt(1),
			BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
			MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
			MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
			MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
			TotalOddsStats: totalOddsStat,
			TotalStats: &sporteventtypes.TotalStats{
				HouseLoss: sdk.NewInt(0),
				BetAmount: sdk.NewInt(0),
			},
		},

		Status: sporteventtypes.SportEventStatus_STATUS_RESULT_DECLARED,
	}

	tApp.SporteventKeeper.SetSportEvent(ctx, resetSportEvent)
	placeTestBet(ctx, t, tApp, inputMsg.BetUIDs[1])
	tApp.SporteventKeeper.SetSportEvent(ctx, *updateSportEvent)

	expected := &types.MsgSettleBetBulkResponse{
		SuccessfulBetUIDsList: []string{inputMsg.BetUIDs[1]},
		FailedBetUIDsErrorMap: map[string]string{
			inputMsg.BetUIDs[0]: types.ErrInvalidBetUID.Error(),
		},
	}

	resp, err := msgk.SettleBetBulk(wctx, inputMsg)
	require.NoError(t, err)
	require.Equal(t, expected, resp)
}
