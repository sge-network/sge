package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/ovm/keeper"
)

func setupKeeperAndApp(t testing.TB) (*simapp.TestApp, *keeper.KeeperTest, sdk.Context) {
	tApp, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)

	return tApp, tApp.OVMKeeper, ctx
}

func setupKeeper(t testing.TB) (*keeper.KeeperTest, sdk.Context) {
	_, k, ctx := setupKeeperAndApp(t)

	return k, ctx
}

func TestLogger(t *testing.T) {
	k, ctx := setupKeeper(t)
	logger := k.Logger(ctx)
	require.True(t, logger != nil)
}
