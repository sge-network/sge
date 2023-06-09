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

func createNParticipationBetPair(
	tApp *simappUtil.TestApp,
	keeper *keeper.KeeperTest,
	ctx sdk.Context,
	n int,
) []types.ParticipationBetPair {
	items := make([]types.ParticipationBetPair, n)

	for i := range items {
		items[i].BetUID = uuid.NewString()
		items[i].OrderBookUID = uuid.NewString()
		items[i].ParticipationIndex = cast.ToUint64(i + 1)

		keeper.SetParticipationBetPair(ctx, items[i], cast.ToUint64(i+10))
	}
	return items
}

func TestParticipationBetPairGetAll(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNParticipationBetPair(tApp, k, ctx, 10)

	betPairs, err := k.GetAllParticipationBetPair(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(betPairs),
	)
}
