package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/require"
)

func TestUnlockedBalanceQuery(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	subAccAddr := sample.NativeAccAddress()
	blockTime := time.Now()
	ctx = ctx.WithBlockTime(blockTime)
	k.SetLockedBalances(ctx, subAccAddr, []types.LockedBalance{
		{
			UnlockTS: cast.ToUint64(blockTime.Unix()),
			Amount:   sdkmath.NewInt(100),
		},
		{
			UnlockTS: cast.ToUint64(blockTime.Add(1 * time.Second).Unix()),
			Amount:   sdkmath.NewInt(200),
		},
		{
			UnlockTS: cast.ToUint64(blockTime.Add(-1 * time.Second).Unix()),
			Amount:   sdkmath.NewInt(1000),
		},
	})

	locked, _ := k.GetBalances(ctx, subAccAddr, types.LockedBalanceStatus_LOCKED_BALANCE_STATUS_LOCKED)
	require.Equal(t, []types.LockedBalance{
		{
			UnlockTS: cast.ToUint64(blockTime.Unix()),
			Amount:   sdkmath.NewInt(100),
		},
		{
			UnlockTS: cast.ToUint64(blockTime.Add(1 * time.Second).Unix()),
			Amount:   sdkmath.NewInt(200),
		},
	}, locked)

	_, unlockedAmount := k.GetBalances(ctx, subAccAddr, types.LockedBalanceStatus_LOCKED_BALANCE_STATUS_UNLOCKED)
	require.Equal(t, int64(1000), unlockedAmount.Int64())
}
