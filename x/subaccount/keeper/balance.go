package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

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

// etBalances returns the locked balances of an account.
func (k Keeper) GetBalances(ctx sdk.Context, subAccountAddress sdk.AccAddress, balanceType types.LockedBalanceStatus) ([]types.LockedBalance, sdkmath.Int) {
	var start, end []byte
	switch balanceType {
	case types.LockedBalanceStatus_LOCKED_BALANCE_STATUS_LOCKED:
		start = utils.Int64ToBytes(ctx.BlockTime().Unix())
	case types.LockedBalanceStatus_LOCKED_BALANCE_STATUS_UNLOCKED:
		end = utils.Int64ToBytes(ctx.BlockTime().Unix())
	}

	iterator := prefix.NewStore(ctx.KVStore(k.storeKey), types.LockedBalancePrefixKey(subAccountAddress)).Iterator(start, end)
	defer iterator.Close()

	var lockedBalances []types.LockedBalance
	totalAmount := sdkmath.ZeroInt()
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
		totalAmount = totalAmount.Add(*amount)
	}

	return lockedBalances, totalAmount
}

// HasLockedBalances returns true if there is an existing locked balance for the sub account address.
func (k Keeper) HasLockedBalances(ctx sdk.Context, subAccAddr sdk.AccAddress, time uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.LockedBalanceKey(subAccAddr, time))
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

	accSummary := types.AccountSummary{}
	k.cdc.MustUnmarshal(bz, &accSummary)

	return accSummary, true
}

// TopUp tops up the subaccount balance.
func (k Keeper) TopUp(ctx sdk.Context, creator, subAccOwnerAddr string,
	topUpBalances []types.LockedBalance,
) (string, error) {
	addedBalance, err := sumLockedBalance(ctx, topUpBalances)
	if err != nil {
		return "", err
	}

	creatorAddr := sdk.MustAccAddressFromBech32(creator)
	subaccountOwner := sdk.MustAccAddressFromBech32(subAccOwnerAddr)

	subAccAddr, exists := k.GetSubaccountByOwner(ctx, subaccountOwner)
	if !exists {
		return "", types.ErrSubaccountDoesNotExist
	}
	accSummary, exists := k.GetAccountSummary(ctx, subAccAddr)
	if !exists {
		panic("data corruption: subaccount exists but balance does not")
	}

	for _, v := range topUpBalances {
		if k.HasLockedBalances(ctx, subAccAddr, v.UnlockTS) {
			return "", sdkerrors.Wrapf(types.ErrLockedBalanceExists, "%d", v.UnlockTS)
		}
	}

	accSummary.DepositedAmount = accSummary.DepositedAmount.Add(addedBalance)
	k.SetAccountSummary(ctx, subAccAddr, accSummary)
	k.SetLockedBalances(ctx, subAccAddr, topUpBalances)

	err = k.sendCoinsToSubaccount(ctx, creatorAddr, subAccAddr, addedBalance)
	if err != nil {
		return "", sdkerrors.Wrapf(types.ErrSendCoinError, "%s", err)
	}
	return subAccAddr.String(), nil
}

// getSubaccountSummary returns the balance, unlocked balance and bank balance of a subaccount
func (k Keeper) getSubaccountSummary(sdkContext sdk.Context, subaccountAddr sdk.AccAddress) (types.AccountSummary, sdkmath.Int, sdk.Coin) {
	accSummary, exists := k.GetAccountSummary(sdkContext, subaccountAddr)
	if !exists {
		panic("data corruption: subaccount exists but balance does not")
	}
	_, unlockedAmount := k.GetBalances(sdkContext, subaccountAddr, types.LockedBalanceStatus_LOCKED_BALANCE_STATUS_UNLOCKED)
	bankBalance := k.bankKeeper.GetBalance(sdkContext, subaccountAddr, params.DefaultBondDenom)

	return accSummary, unlockedAmount, bankBalance
}

// withdrawUnlocked returns the balance, unlocked balance and bank balance of a subaccount
func (k Keeper) withdrawUnlocked(ctx sdk.Context, subAccAddr sdk.AccAddress, ownerAddr sdk.AccAddress) error {
	accSummary, unlockedAmount, bankBalance := k.getSubaccountSummary(ctx, subAccAddr)

	withdrawableUnlockedBalance := accSummary.WithdrawableUnlockedBalance(unlockedAmount, bankBalance.Amount)
	if withdrawableUnlockedBalance.IsZero() {
		return types.ErrNothingToWithdraw
	}

	if err := accSummary.Withdraw(withdrawableUnlockedBalance); err != nil {
		return err
	}

	k.SetAccountSummary(ctx, subAccAddr, accSummary)

	err := k.bankKeeper.SendCoins(ctx, subAccAddr, ownerAddr, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, withdrawableUnlockedBalance)))
	if err != nil {
		return err
	}

	return nil
}

// withdrawLockedAndUnlocked withdraws unlocked balance first, if more balance is needed to be deducted,
// modifies the locked balances accordingly
func (k Keeper) withdrawLockedAndUnlocked(ctx sdk.Context, subAccAddr sdk.AccAddress, ownerAddr sdk.AccAddress, subAmountDeduct sdkmath.Int,
) error {
	accSummary, _, bankBalance := k.getSubaccountSummary(ctx, subAccAddr)
	withdrawableBalance := accSummary.WithdrawableBalance(bankBalance.Amount)

	// take the minimum of the withdrawable balance and the amount that need to be transferred from sub account
	withdrawableToSend := sdkmath.MinInt(withdrawableBalance, subAmountDeduct)

	// check total withdrawable balance to be enough for withdrawal.
	if subAmountDeduct.GT(withdrawableToSend) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest,
			"not enough balance in sub account locked balances, need more %s tokens",
			subAmountDeduct.Sub(withdrawableToSend))
	}

	// send the total calculated amount to the owner
	if err := k.bankKeeper.SendCoins(ctx,
		subAccAddr,
		ownerAddr,
		sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, subAmountDeduct))); err != nil {
		return sdkerrors.Wrapf(types.ErrSendCoinError, "error sending coin from subaccount to main account %s", err)
	}

	// update account summary withdrawn amount
	if err := accSummary.Withdraw(subAmountDeduct); err != nil {
		return sdkerrors.Wrapf(types.ErrWithdrawLocked, "%s", err)
	}
	k.SetAccountSummary(ctx, subAccAddr, accSummary)

	return nil
}
