package keeper

import (
	sdkerrors "cosmossdk.io/errors"
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

// SetAccountSummary saves the balance of an account.
func (k Keeper) SetAccountSummary(ctx sdk.Context, subAccountAddress sdk.AccAddress, accountSummary types.AccountSummary) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&accountSummary)
	store.Set(types.AccountSummaryKey(subAccountAddress), bz)
}

// GetAccountSummary returns the balance of an account.
func (k Keeper) GetAccountSummary(ctx sdk.Context, subAccountAddress sdk.AccAddress) (types.AccountSummary, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AccountSummaryKey(subAccountAddress))
	if bz == nil {
		return types.AccountSummary{}, false
	}

	balance := types.AccountSummary{}
	k.cdc.MustUnmarshal(bz, &balance)

	return balance, true
}

// TopUp tops up the subaccount balance.
func (k Keeper) TopUp(ctx sdk.Context, creator, subAccOwnerAddr string,
	lockedBalance []types.LockedBalance,
) (string, error) {
	addedBalance, err := sumlockedBalance(ctx, lockedBalance)
	if err != nil {
		return "", err
	}

	creatorAddr := sdk.MustAccAddressFromBech32(creator)
	subaccountOwner := sdk.MustAccAddressFromBech32(subAccOwnerAddr)

	subAccAddr, exists := k.GetSubAccountByOwner(ctx, subaccountOwner)
	if !exists {
		return "", types.ErrSubaccountDoesNotExist
	}
	balance, exists := k.GetAccountSummary(ctx, subAccAddr)
	if !exists {
		panic("data corruption: subaccount exists but balance does not")
	}

	balance.DepositedAmount = balance.DepositedAmount.Add(addedBalance)
	k.SetAccountSummary(ctx, subAccAddr, balance)
	k.SetLockedBalances(ctx, subAccAddr, lockedBalance)

	err = k.sendCoinsToSubaccount(ctx, creatorAddr, subAccAddr, addedBalance)
	if err != nil {
		return "", sdkerrors.Wrapf(types.ErrSendCoinError, "%s", err)
	}
	return subAccAddr.String(), nil
}

// getAccountSummary returns the balance, unlocked balance and bank balance of a subaccount
func (k Keeper) getAccountSummary(sdkContext sdk.Context, subaccountAddr sdk.AccAddress) (types.AccountSummary, sdkmath.Int, sdk.Coin) {
	balance, exists := k.GetAccountSummary(sdkContext, subaccountAddr)
	if !exists {
		panic("data corruption: subaccount exists but balance does not")
	}
	unlockedBalance := k.GetUnlockedBalance(sdkContext, subaccountAddr)
	bankBalance := k.bankKeeper.GetBalance(sdkContext, subaccountAddr, params.DefaultBondDenom)

	return balance, unlockedBalance, bankBalance
}

// getBalances returns the balance, unlocked balance and bank balance of a subaccount
func (k Keeper) withdraw(ctx sdk.Context, subAccAddr sdk.AccAddress, ownerAddr sdk.AccAddress) error {
	summary, unlockedBalance, bankBalance := k.getAccountSummary(ctx, subAccAddr)

	withdrawableBalance := summary.WithdrawableBalance(unlockedBalance, bankBalance.Amount)
	if withdrawableBalance.IsZero() {
		return types.ErrNothingToWithdraw
	}

	summary.Withdraw(withdrawableBalance)
	k.SetAccountSummary(ctx, subAccAddr, summary)

	err := k.bankKeeper.SendCoins(ctx, subAccAddr, ownerAddr, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, withdrawableBalance)))
	if err != nil {
		return err
	}

	return nil
}
