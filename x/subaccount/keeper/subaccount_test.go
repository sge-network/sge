package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/types"
)

func TestSubaccountID(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	// Peek from beginning should return 1
	require.Equal(t, uint64(1), k.Peek(ctx))

	// NextID returns the actual ID and increments the ID
	require.Equal(t, uint64(1), k.NextID(ctx))
	require.Equal(t, uint64(2), k.Peek(ctx))

	// We can set an arbitrary ID and continue from there
	k.SetID(ctx, 100)
	require.Equal(t, uint64(100), k.Peek(ctx))
}

func TestSubaccountOwner(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	owner := sample.NativeAccAddress()

	// Account should not have subaccount
	_, exists := k.GetSubaccountOwner(ctx, owner)
	require.False(t, exists)

	// Set subaccount owner
	id := k.NextID(ctx)
	k.SetSubaccountOwner(ctx, types.NewAddressFromSubaccount(id), owner)

	// Account should have subaccount
	_, exists = k.GetSubaccountByOwner(ctx, owner)
	require.True(t, exists)

	// Get subaccount ID
	subAccountAddress, exists := k.GetSubaccountByOwner(ctx, owner)
	require.True(t, exists)
	require.Equal(t, types.NewAddressFromSubaccount(id), subAccountAddress)

	// Get owner of subaccount
	gotOwner, exists := k.GetSubaccountOwner(ctx, subAccountAddress)
	require.True(t, exists)
	require.Equal(t, owner, gotOwner)
	// Get account ID by owner
	gotSubaccount, exists := k.GetSubaccountByOwner(ctx, owner)
	require.True(t, exists)
	require.Equal(t, types.NewAddressFromSubaccount(id), gotSubaccount)
}

func TestSetLockedBalances(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	someUnlockTime := uint64(time.Now().Add(time.Hour * 24 * 365).Unix())
	otherUnlockTime := uint64(time.Now().Add(time.Hour * 24 * 365 * 2).Unix())

	balanceUnlocks := []types.LockedBalance{
		{
			Amount:   sdkmath.NewInt(10000),
			UnlockTS: someUnlockTime,
		},
		{
			Amount:   sdkmath.NewInt(20000),
			UnlockTS: otherUnlockTime,
		},
	}

	addr := types.NewAddressFromSubaccount(1)

	k.SetLockedBalances(ctx, addr, balanceUnlocks)

	// Get locked balances
	lockedBalances, _ := k.GetBalances(ctx, addr, types.LockedBalanceStatus_LOCKED_BALANCE_STATUS_LOCKED)
	for i, lockedBalance := range lockedBalances {
		require.Equal(t, lockedBalance.Amount, balanceUnlocks[i].Amount)
		require.True(t, lockedBalance.UnlockTS == balanceUnlocks[i].UnlockTS)
	}
}

func TestSetBalances(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	balance := types.AccountSummary{
		DepositedAmount: sdkmath.ZeroInt(),
		SpentAmount:     sdkmath.ZeroInt(),
		WithdrawnAmount: sdkmath.ZeroInt(),
		LostAmount:      sdk.OneInt(),
	}

	subAccAddr := types.NewAddressFromSubaccount(1)
	k.SetAccountSummary(ctx, subAccAddr, balance)

	// Get balance
	gotBalance, exists := k.GetAccountSummary(ctx, subAccAddr)
	require.True(t, exists)
	require.Equal(t, balance, gotBalance)
}

func TestKeeper_GetLockedBalances(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	beforeUnlockTime1 := uint64(time.Now().Add(-time.Hour * 24 * 365).Unix())
	beforeUnlockTime2 := uint64(time.Now().Add(-time.Hour * 24 * 365 * 2).Unix())

	afterUnlockTime1 := uint64(time.Now().Add(time.Hour * 24 * 365).Unix())
	afterUnlockTime2 := uint64(time.Now().Add(time.Hour * 24 * 365 * 2).Unix())

	// I added them unordered to make sure they are sorted
	balanceUnlocks := []types.LockedBalance{
		{
			Amount:   sdkmath.NewInt(10000),
			UnlockTS: beforeUnlockTime1,
		},
		{
			Amount:   sdkmath.NewInt(30000),
			UnlockTS: afterUnlockTime1,
		},
		{
			Amount:   sdkmath.NewInt(20000),
			UnlockTS: beforeUnlockTime2,
		},
		{
			Amount:   sdkmath.NewInt(40000),
			UnlockTS: afterUnlockTime2,
		},
	}

	addr := types.NewAddressFromSubaccount(1)
	k.SetLockedBalances(ctx, addr, balanceUnlocks)

	// get unlocked balance
	_, unlockedAmount := k.GetBalances(ctx, addr, types.LockedBalanceStatus_LOCKED_BALANCE_STATUS_UNLOCKED)
	require.True(t, unlockedAmount.Equal(sdkmath.NewInt(10000+20000)))
}
