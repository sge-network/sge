package keeper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
)

func TestSettleBet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	testCreator = simapp.TestParamUsers["user1"].Address.String()
	addTestMarket(t, tApp, ctx)

	tcs := []struct {
		desc         string
		betUID       string
		updateMarket *markettypes.Market
		bet          *types.Bet
		err          error
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
				UID:       "0db09053-2901-4110-8fb5-c14e21f8d400",
				MarketUID: "nonExistMarket",
				Creator:   "invalid creator",
			},
			err: sdkerrtypes.ErrInvalidAddress,
		},
		{
			desc: "failed in checking status",
			bet: &types.Bet{
				MarketUID: testMarketUID,
				OddsValue: "10",
				Amount:    sdkmath.NewInt(1000000),
				Creator:   testCreator,
				OddsUID:   testOddsUID1,
				Status:    types.Bet_STATUS_SETTLED,
			},
			err: types.ErrBetIsSettled,
		},
		{
			desc: "not found market",
			bet: &types.Bet{
				UID:       "0db09053-2901-4110-8fb5-c14e21f8d400",
				MarketUID: "nonExistMarket",
				Creator:   testCreator,
			},
			err: types.ErrNoMatchingMarket,
		},
		{
			desc: "market is aborted",
			bet: &types.Bet{
				MarketUID: testMarketUID,
				OddsValue: "10",
				Amount:    sdkmath.NewInt(1000000),
				Creator:   testCreator,
				OddsUID:   testOddsUID1,
			},
			updateMarket: &markettypes.Market{
				UID:     testMarketUID,
				Creator: testCreator,
				StartTS: 1111111111,
				EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:    testMarketOdds,

				Status: markettypes.MarketStatus_MARKET_STATUS_ABORTED,
			},
		},
		{
			desc: "market is canceled",
			bet: &types.Bet{
				MarketUID: testMarketUID,
				OddsValue: "10",
				Amount:    sdkmath.NewInt(1000000),
				Creator:   testCreator,
				OddsUID:   testOddsUID1,
			},
			updateMarket: &markettypes.Market{
				UID:     testMarketUID,
				Creator: testCreator,
				StartTS: 1111111111,
				EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:    testMarketOdds,

				Status: markettypes.MarketStatus_MARKET_STATUS_CANCELED,
			},
		},
		{
			desc: "result is not declared",
			bet: &types.Bet{
				MarketUID: testMarketUID,
				OddsValue: "10",
				Amount:    sdkmath.NewInt(1000000),
				Creator:   testCreator,
				OddsUID:   testOddsUID1,
			},
			updateMarket: &markettypes.Market{
				UID:     testMarketUID,
				Creator: testCreator,
				StartTS: 1111111111,
				EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:    testMarketOdds,

				Status: markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
			},
			err: types.ErrResultNotDeclared,
		},
		{
			desc: "success",
			bet: &types.Bet{
				MarketUID: testMarketUID,
				OddsValue: "10",
				Amount:    sdkmath.NewInt(1000000),
				Creator:   testCreator,
				OddsUID:   testOddsUID1,

				Result: types.Bet_RESULT_WON,
			},
			updateMarket: &markettypes.Market{
				UID:     testMarketUID,
				Creator: testCreator,
				StartTS: 1111111111,
				EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:    testMarketOdds,

				Status: markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			betUID := uuid.NewString()

			if tc.bet != nil {
				resetMarket := markettypes.Market{
					UID:     testMarketUID,
					Creator: testCreator,
					StartTS: 1111111111,
					EndTS:   uint64(ctx.BlockTime().Unix()) + 1000,
					Odds:    testMarketOdds,
					Status:  markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
				}
				tApp.MarketKeeper.SetMarket(ctx, resetMarket)

				if resetMarket.Status == markettypes.MarketStatus_MARKET_STATUS_ACTIVE {
					_, err := tApp.OrderbookKeeper.InitiateOrderBookParticipation(
						ctx,
						simapp.TestParamUsers["user1"].Address,
						resetMarket.UID,
						sdkmath.NewInt(100000000),
						sdkmath.NewInt(1),
					)
					require.NoError(t, err)
				}

				tc.bet.UID = betUID
				placeTestBet(ctx, t, tApp, betUID, nil)
				k.SetBet(ctx, *tc.bet, 1)
			}

			if tc.updateMarket != nil {
				tApp.MarketKeeper.SetMarket(ctx, *tc.updateMarket)
			}

			if tc.betUID != "" {
				betUID = tc.betUID
			}

			if tc.bet == nil {
				tc.bet = &types.Bet{
					Creator: "",
				}
			}

			err := k.Settle(ctx, tc.bet.Creator, betUID)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestBatchSettleBet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)

	p := k.GetParams(ctx)
	p.BatchSettlementCount = 7
	k.SetParams(ctx, p)

	marketCount := 5
	participationCount := 3
	marketBetCount := 10
	allBetCount := marketCount * marketBetCount
	blockCount := allBetCount/int(p.BatchSettlementCount) + 1

	marketUIDs := addTestMarketBatch(t, tApp, ctx, marketCount)
	for _, marketUID := range marketUIDs {
		market, found := tApp.MarketKeeper.GetMarket(ctx, marketUID)
		require.True(t, found)

		market.Status = markettypes.MarketStatus_MARKET_STATUS_ACTIVE
		tApp.MarketKeeper.SetMarket(ctx, market)

		depositorUser := 2
		for i := depositorUser; i <= depositorUser+participationCount; i++ {
			_, err := tApp.OrderbookKeeper.InitiateOrderBookParticipation(
				ctx,
				simapp.TestParamUsers["user"+cast.ToString(i)].Address,
				market.UID,
				sdkmath.NewInt(100000000),
				sdkmath.NewInt(1),
			)
			require.NoError(t, err)
		}

		for i := 0; i < marketBetCount; i++ {
			placeTestBet(ctx, t, tApp,
				uuid.NewString(),
				&types.BetOdds{
					UID:               testOddsUID1,
					MarketUID:         marketUID,
					Value:             "4.20",
					MaxLossMultiplier: sdkmath.LegacyMustNewDecFromStr("0.1"),
				},
			)
		}
	}

	allPendingBets, err := k.GetPendingBets(ctx)
	require.NoError(t, err)
	require.Equal(t, allBetCount, len(allPendingBets))

	for _, marketUID := range marketUIDs[:len(marketUIDs)-2] {
		market, found := tApp.MarketKeeper.GetMarket(ctx, marketUID)
		require.True(t, found)
		tApp.MarketKeeper.Resolve(ctx, market, &markettypes.MarketResolutionTicketPayload{
			UID:            marketUID,
			ResolutionTS:   uint64(ctx.BlockTime().Unix()) + 10000,
			WinnerOddsUIDs: []string{testOddsUID1, testOddsUID2, testOddsUID3},
			Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
		})
	}

	market, found := tApp.MarketKeeper.GetMarket(ctx, marketUIDs[len(marketUIDs)-2])
	require.True(t, found)
	tApp.MarketKeeper.Resolve(ctx, market, &markettypes.MarketResolutionTicketPayload{
		UID:          marketUIDs[len(marketUIDs)-2],
		ResolutionTS: uint64(ctx.BlockTime().Unix()) + 10000,
		Status:       markettypes.MarketStatus_MARKET_STATUS_CANCELED,
	})

	market, found = tApp.MarketKeeper.GetMarket(ctx, marketUIDs[len(marketUIDs)-1])
	require.True(t, found)
	tApp.MarketKeeper.Resolve(ctx, market, &markettypes.MarketResolutionTicketPayload{
		UID:          marketUIDs[len(marketUIDs)-1],
		ResolutionTS: uint64(ctx.BlockTime().Unix()) + 10000,
		Status:       markettypes.MarketStatus_MARKET_STATUS_ABORTED,
	})

	for i := 1; i <= blockCount; i++ {
		ctx = ctx.WithBlockHeight(int64(i))
		err := k.BatchMarketSettlements(ctx)
		require.NoError(t, err)

		pendingBets, err := k.GetPendingBets(ctx)
		require.NoError(t, err)

		settledBets, err := k.GetSettledBets(ctx)
		require.NoError(t, err)

		marketStats := tApp.MarketKeeper.GetMarketStats(ctx)

		t.Logf(
			"block: %d, pending bets: %d, settled bets: %d, resolved markets: %v\n",
			i,
			len(pendingBets),
			len(settledBets),
			marketStats.ResolvedUnsettled,
		)
		require.GreaterOrEqual(t, int(p.BatchSettlementCount)*i, len(settledBets))
	}

	allPendingBets, err = k.GetPendingBets(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(allPendingBets))

	allSettledBets, err := k.GetSettledBets(ctx)
	require.NoError(t, err)
	require.Equal(t, allBetCount, len(allSettledBets))

	allBets, err := k.GetBets(ctx)
	require.NoError(t, err)
	for _, bet := range allBets {
		require.NotEqual(t, 0, bet.SettlementHeight)
	}
}
