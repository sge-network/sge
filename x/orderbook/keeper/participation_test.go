package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/testutil/sample"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	housetypes "github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func createNParticipation(
	tApp *simappUtil.TestApp,
	keeper *keeper.KeeperTest,
	ctx sdk.Context,
	n int,
) []types.OrderBookParticipation {
	items := make([]types.OrderBookParticipation, n)

	for i := range items {
		items[i].Index = cast.ToUint64(i + 1)
		items[i].ParticipantAddress = simappUtil.TestParamUsers["user1"].Address.String()
		items[i].OrderBookUID = testOrderBookUID
		items[i].ActualProfit = sdk.NewInt(100)
		items[i].CurrentRoundLiquidity = sdk.NewInt(100)
		items[i].CurrentRoundMaxLoss = sdk.NewInt(100)
		items[i].CurrentRoundTotalBetAmount = sdk.NewInt(100)
		items[i].Liquidity = sdk.NewInt(100)
		items[i].Fee = sdk.NewInt(10)
		items[i].MaxLoss = sdk.NewInt(100)
		items[i].TotalBetAmount = sdk.NewInt(100)

		keeper.SetOrderBookParticipation(ctx, items[i])
	}
	return items
}

func createTestMarket(
	tApp *simappUtil.TestApp,
	keeper *keeper.KeeperTest,
	ctx sdk.Context,
	marketUID string,
	status markettypes.MarketStatus,
	oddsUIDs []string,
) markettypes.Market {
	market := markettypes.NewMarket(marketUID,
		uuid.NewString(),
		cast.ToUint64(ctx.BlockTime().Unix()),
		cast.ToUint64(ctx.BlockTime().Add(5*time.Minute).Unix()),
		[]*markettypes.Odds{
			{UID: oddsUIDs[0], Meta: "odds1"},
			{UID: oddsUIDs[1], Meta: "odds2"},
		},
		"test market",
		marketUID,
		status,
	)
	tApp.MarketKeeper.SetMarket(ctx, market)
	return market
}

func TestParticipationGet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNParticipation(tApp, k, ctx, 10)

	rst, found := k.GetOrderBookParticipation(ctx,
		items[0].OrderBookUID,
		10000,
	)
	var expectedResp types.OrderBookParticipation
	require.False(t, found)
	require.Equal(t,
		nullify.Fill(expectedResp),
		nullify.Fill(rst),
	)

	for i, item := range items {
		rst, found := k.GetOrderBookParticipation(ctx,
			items[i].OrderBookUID,
			uint64(i+1),
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(item),
			nullify.Fill(rst),
		)
	}
}

func TestParticipationsOfOrderBookGet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNParticipation(tApp, k, ctx, 10)

	rst, err := k.GetParticipationsOfOrderBook(ctx,
		uuid.NewString(),
	)
	var expectedResp []types.OrderBookParticipation
	require.NoError(t, err)
	require.Equal(t,
		nullify.Fill(expectedResp),
		nullify.Fill(rst),
	)

	rst, err = k.GetParticipationsOfOrderBook(ctx,
		testOrderBookUID,
	)
	require.NoError(t, err)
	require.Equal(t,
		nullify.Fill(items),
		nullify.Fill(rst),
	)
}

func TestParticipationGetAll(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNParticipation(tApp, k, ctx, 10)

	participations, err := k.GetAllOrderBookParticipations(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(participations),
	)
}

func TestInitiateOrderBookParticipationNoMarket(t *testing.T) {
	k, ctx := setupKeeper(t)

	_, err := k.InitiateOrderBookParticipation(ctx,
		simappUtil.TestParamUsers["user1"].Address,
		testOrderBookUID,
		sdk.NewInt(1000),
		sdk.NewInt(100))
	require.ErrorIs(t, types.ErrMarketNotFound, err)
}

