package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/ovm/types"
)

// SetKeyVault sets the key vault and overwrite the old values.
func (k Keeper) SetKeyVault(ctx sdk.Context, ks types.KeyVault) {
	store := k.getKeyVaultStore(ctx)
	b := k.cdc.MustMarshal(&ks)
	store.Set([]byte{0}, b)
}

// GetKeyVault is the helper functions for this keeper to query the key vault.
func (k *Keeper) GetKeyVault(ctx sdk.Context) (keys types.KeyVault, found bool) {
	store := k.getKeyVaultStore(ctx)

	b := store.Get([]byte{0})
	if b == nil {
		return keys, false
	}

	k.cdc.MustUnmarshal(b, &keys)
	return keys, true
}
