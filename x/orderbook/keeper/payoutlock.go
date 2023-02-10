package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/orderbook/types"
)

// getPayoutLock checks if payout lock exists
func (k Keeper) payoutLockExists(ctx sdk.Context, uniqueLock string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PayoutLockKeyPrefix)
	return store.Has(utils.StrBytes(uniqueLock))
}

// setPayoutLock sets a lock for the payout element
func (k Keeper) setPayoutLock(ctx sdk.Context, uniqueLock string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PayoutLockKeyPrefix)
	store.Set(utils.StrBytes(uniqueLock), []byte{1})
}
