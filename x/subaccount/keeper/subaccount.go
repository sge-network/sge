package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/app/params"
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
func (k Keeper) GetSubAccountOwner(ctx sdk.Context, subAccAddr sdk.AccAddress) (sdk.AccAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	addr := store.Get(types.SubAccountKey(subAccAddr))
	return addr, addr != nil
}

// IsSubAccount returns true if the address blongs to a sub account.
func (k Keeper) IsSubAccount(ctx sdk.Context, subAccAddr sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.SubAccountKey(subAccAddr))
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
	k.IterateSubaccounts(ctx, func(subAccountAddress, ownerAddress sdk.AccAddress) (stop bool) {
		balance, exists := k.GetAccountSummary(ctx, subAccountAddress)
		if !exists {
			panic("subaccount balance does not exist")
		}
		lockedBalances, _ := k.GetBalances(ctx, subAccountAddress, types.BalanceType_BALANCE_TYPE_UNSPECIFIED)
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

func (k Keeper) CreateSubAccount(ctx sdk.Context, creator, owner string,
	lockedBalances []types.LockedBalance,
) (string, error) {
	lockedBalance, err := sumlockedBalance(ctx, lockedBalances)
	if err != nil {
		return "", err
	}

	creatorAddr := sdk.MustAccAddressFromBech32(creator)
	subAccOwnerAddr := sdk.MustAccAddressFromBech32(owner)
	if _, exists := k.GetSubAccountByOwner(ctx, subAccOwnerAddr); exists {
		return "", types.ErrSubaccountAlreadyExist
	}

	subaccountID := k.NextID(ctx)

	// ALERT: If someone frontruns the account creation, will be overwritten here
	subAccAddr := types.NewAddressFromSubaccount(subaccountID)
	subaccountAccount := k.accountKeeper.NewAccountWithAddress(ctx, subAccAddr)
	k.accountKeeper.SetAccount(ctx, subaccountAccount)

	err = k.sendCoinsToSubaccount(ctx, creatorAddr, subAccAddr, lockedBalance)
	if err != nil {
		return "", sdkerrors.Wrap(err, "unable to send coins")
	}

	k.SetSubAccountOwner(ctx, subAccAddr, subAccOwnerAddr)
	k.SetLockedBalances(ctx, subAccAddr, lockedBalances)
	k.SetAccountSummary(ctx, subAccAddr, types.AccountSummary{
		DepositedAmount: lockedBalance,
		SpentAmount:     sdk.ZeroInt(),
		WithdrawnAmount: sdk.ZeroInt(),
		LostAmount:      sdk.ZeroInt(),
	})
	return subAccAddr.String(), nil
}

// sendCoinsToSubaccount sends the coins to the subaccount.
func (k Keeper) sendCoinsToSubaccount(ctx sdk.Context, creatorAccount, subAccountAddress sdk.AccAddress, moneyToSend sdkmath.Int) error {
	err := k.bankKeeper.SendCoins(ctx, creatorAccount, subAccountAddress, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, moneyToSend)))
	if err != nil {
		return sdkerrors.Wrap(err, "unable to send coins")
	}

	return nil
}

// sumlockedBalance sums all the balances to unlock and returns the total amount. It
// returns an error if any of to unlock times is expired.
func sumlockedBalance(ctx sdk.Context, lockedBalances []types.LockedBalance) (sdkmath.Int, error) {
	lockedBalance := sdkmath.NewInt(0)

	for _, lb := range lockedBalances {
		// return error if balance is unlocked
		if lb.UnlockTS < cast.ToUint64(ctx.BlockTime().Unix()) {
			return sdkmath.Int{}, types.ErrUnlockTokenTimeExpired
		}

		lockedBalance = lockedBalance.Add(lb.Amount)
	}

	return lockedBalance, nil
}
