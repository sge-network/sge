package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/market/keeper"
	"github.com/sge-network/sge/x/market/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func createNMarket(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Market {
	items := make([]types.Market, n)
	for i := range items {
		items[i].UID = cast.ToString(i)
		items[i].SrContributionForHouse = sdk.NewInt(0)

		keeper.SetMarket(ctx, items[i])
	}
	return items
}

func TestMarketGet(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNMarket(k, ctx, 10)
	_, found := k.GetMarket(ctx,
		"NotExistUid",
	)
	require.False(t, found)

	for _, item := range items {
		rst, found := k.GetMarket(ctx,
			item.UID,
		)
		require.True(t, found)
		require.EqualValues(t, item, rst)
	}
}

func TestGetMarket(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNMarket(k, ctx, 10)

	rst, found := k.GetMarket(ctx, "NotExistUid")
	require.False(t, found)
	require.Equal(t, rst.UID, "")

	for _, item := range items {
		rst, found := k.GetMarket(ctx,
			item.UID,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(item),
			nullify.Fill(rst),
		)
	}
}

func TestMarketRemove(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNMarket(k, ctx, 10)
	for _, item := range items {
		k.RemoveMarket(ctx,
			item.UID,
		)
		_, found := k.GetMarket(ctx,
			item.UID,
		)
		require.False(t, found)
	}
}

func TestMarketGetAll(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNMarket(k, ctx, 10)

	markets, err := k.GetMarketAll(ctx)
	require.NoError(t, err)

	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(markets),
	)
}

func TestResolveMarkets(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		k, ctx := setupKeeper(t)
		resEventsIn := types.MarketResolutionTicketPayload{
			UID: "NotExistUid",
		}
		_, err := k.ResolveMarket(ctx, &resEventsIn)
		require.Equal(t, types.ErrNoMatchingMarket, err)
	})

	t.Run("NotPending", func(t *testing.T) {
		k, ctx := setupKeeper(t)

		item := types.Market{
			UID:    "uid",
			Status: types.MarketStatus_MARKET_STATUS_CANCELED,
		}
		k.SetMarket(ctx, item)

		resEventsIn := types.MarketResolutionTicketPayload{
			UID: item.UID,
		}

		_, err := k.ResolveMarket(ctx, &resEventsIn)
		require.Equal(t, types.ErrCanNotBeAltered, err)
	})

	t.Run("Success", func(t *testing.T) {
		k, ctx := setupKeeper(t)

		item := types.Market{
			UID:    "uid",
			Status: types.MarketStatus_MARKET_STATUS_ACTIVE,
		}
		k.SetMarket(ctx, item)

		resEventsIn := types.MarketResolutionTicketPayload{
			UID:            item.UID,
			ResolutionTS:   123456,
			WinnerOddsUIDs: []string{"oddsUID1", "oddsUID2"},
			Status:         types.MarketStatus_MARKET_STATUS_RESULT_DECLARED,
		}
		_, err := k.ResolveMarket(ctx, &resEventsIn)
		require.Nil(t, err)
		val, found := k.GetMarket(ctx, item.UID)
		require.True(t, found)
		require.Equal(t, resEventsIn.ResolutionTS, val.ResolutionTS)
		require.Equal(t, resEventsIn.WinnerOddsUIDs, val.WinnerOddsUIDs)
		require.Equal(t, resEventsIn.Status, val.Status)
	})
}

func TestMarketExist(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		k, ctx := setupKeeper(t)
		found := k.MarketExists(ctx, "notExistMarketUID")
		require.False(t, found)
	})

	t.Run("Found", func(t *testing.T) {
		k, ctx := setupKeeper(t)
		item := types.Market{
			UID: "uid",
		}
		k.SetMarket(ctx, item)
		found := k.MarketExists(ctx, item.UID)
		require.True(t, found)
	})
}
