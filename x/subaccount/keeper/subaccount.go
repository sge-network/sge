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
		return 1
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
func (k Keeper) SetID(ctx sdk.Context, id uint64) {
	ctx.KVStore(k.storeKey).Set(subaccounttypes.SubaccountIDPrefix, sdk.Uint64ToBigEndian(id))
}

// SetSubAccountOwner sets the owner of a subaccount.
func (k Keeper) SetSubAccountOwner(ctx sdk.Context, subAccountAddress, ownerAddress sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(subaccounttypes.SubAccountOwnerKey(ownerAddress), subAccountAddress)
	// and reverse mapping
	store.Set(subaccounttypes.SubAccountKey(subAccountAddress), ownerAddress)
}

// GetSubAccountByOwner returns the subaccount ID of an owner.
func (k Keeper) GetSubAccountByOwner(ctx sdk.Context, address sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	addr := store.Get(subaccounttypes.SubAccountOwnerKey(address))
	return addr, addr != nil
}

// GetSubAccountOwner returns the owner of a subaccount given the subaccount address.
func (k Keeper) GetSubAccountOwner(ctx sdk.Context, subaccountAddr sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	addr := store.Get(subaccounttypes.SubAccountKey(subaccountAddr))
	return addr, addr != nil
}

// SetLockedBalances saves the locked balances of an account.
func (k Keeper) SetLockedBalances(ctx sdk.Context, subAccountAddress sdk.AccAddress, lockedBalances []subaccounttypes.LockedBalance) {
	store := ctx.KVStore(k.storeKey)

	for _, lockedBalance := range lockedBalances {
		amountBytes, err := lockedBalance.Amount.Marshal()
		if err != nil {
			panic(err)
		}
		store.Set(
			subaccounttypes.LockedBalanceKey(subAccountAddress, lockedBalance.UnlockTime),
			amountBytes,
		)
	}
}

// GetLockedBalances returns the locked balances of an account.
func (k Keeper) GetLockedBalances(ctx sdk.Context, subAccountAddress sdk.AccAddress) []subaccounttypes.LockedBalance {
	iterator := prefix.NewStore(ctx.KVStore(k.storeKey), subaccounttypes.LockedBalancePrefixKey(subAccountAddress)).Iterator(nil, nil)
	defer iterator.Close()

	var lockedBalances []subaccounttypes.LockedBalance
	for ; iterator.Valid(); iterator.Next() {
		unlockTime, err := sdk.ParseTimeBytes(iterator.Key())
		if err != nil {
			panic(err)
		}

		amount := new(sdk.Int)
		err = amount.Unmarshal(iterator.Value())
		if err != nil {
			panic(err)
		}
		lockedBalances = append(lockedBalances, subaccounttypes.LockedBalance{
			UnlockTime: unlockTime,
			Amount:     *amount,
		})
	}

	return lockedBalances
}

// GetUnlockedBalance returns the unlocked balance of an account.
func (k Keeper) GetUnlockedBalance(ctx sdk.Context, subAccountAddress sdk.AccAddress) sdk.Int {
	iterator := prefix.NewStore(ctx.KVStore(k.storeKey), subaccounttypes.LockedBalancePrefixKey(subAccountAddress)).
		Iterator(nil, sdk.FormatTimeBytes(ctx.BlockTime()))

	unlockedBalance := sdk.ZeroInt()
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		amount := new(sdk.Int)
		err := amount.Unmarshal(iterator.Value())
		if err != nil {
			panic(err)
		}
		unlockedBalance = unlockedBalance.Add(*amount)
	}

	return unlockedBalance
}

// SetBalance saves the balance of an account.
func (k Keeper) SetBalance(ctx sdk.Context, subAccountAddress sdk.AccAddress, balance subaccounttypes.Balance) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&balance)
	store.Set(subaccounttypes.BalanceKey(subAccountAddress), bz)
}

// GetBalance returns the balance of an account.
func (k Keeper) GetBalance(ctx sdk.Context, subAccountAddress sdk.AccAddress) (subaccounttypes.Balance, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(subaccounttypes.BalanceKey(subAccountAddress))
	if bz == nil {
		return subaccounttypes.Balance{}, false
	}

	balance := subaccounttypes.Balance{}
	k.cdc.MustUnmarshal(bz, &balance)

	return balance, true
}

func (k Keeper) IterateSubaccounts(ctx sdk.Context, cb func(subAccountAddress, subaccountOwner sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, subaccounttypes.SubAccountOwnerReversePrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if cb(iterator.Key()[len(subaccounttypes.SubAccountOwnerReversePrefix):], iterator.Value()) {
			break
		}
	}
}

func (k Keeper) GetAllSubaccounts(ctx sdk.Context) []subaccounttypes.GenesisSubaccount {
	var subaccounts []subaccounttypes.GenesisSubaccount
	k.IterateSubaccounts(ctx, func(subAccountAddress sdk.AccAddress, ownerAddress sdk.AccAddress) (stop bool) {
		balance, exists := k.GetBalance(ctx, subAccountAddress)
		if !exists {
			panic("subaccount balance does not exist")
		}
		lockedBalances := k.GetLockedBalances(ctx, subAccountAddress)
		subaccounts = append(subaccounts, subaccounttypes.GenesisSubaccount{
			Address:        subAccountAddress.String(),
			Owner:          ownerAddress.String(),
			Balance:        balance,
			LockedBalances: lockedBalances,
		})
		return false
	})
	return subaccounts
}
