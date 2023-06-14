package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/nullify"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/stretchr/testify/require"
)

func createNOrderBookOddsExposure(
	tApp *simappUtil.TestApp,
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

func TestOrderBookOddsExposureGet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNOrderBookOddsExposure(tApp, k, ctx, 10)

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
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNOrderBookOddsExposure(tApp, k, ctx, 10)

	exposures, err := k.GetAllOrderBookExposures(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(exposures),
	)
}
