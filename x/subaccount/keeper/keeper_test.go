package keeper_test

import (
	"testing"
	"time"

	"github.com/sge-network/sge/x/subaccount/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/stretchr/testify/require"
)

func setupKeeperAndApp(t testing.TB) (*simappUtil.TestApp, keeper.Keeper, sdk.Context) {
	tApp, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	return tApp, tApp.SubaccountKeeper, ctx.WithBlockTime(time.Now())
}
