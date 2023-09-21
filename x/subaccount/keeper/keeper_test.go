package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/stretchr/testify/require"

	"github.com/sge-network/sge/x/subaccount/keeper"
)

func setupKeeperAndApp(t testing.TB) (*simapp.TestApp, *keeper.Keeper, sdk.Context) {
	tApp, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)

	return tApp, tApp.SubaccountKeeper, ctx.WithBlockTime(time.Now())
}
