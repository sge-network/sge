package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/testutil/nullify"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/bet/keeper"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNBet(tApp *simappUtil.TestApp, keeper *keeper.KeeperTest, ctx sdk.Context, n int) []types.Bet {
	items := make([]types.Bet, n)
	testCreator = simappUtil.TestParamUsers["user1"].Address.String()
	tApp.SporteventKeeper.SetSportEvent(ctx, testSportEvent)

	for i := range items {
		items[i].UID = strconv.Itoa(i)
		items[i].Creator = testCreator
		items[i].OddsValue = "10"
		items[i].OddsType = types.OddsType_ODDS_TYPE_DECIMAL
		items[i].Amount = sdk.NewInt(10)
		items[i].BetFee = sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(1))
		items[i].SportEventUID = testSportEventUID

		id := uint64(i + 1)
		keeper.SetBet(ctx, items[i], id)
		keeper.SetActiveBet(ctx, &types.ActiveBet{
			ID:      id,
			Creator: testCreator,
		}, testSportEventUID)
	}
	return items
}

func TestBetGet(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNBet(tApp, k, ctx, 10)
	testCreator = simappUtil.TestParamUsers["user1"].Address.String()

	rst, found := k.GetBet(ctx,
		testCreator,
		10000,
	)
	var expectedResp types.Bet
	require.False(t, found)
	require.Equal(t,
		nullify.Fill(expectedResp),
		nullify.Fill(rst),
	)

	for i, item := range items {
		rst, found := k.GetBet(ctx,
			testCreator,
			uint64(i+1),
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(item),
			nullify.Fill(rst),
		)
	}
}

func TestBetGetAll(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNBet(tApp, k, ctx, 10)

	bets, err := k.GetBets(ctx)
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(bets),
	)
}

// TestSortBetGetAll checks if incremental id is genereted correctly
func TestSortBetGetAll(t *testing.T) {
	tApp, k, ctx := setupKeeperAndApp(t)
	items := createNBet(tApp, k, ctx, 10000)

	bets, err := k.GetBets(ctx)
	lastBetID := uint64(0)
	for _, b := range bets {
		uuid2ID, found := k.GetBetID(ctx, b.UID)
		require.True(t, found)
		require.Greater(t, uuid2ID.ID, lastBetID)
		lastBetID = uuid2ID.ID
	}
	require.NoError(t, err)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(bets),
	)
}
