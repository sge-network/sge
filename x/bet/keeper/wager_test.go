package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/require"
)

func TestWager(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	ctx = ctx.WithBlockTime(time.Now())
	tcs := []struct {
		desc          string
		bet           *types.Bet
		err           error
		market        *markettypes.Market
		activeBetOdds []*types.BetOdds
		betOdds       map[string]*types.BetOdds
	}{
		{
			desc: "invalid creator address",
			bet: &types.Bet{
				UID:       "betUID",
				MarketUID: "notExistMarketUID",
			},
			betOdds: map[string]*types.BetOdds{
				"odds1": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("0.1")},
				"odds2": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("1.0")},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "not found market",
			bet: &types.Bet{
				UID:       "betUID",
				MarketUID: "notExistMarketUID",
				Creator:   simappUtil.TestParamUsers["user1"].Address.String(),
			},
			betOdds: map[string]*types.BetOdds{
				"odds1": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("0.1")},
				"odds2": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("1.0")},
			},
			err: types.ErrNoMatchingMarket,
		},
		{
			desc: "inactive market",
			market: &markettypes.Market{
				UID:    "uid_inactive",
				Status: markettypes.MarketStatus_MARKET_STATUS_INACTIVE,
			},
			bet: &types.Bet{
				UID:       "betUID",
				MarketUID: "uid_inactive",
				Creator:   simappUtil.TestParamUsers["user1"].Address.String(),
			},
			betOdds: map[string]*types.BetOdds{
				"odds1": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("0.1")},
				"odds2": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("1.0")},
			},
			err: types.ErrInactiveMarket,
		},
		{
			desc: "not active market",
			market: &markettypes.Market{
				UID:    "uid_declared",
				Status: markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
			bet: &types.Bet{
				UID:       "betUID",
				MarketUID: "uid_declared",
				Creator:   simappUtil.TestParamUsers["user1"].Address.String(),
			},
			betOdds: map[string]*types.BetOdds{
				"odds1": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("0.1")},
				"odds2": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("1.0")},
			},
			err: types.ErrInactiveMarket,
		},
		{
			desc: "expired market",
			market: &markettypes.Market{
				UID:    "uid_expired",
				Status: markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
				EndTS:  0o00000000,
			},
			bet: &types.Bet{
				UID:       "betUID",
				MarketUID: "uid_expired",
				Creator:   simappUtil.TestParamUsers["user1"].Address.String(),
			},
			betOdds: map[string]*types.BetOdds{
				"odds1": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("0.1")},
				"odds2": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("1.0")},
			},
			err: types.ErrEndTSIsPassed,
		},
		{
			desc: "not exist odds",
			market: &markettypes.Market{
				UID:    "uid_oddsNotexist",
				Status: markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
				EndTS:  uint64(ctx.BlockTime().Unix()) + 1000,
				Odds: []*markettypes.Odds{
					{UID: "odds1"},
					{UID: "odds2"},
				},
			},
			activeBetOdds: []*types.BetOdds{
				{UID: "odds1", MarketUID: "uid_oddsNotexist", Value: "2.52"},
				{UID: "odds2", MarketUID: "uid_oddsNotexist", Value: "1.50"},
			},
			bet: &types.Bet{
				UID:       "betUID",
				MarketUID: "uid_oddsNotexist",
				OddsUID:   "notExistOdds",
				Amount:    sdk.NewInt(1000),
				OddsValue: "5",
				OddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
				Creator:   simappUtil.TestParamUsers["user1"].Address.String(),
			},
			betOdds: map[string]*types.BetOdds{
				"odds1": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("0.1")},
				"odds2": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("1.0")},
			},
			err: types.ErrOddsUIDNotExist,
		},
		{
			desc: "low bet amount",
			market: &markettypes.Market{
				UID:    "uid_lowBetAmount",
				Status: markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
				EndTS:  uint64(ctx.BlockTime().Unix()) + 1000,
				Odds: []*markettypes.Odds{
					{UID: "odds1"},
					{UID: "odds2"},
				},
			},
			activeBetOdds: []*types.BetOdds{
				{UID: "odds1", MarketUID: "uid_lowBetAmount", Value: "2.52"},
				{UID: "odds2", MarketUID: "uid_lowBetAmount", Value: "1.50"},
			},
			bet: &types.Bet{
				UID:       "betUID",
				MarketUID: "uid_lowBetAmount",
				OddsUID:   "odds1",
				Amount:    sdk.NewInt(100),
				OddsValue: "5",
				OddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
				Creator:   simappUtil.TestParamUsers["user1"].Address.String(),
			},
			betOdds: map[string]*types.BetOdds{
				"odds1": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("0.1")},
				"odds2": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("1.0")},
			},
			err: types.ErrBetAmountIsLow,
		},
		{
			desc: "success",
			market: &markettypes.Market{
				UID:    "uid_success",
				Status: markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
				EndTS:  uint64(ctx.BlockTime().Unix()) + 1000,
				Odds: []*markettypes.Odds{
					{UID: "odds1"},
					{UID: "odds2"},
				},
			},
			activeBetOdds: []*types.BetOdds{
				{
					UID:               "odds1",
					MarketUID:         "uid_success",
					Value:             "2.52",
					MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
				},
				{
					UID:               "odds2",
					MarketUID:         "uid_success",
					Value:             "1.50",
					MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
				},
			},
			bet: &types.Bet{
				UID:               "betUID",
				MarketUID:         "uid_success",
				OddsUID:           "odds1",
				Amount:            sdk.NewInt(1000000),
				OddsValue:         "5",
				OddsType:          types.OddsType_ODDS_TYPE_DECIMAL,
				Creator:           simappUtil.TestParamUsers["user1"].Address.String(),
				MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
			},
			betOdds: map[string]*types.BetOdds{
				"odds1": {UID: "odds1", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("0.1")},
				"odds2": {UID: "odds2", MarketUID: "uid_success", Value: "5", MaxLossMultiplier: sdk.MustNewDecFromStr("1.0")},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.market != nil {
				tApp.MarketKeeper.SetMarket(ctx, *tc.market)

				var oddsUIDs []string
				for _, v := range tc.market.Odds {
					oddsUIDs = append(oddsUIDs, v.UID)
				}
				err := tApp.OrderbookKeeper.InitiateOrderBook(ctx, tc.market.UID, oddsUIDs)
				require.NoError(t, err)

				if tc.market.Status == markettypes.MarketStatus_MARKET_STATUS_ACTIVE {
					_, err = tApp.OrderbookKeeper.InitiateOrderBookParticipation(
						ctx,
						simappUtil.TestParamUsers["user1"].Address,
						tc.market.UID,
						sdk.NewInt(100000000),
						sdk.NewInt(1),
					)
					require.NoError(t, err)
				}
			}

			err := k.Wager(ctx, tc.bet, tc.betOdds)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
		})
	}
}
