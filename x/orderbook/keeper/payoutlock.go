package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/utils"
)

// SetPayoutLock sets a lock for the payout element
func (k Keeper) SetPayoutLock(ctx sdk.Context, uniqueLock string) {
	store := k.getPayoutLockStore(ctx)
	store.Set(utils.StrBytes(uniqueLock), []byte{1})
}

// removePayoutLock removes a lock from payout
func (k Keeper) removePayoutLock(ctx sdk.Context, uniqueLock string) {
	store := k.getPayoutLockStore(ctx)
	store.Delete(utils.StrBytes(uniqueLock))
}

// GetAllPayoutLock returns all payout locks used during genesis dump.
func (k Keeper) GetAllPayoutLock(ctx sdk.Context) (list [][]byte, err error) {
	store := k.getPayoutLockStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Value())
	}

	return
}

// getPayoutLock checks if payout lock exists
func (k Keeper) payoutLockExists(ctx sdk.Context, uniqueLock string) bool {
	store := k.getPayoutLockStore(ctx)
	return store.Has(utils.StrBytes(uniqueLock))
}
