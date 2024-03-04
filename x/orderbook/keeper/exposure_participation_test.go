package keeper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/sge-network/sge/x/orderbook/types"
)

func createNParticipationExposure(
	keeper *keeper.KeeperTest,
	ctx sdk.Context,
	n int,
) []types.ParticipationExposure {
	items := make([]types.ParticipationExposure, n)

	for i := range items {
		items[i].ParticipationIndex = testParticipationIndex
		items[i].OrderBookUID = testOrderBookUID
		items[i].OddsUID = uuid.NewString()
		items[i].Round = 1
		items[i].IsFulfilled = false
		items[i].BetAmount = sdkmath.NewInt(100)
		items[i].Exposure = sdkmath.NewInt(100)

		keeper.SetParticipationExposure(ctx, items[i])
		keeper.SetParticipationExposureByIndex(ctx, items[i])
		keeper.SetHistoricalParticipationExposure(ctx, items[i])
	}
	return items
}

func TestParticipationExposureGet(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNParticipationExposure(k, ctx, 10)

	rst, err := k.GetExposureByOrderBookAndOdds(ctx,
		items[0].OrderBookUID,
		uuid.NewString(),
	)
	var expectedResp []types.ParticipationExposure
	require.NoError(t, err)
	require.Equal(t,
		nullify.Fill(expectedResp),
		nullify.Fill(rst),
	)

	for _, item := range items {
		rst, err := k.GetExposureByOrderBookAndOdds(ctx,
			item.OrderBookUID,
			item.OddsUID,
		)
		require.NoError(t, err)
		require.Equal(t,
			nullify.Fill([]types.ParticipationExposure{item}),
			nullify.Fill(rst),
		)
	}
}

func TestExposureByOrderBookAndParticipationIndexGet(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNParticipationExposure(k, ctx, 10)

	rst, err := k.GetExposureByOrderBookAndParticipationIndex(ctx,
		items[0].OrderBookUID,
		1000,
	)
	var expectedResp []types.ParticipationExposure
	require.NoError(t, err)
	require.Equal(t,
		nullify.Fill(expectedResp),
		nullify.Fill(rst),
	)

	rst, err = k.GetExposureByOrderBookAndParticipationIndex(ctx,
		testOrderBookUID,
		testParticipationIndex,
	)
	require.NoError(t, err)
	require.Equal(t, len(items), len(rst))
}

func TestAllHistoricalParticipationExposuresGet(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNParticipationExposure(k, ctx, 10)

	rst, err := k.GetAllHistoricalParticipationExposures(ctx)
	require.NoError(t, err)
	require.Equal(t, len(items), len(rst))
}

func TestParticipationExposureGetAll(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNParticipationExposure(k, ctx, 10)

	exposures, err := k.GetAllParticipationExposures(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(exposures),
	)
}
