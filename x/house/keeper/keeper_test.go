package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/house/keeper"
	"github.com/stretchr/testify/require"
)

var (
	testMarketUID        = uuid.NewString()
	testDepositorAddress = sample.AccAddress()
)

func setupKeeper(t testing.TB) (*keeper.KeeperTest, sdk.Context) {
	_, k, ctx := setupKeeperAndApp(t)

	return k, ctx
}

func setupKeeperAndApp(t testing.TB) (*simappUtil.TestApp, *keeper.KeeperTest, sdk.Context) {
	tApp, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	return tApp, &tApp.HouseKeeper, ctx.WithBlockTime(time.Now())
}
