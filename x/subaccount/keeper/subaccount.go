package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

// GetSubAccountOwner returns the owner of a subaccount.
func (k Keeper) GetSubAccountOwner(ctx sdk.Context, id uint64) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get(subaccounttypes.SubAccountKey(id))
}

// SetLockedBalances saves the locked balances of an account.
func (k Keeper) SetLockedBalances(ctx sdk.Context, account sdk.AccAddress, lockedBalances []*subaccounttypes.LockedBalance) {
	store := ctx.KVStore(k.storeKey)

	for _, lockedBalance := range lockedBalances {
		store.Set(subaccounttypes.LockedBalanceKey(account, lockedBalance.UnlockTime), sdk.Uint64ToBigEndian(lockedBalance.Amount.Uint64()))
	}
}

// GetLockedBalances returns the locked balances of an account.
func (k Keeper) GetLockedBalances(ctx sdk.Context, account sdk.AccAddress) []subaccounttypes.LockedBalance {
	iterator := prefix.NewStore(ctx.KVStore(k.storeKey), subaccounttypes.LockedBalancePrefixKey(account)).Iterator(nil, nil)

	lockedBalances := []subaccounttypes.LockedBalance{}

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		unlockTime, err := sdk.ParseTimeBytes(iterator.Key())
		if err != nil {
			panic(err)
		}

		amount := sdk.BigEndianToUint64(iterator.Value())
		lockedBalances = append(lockedBalances, subaccounttypes.LockedBalance{
			UnlockTime: unlockTime,
			Amount:     sdk.NewIntFromUint64(amount),
		})
	}

	return lockedBalances
}
