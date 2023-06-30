package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey) Keeper {
	return Keeper{storeKey: storeKey}
}

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

func (k Keeper) SetID(ctx sdk.Context, ID uint64) {
	ctx.KVStore(k.storeKey).Set(subaccounttypes.SubaccountIDPrefix, sdk.Uint64ToBigEndian(ID))
}
