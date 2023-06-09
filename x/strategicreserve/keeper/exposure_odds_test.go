package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/nullify"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/strategicreserve/keeper"
	"github.com/sge-network/sge/x/strategicreserve/types"
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
		items[i].OrderBookUID = uuid.NewString()
		items[i].ActualProfit = sdk.NewInt(100)
		items[i].CurrentRoundLiquidity = sdk.NewInt(100)
		items[i].CurrentRoundMaxLoss = sdk.NewInt(100)
		items[i].CurrentRoundTotalBetAmount = sdk.NewInt(100)
		items[i].Liquidity = sdk.NewInt(100)
		items[i].MaxLoss = sdk.NewInt(100)
		items[i].TotalBetAmount = sdk.NewInt(100)

		keeper.SetOrderBookParticipation(ctx, items[i])
	}
	return items
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
