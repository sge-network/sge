package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/nullify"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func createNOrderBook(
	tApp *simappUtil.TestApp,
	keeper *keeper.KeeperTest,
	ctx sdk.Context,
	n int,
) []types.OrderBook {
	items := make([]types.OrderBook, n)

	for i := range items {
		items[i].OddsCount = cast.ToUint64(i + 1)
		items[i].ParticipationCount = cast.ToUint64(i + 10)
		items[i].Status = types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE
		items[i].UID = uuid.NewString()

		keeper.SetOrderBook(ctx, items[i])
	}
	return items
}

func TestOrderBookGet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNOrderBook(tApp, k, ctx, 10)

	rst, found := k.GetOrderBook(ctx,
		uuid.NewString(),
	)
	var expectedResp types.OrderBook
	require.False(t, found)
	require.Equal(t,
		nullify.Fill(expectedResp),
		nullify.Fill(rst),
	)

	for _, item := range items {
		rst, found := k.GetOrderBook(ctx,
			item.UID,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(item),
			nullify.Fill(rst),
		)
	}
}

func TestOrderBookGetAll(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNOrderBook(tApp, k, ctx, 10)

	orderBooks, err := k.GetAllOrderBooks(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(orderBooks),
	)
}