func TestInitiateOrderBookParticipation(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)

	oddsUIDs := []string{uuid.NewString(), uuid.NewString()}

	for _, tc := range []struct {
		desc          string
		depositorAddr sdk.AccAddress
		marketStatus  markettypes.MarketStatus

		err error
	}{
		{
			desc:          "not active market",
			depositorAddr: simappUtil.TestParamUsers["user1"].Address,
			marketStatus:  markettypes.MarketStatus_MARKET_STATUS_INACTIVE,
			err:           types.ErrParticipationOnInactiveMarket,
		},
		{
			desc:          "not enough fund",
			depositorAddr: sdk.MustAccAddressFromBech32(sample.AccAddress()),
			marketStatus:  markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
			err:           types.ErrInsufficientAccountBalance,
		},
		{
			desc:          "success",
			depositorAddr: simappUtil.TestParamUsers["user1"].Address,
			marketStatus:  markettypes.MarketStatus_MARKET_STATUS_ACTIVE,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			marketUID := uuid.NewString()
			createTestMarket(tApp, k, ctx, marketUID, tc.marketStatus, oddsUIDs)

			err := k.InitiateOrderBook(ctx, marketUID, oddsUIDs)
			require.NoError(t, err)

			participationIndex, err := k.InitiateOrderBookParticipation(ctx,
				tc.depositorAddr,
				marketUID,
				sdk.NewInt(1000),
				sdk.NewInt(100))

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testParticipationIndex, participationIndex)
			}
		})
	}
}

func TestWithdrawOrderBookParticipation(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)

	oddsUIDs := []string{uuid.NewString(), uuid.NewString()}
	fee := sdk.NewInt(100)

	for _, tc := range []struct {
		desc           string
		depositorAddr  sdk.AccAddress
		depositAmount  sdk.Int
		withdrawMode   housetypes.WithdrawalMode
		withdrawAmount sdk.Int

		expWithdrawnAmount sdk.Int
		err                error
	}{
		{
			desc:          "no participation",
			depositorAddr: simappUtil.TestParamUsers["user1"].Address,
			depositAmount: sdk.ZeroInt(),
			withdrawMode:  housetypes.WithdrawalMode_WITHDRAWAL_MODE_FULL,
			err:           types.ErrOrderBookParticipationNotFound,
		},
		{
			desc:           "withdraw amount too large",
			depositorAddr:  simappUtil.TestParamUsers["user1"].Address,
			depositAmount:  sdk.NewInt(1000),
			withdrawMode:   housetypes.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL,
			withdrawAmount: sdk.NewInt(10000),
			err:            types.ErrWithdrawalTooLarge,
		},
		{
			desc:               "success full",
			depositorAddr:      simappUtil.TestParamUsers["user1"].Address,
			depositAmount:      sdk.NewInt(1000),
			withdrawMode:       housetypes.WithdrawalMode_WITHDRAWAL_MODE_FULL,
			expWithdrawnAmount: sdk.NewInt(1000).Sub(fee),
		},
		{
			desc:               "success partial",
			depositorAddr:      simappUtil.TestParamUsers["user1"].Address,
			depositAmount:      sdk.NewInt(500),
			withdrawMode:       housetypes.WithdrawalMode_WITHDRAWAL_MODE_FULL,
			expWithdrawnAmount: sdk.NewInt(500).Sub(fee),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			marketUID := uuid.NewString()
			createTestMarket(tApp, k, ctx, marketUID, markettypes.MarketStatus_MARKET_STATUS_ACTIVE, oddsUIDs)

			err := k.InitiateOrderBook(ctx, marketUID, oddsUIDs)
			require.NoError(t, err)

			var participationIndex uint64
			if !tc.depositAmount.IsZero() {
				participationIndex, err = k.InitiateOrderBookParticipation(ctx,
					tc.depositorAddr,
					marketUID,
					tc.depositAmount,
					sdk.NewInt(100))
				require.NoError(t, err)
			}

			withdrawnAmount, err := k.CalcWithdrawalAmount(ctx,
				tc.depositorAddr.String(),
				marketUID,
				participationIndex,
				tc.withdrawMode,
				sdk.ZeroInt(),
				tc.depositAmount,
			)

			if err == nil {
				err = k.WithdrawOrderBookParticipation(ctx,
					marketUID,
					participationIndex,
					withdrawnAmount,
				)
			}

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expWithdrawnAmount, withdrawnAmount)
			}
		})
	}
}
