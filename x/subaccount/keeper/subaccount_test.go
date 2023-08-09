package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/types"

	"github.com/stretchr/testify/require"
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

func TestSubAccountOwner(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	owner := sample.NativeAccAddress()

	// Account should not have subaccount
	_, exists := k.GetSubAccountOwner(ctx, owner)
	require.False(t, exists)

	// Set subaccount owner
	id := k.NextID(ctx)
	k.SetSubAccountOwner(ctx, types.NewAddressFromSubaccount(id), owner)

	// Account should have subaccount
	_, exists = k.GetSubAccountByOwner(ctx, owner)
	require.True(t, exists)

	// Get subaccount ID
	subAccountAddress, exists := k.GetSubAccountByOwner(ctx, owner)
	require.True(t, exists)
	require.Equal(t, types.NewAddressFromSubaccount(id), subAccountAddress)

	// Get owner of subaccount
	gotOwner, exists := k.GetSubAccountOwner(ctx, subAccountAddress)
	require.True(t, exists)
	require.Equal(t, owner, gotOwner)
	// Get account ID by owner
	gotSubAccount, exists := k.GetSubAccountByOwner(ctx, owner)
	require.True(t, exists)
	require.Equal(t, types.NewAddressFromSubaccount(id), gotSubAccount)
}

func TestSetLockedBalances(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	someUnlockTime := time.Now().Add(time.Hour * 24 * 365)
	otherUnlockTime := time.Now().Add(time.Hour * 24 * 365 * 2)

	balanceUnlocks := []types.LockedBalance{
		{
			Amount:     sdk.NewInt(10000),
			UnlockTime: someUnlockTime,
		},
		{
			Amount:     sdk.NewInt(20000),
			UnlockTime: otherUnlockTime,
		},
	}

	addr := types.NewAddressFromSubaccount(1)

	k.SetLockedBalances(ctx, addr, balanceUnlocks)

	// Get locked balances
	lockedBalances := k.GetLockedBalances(ctx, addr)
	for i, lockedBalance := range lockedBalances {
		require.Equal(t, lockedBalance.Amount, balanceUnlocks[i].Amount)
		require.True(t, lockedBalance.UnlockTime.Equal(balanceUnlocks[i].UnlockTime))
	}
}

func TestSetBalances(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	balance := types.Balance{
		DepositedAmount: sdk.ZeroInt(),
		SpentAmount:     sdk.ZeroInt(),
		WithdrawmAmount: sdk.ZeroInt(),
		LostAmount:      sdk.OneInt(),
	}

	subAccAddr := types.NewAddressFromSubaccount(1)
	k.SetBalance(ctx, subAccAddr, balance)

	// Get balance
	gotBalance, exists := k.GetBalance(ctx, subAccAddr)
	require.True(t, exists)
	require.Equal(t, balance, gotBalance)
}

func TestKeeper_GetLockedBalances(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	beforeUnlockTime1 := time.Now().Add(-time.Hour * 24 * 365)
	beforeUnlockTime2 := time.Now().Add(-time.Hour * 24 * 365 * 2)

	afterUnlockTime1 := time.Now().Add(time.Hour * 24 * 365)
	afterUnlockTime2 := time.Now().Add(time.Hour * 24 * 365 * 2)

	// I added them unordered to make sure they are sorted
	balanceUnlocks := []types.LockedBalance{
		{
			Amount:     sdk.NewInt(10000),
			UnlockTime: beforeUnlockTime1,
		},
		{
			Amount:     sdk.NewInt(30000),
			UnlockTime: afterUnlockTime1,
		},
		{
			Amount:     sdk.NewInt(20000),
			UnlockTime: beforeUnlockTime2,
		},
		{
			Amount:     sdk.NewInt(40000),
			UnlockTime: afterUnlockTime2,
		},
	}

	addr := types.NewAddressFromSubaccount(1)
	k.SetLockedBalances(ctx, addr, balanceUnlocks)

	// get unlocked balance
	unlockedBalance := k.GetUnlockedBalance(ctx, addr)
	require.True(t, unlockedBalance.Equal(sdk.NewInt(10000+20000)))
}
