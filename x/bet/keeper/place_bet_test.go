package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/app/params"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/types"
	sporteventtypes "github.com/sge-network/sge/x/sportevent/types"
	"github.com/stretchr/testify/require"
)

func TestPlaceBet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)

	tcs := []struct {
		desc       string
		bet        *types.Bet
		err        error
		sportEvent *sporteventtypes.SportEvent
	}{
		{
			desc: "invalid creator address",
			bet: &types.Bet{
				UID:           "betUID",
				SportEventUID: "notExistSportEventUID",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "not found sport event",
			bet: &types.Bet{
				UID:           "betUID",
				SportEventUID: "notExistSportEventUID",
				Creator:       simappUtil.TestParamUsers["user1"].Address.String(),
			},
			err: types.ErrNoMatchingSportEvent,
		},
		{
			desc: "inactive sport event",
			sportEvent: &sporteventtypes.SportEvent{
				UID:    "uid_inactive",
				Status: sporteventtypes.SportEventStatus_STATUS_RESULT_DECLARED,
				BetConstraints: &sporteventtypes.EventBetConstraints{
					MaxBetCap:      sdk.NewInt(10000000000000),
					MinAmount:      sdk.NewInt(1),
					BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
					MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
					MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
					MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
					TotalOddsStats: defaultTotalStats(),
					TotalStats: &sporteventtypes.TotalStats{
						HouseLoss: sdk.NewInt(0),
						BetAmount: sdk.NewInt(0),
					},
				},
			},
			bet: &types.Bet{
				UID:           "betUID",
				SportEventUID: "uid_inactive",
				Creator:       simappUtil.TestParamUsers["user1"].Address.String(),
			},
			err: types.ErrInactiveSportEvent,
		},
		{
			desc: "not pending sport event",
			sportEvent: &sporteventtypes.SportEvent{
				UID:    "uid_declared",
				Status: sporteventtypes.SportEventStatus_STATUS_RESULT_DECLARED,
				Active: true,
				BetConstraints: &sporteventtypes.EventBetConstraints{
					MaxBetCap:      sdk.NewInt(10000000000000),
					MinAmount:      sdk.NewInt(1),
					BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
					MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
					MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
					MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
					TotalOddsStats: defaultTotalStats(),
					TotalStats: &sporteventtypes.TotalStats{
						HouseLoss: sdk.NewInt(0),
						BetAmount: sdk.NewInt(0),
					},
				},
			},
			bet: &types.Bet{
				UID:           "betUID",
				SportEventUID: "uid_declared",
				Creator:       simappUtil.TestParamUsers["user1"].Address.String(),
			},
			err: types.ErrSportEventStatusNotPending,
		},
		{
			desc: "expired sport event",
			sportEvent: &sporteventtypes.SportEvent{
				UID:    "uid_expired",
				Status: sporteventtypes.SportEventStatus_STATUS_PENDING,
				EndTS:  000000000,
				Active: true,
				BetConstraints: &sporteventtypes.EventBetConstraints{
					MaxBetCap:      sdk.NewInt(10000000000000),
					MinAmount:      sdk.NewInt(1),
					BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
					MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
					MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
					MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
					TotalOddsStats: defaultTotalStats(),
					TotalStats: &sporteventtypes.TotalStats{
						HouseLoss: sdk.NewInt(0),
						BetAmount: sdk.NewInt(0),
					},
				},
			},
			bet: &types.Bet{
				UID:           "betUID",
				SportEventUID: "uid_expired",
				Creator:       simappUtil.TestParamUsers["user1"].Address.String(),
			},
			err: types.ErrEndTSIsPassed,
		},
		{
			desc: "not exist odds",
			sportEvent: &sporteventtypes.SportEvent{
				UID:      "uid_oddsNotexist",
				Status:   sporteventtypes.SportEventStatus_STATUS_PENDING,
				EndTS:    uint64(time.Now().Add(time.Hour).Unix()),
				OddsUIDs: []string{"odds1", "odds2"},
				Active:   true,
				BetConstraints: &sporteventtypes.EventBetConstraints{
					MaxBetCap:      sdk.NewInt(10000000000000),
					MinAmount:      sdk.NewInt(1),
					BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
					MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
					MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
					MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
					TotalOddsStats: defaultTotalStats(),
					TotalStats: &sporteventtypes.TotalStats{
						HouseLoss: sdk.NewInt(0),
						BetAmount: sdk.NewInt(0),
					},
				},
			},
			bet: &types.Bet{
				UID:           "betUID",
				SportEventUID: "uid_oddsNotexist",
				OddsUID:       "notExistOdds",
				Amount:        sdk.NewInt(1000),
				OddsValue:     sdk.NewDec(5),
				Creator:       simappUtil.TestParamUsers["user1"].Address.String(),
			},
			err: types.ErrOddsUIDNotExist,
		},
		{
			desc: "low bet amount",
			sportEvent: &sporteventtypes.SportEvent{
				UID:      "uid_lowBetAmount",
				Status:   sporteventtypes.SportEventStatus_STATUS_PENDING,
				EndTS:    uint64(time.Now().Add(time.Hour).Unix()),
				OddsUIDs: []string{"odds1", "odds2"},
				Active:   true,
				BetConstraints: &sporteventtypes.EventBetConstraints{
					MaxBetCap:      sdk.NewInt(10000000000000),
					MinAmount:      sdk.NewInt(1000),
					BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
					MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
					MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
					MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
					TotalOddsStats: defaultTotalStats(),
					TotalStats: &sporteventtypes.TotalStats{
						HouseLoss: sdk.NewInt(0),
						BetAmount: sdk.NewInt(0),
					},
				},
			},
			bet: &types.Bet{
				UID:           "betUID",
				SportEventUID: "uid_lowBetAmount",
				OddsUID:       "odds1",
				Amount:        sdk.NewInt(100),
				OddsValue:     sdk.NewDec(5),
				Creator:       simappUtil.TestParamUsers["user1"].Address.String(),
			},
			err: types.ErrBetAmountIsLow,
		},
		{
			desc: "success",
			sportEvent: &sporteventtypes.SportEvent{
				UID:      "uid_success",
				Status:   sporteventtypes.SportEventStatus_STATUS_PENDING,
				EndTS:    uint64(time.Now().Add(time.Hour).Unix()),
				OddsUIDs: []string{"odds1", "odds2"},
				Active:   true,
				BetConstraints: &sporteventtypes.EventBetConstraints{
					MaxBetCap:      sdk.NewInt(10000000000000),
					MinAmount:      sdk.NewInt(1),
					BetFee:         sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1)),
					MaxLoss:        sdk.NewInt(sporteventtypes.DefaultMaxEventLoss),
					MaxVig:         sdk.NewDec(sporteventtypes.DefaultMaxVig),
					MinVig:         sdk.NewDec(sporteventtypes.DefaultMinVig),
					TotalOddsStats: defaultTotalStats(),
					TotalStats: &sporteventtypes.TotalStats{
						HouseLoss: sdk.NewInt(0),
						BetAmount: sdk.NewInt(0),
					},
				},
			},
			bet: &types.Bet{
				UID:           "betUID",
				SportEventUID: "uid_success",
				OddsUID:       "odds1",
				Amount:        sdk.NewInt(1000),
				OddsValue:     sdk.NewDec(5),
				Creator:       simappUtil.TestParamUsers["user1"].Address.String(),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.sportEvent != nil {
				tApp.SporteventKeeper.SetSportEvent(ctx, *tc.sportEvent)
			}
			err := k.PlaceBet(ctx, tc.bet, testActiveBetOdds)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
