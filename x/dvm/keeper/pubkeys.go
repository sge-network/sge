package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/types"
)

// SetPublicKeys sets the list of keys, overwrite the old values.
func (k Keeper) SetPublicKeys(ctx sdk.Context, ks types.PublicKeys) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PubKeysListKey)
	b := k.cdc.MustMarshal(&ks)
	store.Set([]byte{0}, b)
}

// GetPublicKeys is the helper functions for this keeper to query the list of public keys.
func (k *Keeper) GetPublicKeys(ctx sdk.Context) (keys types.PublicKeys, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PubKeysListKey)

	b := store.Get([]byte{0})
	if b == nil {
		return keys, false
	}

	k.cdc.MustUnmarshal(b, &keys)
	return keys, true
}
