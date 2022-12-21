package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/x/bet/types"
	sporteventtypes "github.com/sge-network/sge/x/sportevent/types"
	"github.com/stretchr/testify/require"
)

func TestSettleBet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	addSportEvent(t, tApp, ctx)

	tcs := []struct {
		desc             string
		betUID           string
		updateSportEvent *sporteventtypes.SportEvent
		bet              *types.Bet
		err              error
	}{
		{
			desc:   "invalid betUID",
			betUID: "invalidUID",
			err:    types.ErrInvalidBetUID,
		},
		{
			desc:   "not found bet",
			betUID: "0db09053-2901-4110-8fb5-c14e21f8d400",
			err:    types.ErrNoMatchingBet,
		},
		{
			desc: "invalid creator",
			bet: &types.Bet{
				UID:           "0db09053-2901-4110-8fb5-c14e21f8d400",
				SportEventUID: "nonExistSportEvent",
				Creator:       "invalid creator",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "failed in checking status",
			bet: &types.Bet{
				SportEventUID: testSportEventUID,
				OddsValue:     sdk.NewDec(10),
				Amount:        sdk.NewInt(500),
				Creator:       testCreator,
				OddsUID:       testOddsUID1,
				Ticket:        "Ticket",

				Verified: true,
				Status:   types.Bet_STATUS_SETTLED,
			},
			err: types.ErrBetIsSettled,
		},
		{
			desc: "not found sport event",
			bet: &types.Bet{
				UID:           "0db09053-2901-4110-8fb5-c14e21f8d400",
				SportEventUID: "nonExistSportEvent",
				Creator:       testCreator,
			},
			err: types.ErrNoMatchingSportEvent,
		},
		{
			desc: "sport event is aborted",
			bet: &types.Bet{
				SportEventUID: testSportEventUID,
				OddsValue:     sdk.NewDec(10),
				Amount:        sdk.NewInt(500),
				Creator:       testCreator,
				OddsUID:       testOddsUID1,
				Ticket:        "Ticket",
			},
			updateSportEvent: &sporteventtypes.SportEvent{
				UID:     testSportEventUID,
				Creator: testCreator,
				StartTS: 1111111111,
				EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:    testEventOdds,

				Status: sporteventtypes.SportEventStatus_STATUS_ABORTED,
			},
		},
		{
			desc: "sport event is canceled",
			bet: &types.Bet{
				SportEventUID: testSportEventUID,
				OddsValue:     sdk.NewDec(10),
				Amount:        sdk.NewInt(300),
				Creator:       testCreator,
				OddsUID:       testOddsUID1,
				Ticket:        "Ticket",
			},
			updateSportEvent: &sporteventtypes.SportEvent{
				UID:     testSportEventUID,
				Creator: testCreator,
				StartTS: 1111111111,
				EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:    testEventOdds,

				Status: sporteventtypes.SportEventStatus_STATUS_CANCELLED,
			},
		},
		{
			desc: "result is not declared",
			bet: &types.Bet{
				SportEventUID: testSportEventUID,
				OddsValue:     sdk.NewDec(10),
				Amount:        sdk.NewInt(500),
				Creator:       testCreator,
				OddsUID:       testOddsUID1,
				Ticket:        "Ticket",
			},
			updateSportEvent: &sporteventtypes.SportEvent{
				UID:     testSportEventUID,
				Creator: testCreator,
				StartTS: 1111111111,
				EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:    testEventOdds,

				Status: sporteventtypes.SportEventStatus_STATUS_PENDING,
			},
			err: types.ErrResultNotDeclared,
		},
		{
			desc: "success",
			bet: &types.Bet{
				SportEventUID: testSportEventUID,
				OddsValue:     sdk.NewDec(10),
				Amount:        sdk.NewInt(500),
				Creator:       testCreator,
				OddsUID:       testOddsUID1,
				Ticket:        "Ticket",

				Result: types.Bet_RESULT_WON,
			},
			updateSportEvent: &sporteventtypes.SportEvent{
				UID:     testSportEventUID,
				Creator: testCreator,
				StartTS: 1111111111,
				EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:    testEventOdds,

				Status: sporteventtypes.SportEventStatus_STATUS_RESULT_DECLARED,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			betUID := uuid.NewString()

			if tc.bet != nil {
				resetSportEvent := sporteventtypes.SportEvent{
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
				tApp.SporteventKeeper.SetSportEvent(ctx, resetSportEvent)
				tc.bet.UID = betUID
				placeTestBet(ctx, t, tApp, betUID)
				k.SetBet(ctx, *tc.bet)
			}
			if tc.updateSportEvent != nil {
				tApp.SporteventKeeper.SetSportEvent(ctx, *tc.updateSportEvent)
			}
			if tc.betUID != "" {
				betUID = tc.betUID
			}
			err := k.SettleBet(ctx, betUID)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestCheckBetStatus(t *testing.T) {
	k, _ := setupKeeper(t)
	tcs := []struct {
		desc string
		bet  *types.Bet
		err  error
	}{
		{
			desc: "pending bet",
			bet: &types.Bet{
				Status: types.Bet_STATUS_PENDING,
			},
		},
		{
			desc: "canceled bet",
			bet: &types.Bet{
				Status: types.Bet_STATUS_CANCELLED,
			},
			err: types.ErrBetIsCanceled,
		},
		{
			desc: "settled bet",
			bet: &types.Bet{
				Status: types.Bet_STATUS_SETTLED,
			},
			err: types.ErrBetIsSettled,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := k.CheckBetStatus(tc.bet)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestResolveBetResult(t *testing.T) {
	k, _ := setupKeeper(t)
	tcs := []struct {
		desc       string
		bet        *types.Bet
		sportEvent sporteventtypes.SportEvent
		err        error
		result     types.Bet_Result
	}{
		{
			desc: "not declared",
			sportEvent: sporteventtypes.SportEvent{
				Status: sporteventtypes.SportEventStatus_STATUS_PENDING,
			},
			bet: &types.Bet{},
			err: types.ErrResultNotDeclared,
		},
		{
			desc: "won",
			sportEvent: sporteventtypes.SportEvent{
				Status:         sporteventtypes.SportEventStatus_STATUS_RESULT_DECLARED,
				WinnerOddsUIDs: []string{"oddsUID"},
			},
			bet: &types.Bet{
				OddsUID: "oddsUID",
			},
			result: types.Bet_RESULT_WON,
		},
		{
			desc: "lost",
			sportEvent: sporteventtypes.SportEvent{
				Status:         sporteventtypes.SportEventStatus_STATUS_RESULT_DECLARED,
				WinnerOddsUIDs: []string{"oddsUID"},
			},
			bet:    &types.Bet{},
			result: types.Bet_RESULT_LOST,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := k.ResolveBetResult(tc.bet, tc.sportEvent)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.Nil(t, err)
				require.Equal(t, tc.bet.Result, tc.result)
				require.Equal(t, tc.bet.Status, types.Bet_STATUS_RESULT_DECLARED)
			}
		})
	}
}
