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

func createNParticipationExposure(
	tApp *simappUtil.TestApp,
	keeper *keeper.KeeperTest,
	ctx sdk.Context,
	n int,
) []types.ParticipationExposure {
	items := make([]types.ParticipationExposure, n)

	for i := range items {
		items[i].ParticipationIndex = cast.ToUint64(i + 1)
		items[i].OrderBookUID = uuid.NewString()
		items[i].OddsUID = uuid.NewString()
		items[i].Round = 1
		items[i].IsFulfilled = false
		items[i].BetAmount = sdk.NewInt(100)
		items[i].Exposure = sdk.NewInt(100)

		keeper.SetParticipationExposure(ctx, items[i])
	}
	return items
}

func TestParticipationExposureGet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNParticipationExposure(tApp, k, ctx, 10)

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

func TestParticipationExposureGetAll(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNParticipationExposure(tApp, k, ctx, 10)

	exposures, err := k.GetAllParticipationExposures(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(exposures),
	)
}
