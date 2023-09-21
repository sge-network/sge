package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/stretchr/testify/require"
)

var (
	testOrderBookUID       = uuid.NewString()
	testParticipationIndex = uint64(1)
)

func setupKeeper(t testing.TB) (*keeper.KeeperTest, sdk.Context) {
	_, k, ctx := setupKeeperAndApp(t)

	return k, ctx
}

func setupKeeperAndApp(t testing.TB) (*simapp.TestApp, *keeper.KeeperTest, sdk.Context) {
	tApp, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)

	return tApp, tApp.OrderbookKeeper, ctx.WithBlockTime(time.Now())
}
