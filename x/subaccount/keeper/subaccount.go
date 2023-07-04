package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

// Peek returns the next value without advancing the subaccount ID.
func (k Keeper) Peek(ctx sdk.Context) uint64 {
	key := ctx.KVStore(k.storeKey).Get(subaccounttypes.SubaccountIDPrefix)
	if key == nil {
		return 0
	}

	return sdk.BigEndianToUint64(key)
}

// NextID returns the actual value, same as Peek, but also advances the subaccount ID.
func (k Keeper) NextID(ctx sdk.Context) uint64 {
	actualID := k.Peek(ctx)

	ctx.KVStore(k.storeKey).Set(subaccounttypes.SubaccountIDPrefix, sdk.Uint64ToBigEndian(actualID+1))

	return actualID
}

// SetID sets the ID to a given value.
func (k Keeper) SetID(ctx sdk.Context, ID uint64) {
	ctx.KVStore(k.storeKey).Set(subaccounttypes.SubaccountIDPrefix, sdk.Uint64ToBigEndian(ID))
}

// HasSubAccount returns true if the account has a subaccount.
func (k Keeper) HasSubAccount(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(subaccounttypes.SubAccountOwnerKey(address))
}

// SetSubAccountOwner sets the owner of a subaccount.
func (k Keeper) SetSubAccountOwner(ctx sdk.Context, id uint64, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(subaccounttypes.SubAccountOwnerKey(address), sdk.Uint64ToBigEndian(id))
	// and reverse mapping
	store.Set(subaccounttypes.SubAccountKey(id), address.Bytes())
}

// GetSubAccountByOwner returns the subaccount ID of an owner.
func (k Keeper) GetSubAccountByOwner(ctx sdk.Context, address sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)

	return sdk.BigEndianToUint64(store.Get(subaccounttypes.SubAccountOwnerKey(address)))
}

func (k Keeper) GetSubAccountOwner(ctx sdk.Context, id uint64) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)

	return store.Get(subaccounttypes.SubAccountKey(id))
}
