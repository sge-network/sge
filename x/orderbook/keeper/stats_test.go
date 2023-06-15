package keeper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/stretchr/testify/require"
)

func TestStatsGet(t *testing.T) {
	k, ctx := setupKeeper(t)

	stats := types.OrderBookStats{
		ResolvedUnsettled: []string{uuid.NewString()},
	}
	k.SetOrderBookStats(ctx, stats)

	rst := k.GetOrderBookStats(ctx)
	require.Equal(t,
		nullify.Fill(stats),
		nullify.Fill(rst),
	)
}
