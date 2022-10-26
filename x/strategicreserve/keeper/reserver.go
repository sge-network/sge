package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// GetReserver returns the reserver
func (k Keeper) GetReserver(ctx sdk.Context) (reserver types.Reserver) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.ReserverKey)
	if b == nil {
		panic(ErrTextNilReserver)
	}

	k.cdc.MustUnmarshal(b, &reserver)
	return
}

// SetReserver sets the reserver
func (k Keeper) SetReserver(ctx sdk.Context, reserver types.Reserver) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&reserver)
	store.Set(types.ReserverKey, b)
}
