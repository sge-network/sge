package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/reward/types"
)

// SetPromoter set a specific promoter in the store from its index
func (k Keeper) SetPromoter(ctx sdk.Context, promoter types.Promoter) {
	store := k.getPromoterStore(ctx)
	b := k.cdc.MustMarshal(&promoter)
	store.Set(types.GetPromoterKey(promoter.UID), b)
}

// GetPromoter returns a promoter,
func (k Keeper) GetPromoter(
	ctx sdk.Context,
	uid string,
) (val types.Promoter, found bool) {
	store := k.getPromoterStore(ctx)

	b := store.Get(types.GetPromoterKey(uid))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllPromoter returns all promoters
func (k Keeper) GetAllPromoter(ctx sdk.Context) (list []types.Promoter) {
	store := k.getPromoterStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Promoter
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// SetPromoterByAddress set a specific promoter in the store from its address
func (k Keeper) SetPromoterByAddress(ctx sdk.Context, promoterByAddress types.PromoterByAddress) {
	store := k.getPromoterByAddressStore(ctx)
	b := k.cdc.MustMarshal(&promoterByAddress)
	store.Set(types.GetPromoterByAddressKey(promoterByAddress.Address), b)
}

// GetPromoter returns a promoter,
func (k Keeper) GetPromoterByAddress(ctx sdk.Context, address string) (val types.PromoterByAddress, found bool) {
	store := k.getPromoterByAddressStore(ctx)

	b := store.Get(types.GetPromoterByAddressKey(address))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
