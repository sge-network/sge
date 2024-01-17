package keeper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/sge-network/sge/x/orderbook/types"
)

func createNOrderBookOddsExposure(
	keeper *keeper.KeeperTest,
	ctx sdk.Context,
	n int,
) []types.OrderBookOddsExposure {
	items := make([]types.OrderBookOddsExposure, n)

	for i := range items {
		items[i].FulfillmentQueue = []uint64{1}
		items[i].OddsUID = uuid.NewString()
		items[i].OrderBookUID = testOrderBookUID

		keeper.SetOrderBookOddsExposure(ctx, items[i])
	}
	return items
}

func TestOddsExposuresByOrderBookGet(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNOrderBookOddsExposure(k, ctx, 10)

	rst, err := k.GetOddsExposuresByOrderBook(ctx,
		uuid.NewString(),
	)
	var expectedResp []types.OrderBookOddsExposure
	require.NoError(t, err)
	require.Equal(t,
		nullify.Fill(expectedResp),
		nullify.Fill(rst),
	)

	rst, err = k.GetOddsExposuresByOrderBook(ctx,
		testOrderBookUID,
	)

	require.NoError(t, err)
	require.Equal(t, len(items), len(rst))
}

func TestOrderBookOddsExposureGet(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNOrderBookOddsExposure(k, ctx, 10)

	rst, found := k.GetOrderBookOddsExposure(ctx,
		uuid.NewString(),
		uuid.NewString(),
	)
	var expectedResp types.OrderBookOddsExposure
	require.False(t, found)
	require.Equal(t,
		nullify.Fill(expectedResp),
		nullify.Fill(rst),
	)

	for _, item := range items {
		rst, found := k.GetOrderBookOddsExposure(ctx,
			item.OrderBookUID,
			item.OddsUID,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(item),
			nullify.Fill(rst),
		)
	}
}

func TestOrderBookOddsExposureGetAll(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNOrderBookOddsExposure(k, ctx, 10)

	exposures, err := k.GetAllOrderBookExposures(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(exposures),
	)
}
