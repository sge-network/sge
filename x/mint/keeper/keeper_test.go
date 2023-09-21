package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/simapp"
	"github.com/sge-network/sge/x/mint/keeper"
	"github.com/stretchr/testify/require"
)

func setupKeeperAndApp(t testing.TB) (*simapp.TestApp, *keeper.KeeperTest, sdk.Context) {
	tApp, ctx, err := simapp.GetTestObjects()
	require.NoError(t, err)

	return tApp, &tApp.MintKeeper, ctx
}

func setupKeeper(t testing.TB) (*keeper.KeeperTest, sdk.Context) {
	_, k, ctx := setupKeeperAndApp(t)

	return k, ctx
}

func TestLogger(t *testing.T) {
	k, ctx := setupKeeper(t)
	logger := k.Logger(ctx)
	require.NotNil(t, logger)

	logger.Debug("test log")
}
