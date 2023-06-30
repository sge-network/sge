package keeper_test

import (
	"github.com/sge-network/sge/x/subaccount/keeper"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simappUtil "github.com/sge-network/sge/testutil/simapp"
	"github.com/stretchr/testify/require"
)

func setupKeeperAndApp(t testing.TB) (*simappUtil.TestApp, keeper.Keeper, sdk.Context) {
	tApp, ctx, err := simappUtil.GetTestObjects()
	require.NoError(t, err)

	return tApp, tApp.SubAccountKeeper, ctx.WithBlockTime(time.Now())
}

func TestSubaccountID(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	// Peek from beginning should return 0
	require.Equal(t, uint64(0), k.Peek(ctx))

	// NextID returns the actual ID and increments the ID
	require.Equal(t, uint64(0), k.NextID(ctx))
	require.Equal(t, uint64(1), k.Peek(ctx))

	// We can set an arbitrary ID and continue from there
	k.SetID(ctx, 100)
	require.Equal(t, uint64(100), k.Peek(ctx))
}
