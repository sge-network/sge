package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/orderbook/keeper"
	"github.com/stretchr/testify/require"
)

func setupKeeperAndApp(t testing.TB) (*simappUtil.TestApp, *keeper.KeeperTest, sdk.Context) {
	tApp, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	return tApp, &tApp.OrderbookKeeper, ctx.WithBlockTime(time.Now())
}
