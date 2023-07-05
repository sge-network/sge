package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/testutil/sample"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/house/keeper"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/stretchr/testify/require"
)

var (
	testMarketUID        = uuid.NewString()
	testDepositorAddress = sample.AccAddress()
	testMarketOdds       = []*markettypes.Odds{
		{UID: uuid.NewString(), Meta: "Odds 1"},
		{UID: uuid.NewString(), Meta: "Odds 2"},
		{UID: uuid.NewString(), Meta: "Odds 3"},
	}
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
