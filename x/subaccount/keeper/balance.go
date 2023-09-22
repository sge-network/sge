package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/subaccount/types"
)

// SetLockedBalances saves the locked balances of an account.
func (k Keeper) SetLockedBalances(ctx sdk.Context, subAccountAddress sdk.AccAddress, lockedBalances []types.LockedBalance) {
	store := ctx.KVStore(k.storeKey)

	for _, lockedBalance := range lockedBalances {
		amountBytes, err := lockedBalance.Amount.Marshal()
		if err != nil {
			panic(err)
		}
		store.Set(
			types.LockedBalanceKey(subAccountAddress, lockedBalance.UnlockTS),
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
		unlockTime := utils.Uint64FromBytes(iterator.Key())

		amount := new(sdkmath.Int)
		err := amount.Unmarshal(iterator.Value())
		if err != nil {
			panic(err)
		}
		lockedBalances = append(lockedBalances, types.LockedBalance{
			UnlockTS: unlockTime,
			Amount:   *amount,
		})
	}

	return lockedBalances
}

// GetUnlockedBalance returns the unlocked balance of an account.
func (k Keeper) GetUnlockedBalance(ctx sdk.Context, subAccountAddress sdk.AccAddress) sdkmath.Int {
	iterator := prefix.NewStore(ctx.KVStore(k.storeKey), types.LockedBalancePrefixKey(subAccountAddress)).
		Iterator(nil, utils.Int64ToBytes(ctx.BlockTime().Unix()))

	unlockedBalance := sdk.ZeroInt()
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		amount := new(sdkmath.Int)
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

// getBalances returns the balance, unlocked balance and bank balance of a subaccount
func (k Keeper) getBalances(sdkContext sdk.Context, subaccountAddr sdk.AccAddress) (types.Balance, sdkmath.Int, sdk.Coin) {
	balance, exists := k.GetBalance(sdkContext, subaccountAddr)
	if !exists {
		panic("data corruption: subaccount exists but balance does not")
	}
	unlockedBalance := k.GetUnlockedBalance(sdkContext, subaccountAddr)
	bankBalance := k.bankKeeper.GetBalance(sdkContext, subaccountAddr, params.DefaultBondDenom)

	return balance, unlockedBalance, bankBalance
}
