package keeper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/testutil/nullify"
	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/sge-network/sge/x/orderbook/types"
)

func createNParticipationBetPair(
	keeper *keeper.KeeperTest,
	ctx sdk.Context,
	n int,
) []types.ParticipationBetPair {
	items := make([]types.ParticipationBetPair, n)

	for i := range items {
		items[i].BetUID = uuid.NewString()
		items[i].OrderBookUID = testOrderBookUID
		items[i].ParticipationIndex = testParticipationIndex

		keeper.SetParticipationBetPair(ctx, items[i], cast.ToUint64(i+10))
	}
	return items
}

func TestParticipationBetPairGetAll(t *testing.T) {
	k, ctx := setupKeeper(t)
	items := createNParticipationBetPair(k, ctx, 10)

	betPairs, err := k.GetAllParticipationBetPair(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(betPairs),
	)
}
