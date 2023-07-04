package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

type Keeper struct {
	cdc      codec.Codec
	storeKey sdk.StoreKey
}

func (k Keeper) SetLockedBalances(ctx sdk.Context, subaccountID uint64, lockedBalances []types.LockedBalance) {
	store := ctx.KVStore(k.storeKey)

	for _, lockedBalance := range lockedBalances {
		account := types.NewModuleAccountFromSubAccount(subaccountID)
		store.Set(types.LockedBalanceKey(account, lockedBalance.UnlockTime), sdk.Uint64ToBigEndian(lockedBalance.Amount.Uint64()))
	}
}

func NewKeeper(storeKey sdk.StoreKey) Keeper {
	return Keeper{storeKey: storeKey}
}
