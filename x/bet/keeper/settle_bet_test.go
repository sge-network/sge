package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/google/uuid"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/require"
)

func TestSettleBet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	testCreator = simappUtil.TestParamUsers["user1"].Address.String()
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
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			desc: "failed in checking status",
			bet: &types.Bet{
				MarketUID: testMarketUID,
				OddsValue: "10",
				OddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
				Amount:    sdk.NewInt(500),
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
				OddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
				Amount:    sdk.NewInt(500),
				Creator:   testCreator,
				OddsUID:   testOddsUID1,
			},
			updateMarket: &markettypes.Market{
				UID:                    testMarketUID,
				Creator:                testCreator,
				StartTS:                1111111111,
				EndTS:                  uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:                   testMarketOdds,
				SrContributionForHouse: sdk.NewInt(500000),

				Status: markettypes.MarketStatus_MARKET_STATUS_ABORTED,
			},
		},
		{
			desc: "market is canceled",
			bet: &types.Bet{
				MarketUID: testMarketUID,
				OddsValue: "10",
				OddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
				Amount:    sdk.NewInt(300),
				Creator:   testCreator,
				OddsUID:   testOddsUID1,
			},
			updateMarket: &markettypes.Market{
				UID:                    testMarketUID,
				Creator:                testCreator,
				StartTS:                1111111111,
				EndTS:                  uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:                   testMarketOdds,
				SrContributionForHouse: sdk.NewInt(500000),

				Status: markettypes.MarketStatus_MARKET_STATUS_CANCELED,
			},
		},
		{
			desc: "result is not declared",
			bet: &types.Bet{
				MarketUID: testMarketUID,
				OddsValue: "10",
				OddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
				Amount:    sdk.NewInt(500),
				Creator:   testCreator,
				OddsUID:   testOddsUID1,
			},
			updateMarket: &markettypes.Market{
				UID:                    testMarketUID,
				Creator:                testCreator,
				StartTS:                1111111111,
				EndTS:                  uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:                   testMarketOdds,
				SrContributionForHouse: sdk.NewInt(500000),

				Status: markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
			},
			err: types.ErrResultNotDeclared,
		},
		{
			desc: "success",
			bet: &types.Bet{
				MarketUID: testMarketUID,
				OddsValue: "10",
				OddsType:  types.OddsType_ODDS_TYPE_DECIMAL,
				Amount:    sdk.NewInt(500),
				Creator:   testCreator,
				OddsUID:   testOddsUID1,

				Result: types.Bet_RESULT_WON,
			},
			updateMarket: &markettypes.Market{
				UID:                    testMarketUID,
				Creator:                testCreator,
				StartTS:                1111111111,
				EndTS:                  uint64(ctx.BlockTime().Unix()) + 1000,
				Odds:                   testMarketOdds,
				SrContributionForHouse: sdk.NewInt(500000),

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
					BetConstraints: &markettypes.MarketBetConstraints{
						MinAmount: sdk.NewInt(1),
						BetFee:    sdk.NewInt(1),
					},
					SrContributionForHouse: sdk.NewInt(500000),
				}
				tApp.MarketKeeper.SetMarket(ctx, resetMarket)
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

			err := k.SettleBet(ctx, tc.bet.Creator, betUID)
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
	marketBetCount := 10
	allBetCount := marketCount * marketBetCount
	blockCount := allBetCount/int(p.BatchSettlementCount) + 1

	marketUIDs := addTestMarketBatch(t, tApp, ctx, marketCount)
	for _, marketUID := range marketUIDs {
		market, found := tApp.MarketKeeper.GetMarket(ctx, marketUID)
		require.True(t, found)

		market.Status = markettypes.MarketStatus_MARKET_STATUS_ACTIVE
		market.BetConstraints = &markettypes.MarketBetConstraints{
			MinAmount: sdk.NewInt(1),
			BetFee:    sdk.NewInt(1),
		}
		tApp.MarketKeeper.SetMarket(ctx, market)

		for i := 0; i < marketBetCount; i++ {
			placeTestBet(ctx, t, tApp,
				uuid.NewString(),
				&types.BetOdds{
					UID:               testOddsUID1,
					MarketUID:         marketUID,
					Value:             "4.20",
					MaxLossMultiplier: sdk.MustNewDecFromStr("0.1"),
				},
			)
		}
	}

	allPendingBets, err := k.GetPendingBets(ctx)
	require.NoError(t, err)
	require.Equal(t, allBetCount, len(allPendingBets))

	for _, marketUID := range marketUIDs[:len(marketUIDs)-2] {
		_, err := tApp.MarketKeeper.ResolveMarket(ctx, &markettypes.MarketResolutionTicketPayload{
			UID:            marketUID,
			ResolutionTS:   uint64(ctx.BlockTime().Unix()) + 10000,
			WinnerOddsUIDs: []string{testOddsUID1, testOddsUID2, testOddsUID3},
			Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
		})
		require.NoError(t, err)
	}

	_, err = tApp.MarketKeeper.ResolveMarket(ctx, &markettypes.MarketResolutionTicketPayload{
		UID:          marketUIDs[len(marketUIDs)-2],
		ResolutionTS: uint64(ctx.BlockTime().Unix()) + 10000,
		Status:       markettypes.MarketStatus_MARKET_STATUS_CANCELED,
	})
	require.NoError(t, err)

	_, err = tApp.MarketKeeper.ResolveMarket(ctx, &markettypes.MarketResolutionTicketPayload{
		UID:          marketUIDs[len(marketUIDs)-1],
		ResolutionTS: uint64(ctx.BlockTime().Unix()) + 10000,
		Status:       markettypes.MarketStatus_MARKET_STATUS_ABORTED,
	})
	require.NoError(t, err)

	for i := 1; i <= blockCount; i++ {
		ctx = ctx.WithBlockHeight(int64(i))
		err := k.BatchMarketSettlements(ctx)
		require.NoError(t, err)

		pendingBets, err := k.GetPendingBets(ctx)
		require.NoError(t, err)

		settledBets, err := k.GetSettledBets(ctx)
		require.NoError(t, err)

		marketStats := tApp.MarketKeeper.GetMarketStats(ctx)

		t.Logf("block: %d, pending bets: %d, settled bets: %d, resolved markets: %v\n", i, len(pendingBets), len(settledBets), marketStats.ResolvedUnsettled)
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
				Status: types.Bet_STATUS_CANCELED,
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

func TestProcessBetResultAndStatus(t *testing.T) {
	k, _ := setupKeeper(t)
	tcs := []struct {
		desc   string
		bet    *types.Bet
		market markettypes.Market
		err    error
		result types.Bet_Result
	}{
		{
			desc: "not declared",
			market: markettypes.Market{
				Status: markettypes.MarketStatus_MARKET_STATUS_UNSPECIFIED,
			},
			bet: &types.Bet{},
			err: types.ErrResultNotDeclared,
		},
		{
			desc: "won",
			market: markettypes.Market{
				Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
				WinnerOddsUIDs: []string{"oddsUID"},
			},
			bet: &types.Bet{
				OddsUID: "oddsUID",
			},
			result: types.Bet_RESULT_WON,
		},
		{
			desc: "lost",
			market: markettypes.Market{
				Status:         markettypes.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
				WinnerOddsUIDs: []string{"oddsUID"},
			},
			bet:    &types.Bet{},
			result: types.Bet_RESULT_LOST,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			err := k.ProcessBetResultAndStatus(tc.bet, tc.market)
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
