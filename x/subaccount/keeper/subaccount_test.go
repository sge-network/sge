package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/subaccount/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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

func TestSubAccountOwner(t *testing.T) {
	_, k, ctx := setupKeeperAndApp(t)

	address := sample.AccAddress()

	// Account should not have subaccount
	require.Equal(t, false, k.HasSubAccount(ctx, address))

	// Set subaccount owner
	k.NextID(ctx)
	ID := k.Peek(ctx)
	k.SetSubAccountOwner(ctx, ID, address)

	// Account should have subaccount
	require.True(t, k.HasSubAccount(ctx, address))

	// Get subaccount ID
	require.Equal(t, ID, k.GetSubAccountByOwner(ctx, address))

	// Get owner of subaccount
	require.Equal(t, address, k.GetSubAccountOwner(ctx, ID))
	// Get account ID by owner
	require.Equal(t, ID, k.GetSubAccountByOwner(ctx, address))
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

	k.SetLockedBalances(ctx, 1, balanceUnlocks)
}
