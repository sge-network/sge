package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/subaccount/types"
)

// Peek returns the next value without advancing the subaccount ID.
func (k Keeper) Peek(ctx sdk.Context) uint64 {
	key := ctx.KVStore(k.storeKey).Get(types.SubaccountIDPrefix)
	if key == nil {
		return 1
	}

	return sdk.BigEndianToUint64(key)
}

// NextID returns the actual value, same as Peek, but also advances the subaccount ID.
func (k Keeper) NextID(ctx sdk.Context) uint64 {
	actualID := k.Peek(ctx)

	ctx.KVStore(k.storeKey).Set(types.SubaccountIDPrefix, sdk.Uint64ToBigEndian(actualID+1))

	return actualID
}

// SetID sets the ID to a given value.
func (k Keeper) SetID(ctx sdk.Context, id uint64) {
	ctx.KVStore(k.storeKey).Set(types.SubaccountIDPrefix, sdk.Uint64ToBigEndian(id))
}

// SetSubAccountOwner sets the owner of a subaccount.
func (k Keeper) SetSubAccountOwner(ctx sdk.Context, subAccountAddress, ownerAddress sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.SubAccountOwnerKey(ownerAddress), subAccountAddress)
	// and reverse mapping
	store.Set(types.SubAccountKey(subAccountAddress), ownerAddress)
}

// GetSubAccountByOwner returns the subaccount ID of an owner.
func (k Keeper) GetSubAccountByOwner(ctx sdk.Context, address sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	addr := store.Get(types.SubAccountOwnerKey(address))
	return addr, addr != nil
}

// GetSubAccountOwner returns the owner of a subaccount given the subaccount address.
func (k Keeper) GetSubAccountOwner(ctx sdk.Context, subaccountAddr sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	addr := store.Get(types.SubAccountKey(subaccountAddr))
	return addr, addr != nil
}

// SetLockedBalances saves the locked balances of an account.
func (k Keeper) SetLockedBalances(ctx sdk.Context, subAccountAddress sdk.AccAddress, lockedBalances []types.LockedBalance) {
	store := ctx.KVStore(k.storeKey)

	for _, lockedBalance := range lockedBalances {
		amountBytes, err := lockedBalance.Amount.Marshal()
		if err != nil {
			panic(err)
		}
		store.Set(
			types.LockedBalanceKey(subAccountAddress, lockedBalance.UnlockTime),
			amountBytes,
		)
	}
}

// GetLockedBalances returns the locked balances of an account.
func (k Keeper) GetLockedBalances(ctx sdk.Context, subAccountAddress sdk.AccAddress) []types.LockedBalance {
	iterator := prefix.NewStore(ctx.KVStore(k.storeKey), types.LockedBalancePrefixKey(subAccountAddress)).Iterator(nil, nil)
	defer iterator.Close()

	var lockedBalances []types.LockedBalance
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
		lockedBalances = append(lockedBalances, types.LockedBalance{
			UnlockTime: unlockTime,
			Amount:     *amount,
		})
	}

	return lockedBalances
}

// GetUnlockedBalance returns the unlocked balance of an account.
func (k Keeper) GetUnlockedBalance(ctx sdk.Context, subAccountAddress sdk.AccAddress) sdk.Int {
	iterator := prefix.NewStore(ctx.KVStore(k.storeKey), types.LockedBalancePrefixKey(subAccountAddress)).
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
func (k Keeper) SetBalance(ctx sdk.Context, subAccountAddress sdk.AccAddress, balance types.Balance) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&balance)
	store.Set(types.BalanceKey(subAccountAddress), bz)
}

// GetBalance returns the balance of an account.
func (k Keeper) GetBalance(ctx sdk.Context, subAccountAddress sdk.AccAddress) (types.Balance, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BalanceKey(subAccountAddress))
	if bz == nil {
		return types.Balance{}, false
	}

	balance := types.Balance{}
	k.cdc.MustUnmarshal(bz, &balance)

	return balance, true
}

func (k Keeper) IterateSubaccounts(ctx sdk.Context, cb func(subAccountAddress, subaccountOwner sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.SubAccountOwnerReversePrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		if cb(iterator.Key()[len(types.SubAccountOwnerReversePrefix):], iterator.Value()) {
			break
		}
	}
}

func (k Keeper) GetAllSubaccounts(ctx sdk.Context) []types.GenesisSubaccount {
	var subaccounts []types.GenesisSubaccount
	k.IterateSubaccounts(ctx, func(subAccountAddress sdk.AccAddress, ownerAddress sdk.AccAddress) (stop bool) {
		balance, exists := k.GetBalance(ctx, subAccountAddress)
		if !exists {
			panic("subaccount balance does not exist")
		}
		lockedBalances := k.GetLockedBalances(ctx, subAccountAddress)
		subaccounts = append(subaccounts, types.GenesisSubaccount{
			Address:        subAccountAddress.String(),
			Owner:          ownerAddress.String(),
			Balance:        balance,
			LockedBalances: lockedBalances,
		})
		return false
	})
	return subaccounts
}

// sendCoinsToSubaccount sends the coins to the subaccount.
func (k Keeper) sendCoinsToSubaccount(ctx sdk.Context, senderAccount sdk.AccAddress, subAccountAddress sdk.AccAddress, moneyToSend sdk.Int) error {
	denom := k.GetParams(ctx).LockedBalanceDenom
	err := k.bankKeeper.SendCoins(ctx, senderAccount, subAccountAddress, sdk.NewCoins(sdk.NewCoin(denom, moneyToSend)))
	if err != nil {
		return errors.Wrap(err, "unable to send coins")
	}

	return nil
}

// sumBalanceUnlocks sums all the balances to unlock and returns the total amount. It
// returns an error if any of the unlock times is expired.
func sumBalanceUnlocks(ctx sdk.Context, balanceUnlocks []types.LockedBalance) (sdk.Int, error) {
	moneyToSend := sdk.NewInt(0)

	for _, balanceUnlock := range balanceUnlocks {
		if balanceUnlock.UnlockTime.Unix() < ctx.BlockTime().Unix() {
			return sdk.Int{}, types.ErrUnlockTokenTimeExpired
		}

		moneyToSend = moneyToSend.Add(balanceUnlock.Amount)
	}

	return moneyToSend, nil
}
