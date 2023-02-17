package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/types"
)

// GetBookOddExposure returns a specific book odd exposure.
func (k Keeper) GetBookOddExposure(ctx sdk.Context, bookID, oddsID string) (boe types.BookOddExposure, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookOddExposureKeyPrefix)
	value := store.Get(types.GetBookOddExposureKey(bookID, oddsID))
	if value == nil {
		return boe, false
	}

	boe = types.MustUnmarshalBookOddExposure(k.cdc, value)

	return boe, true
}

// GetOddExposuresByBook returns all exposures for a book
func (k Keeper) GetOddExposuresByBook(ctx sdk.Context, bookID string) (boes []types.BookOddExposure) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookOddExposureKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.GetBookOddExposuresKey(bookID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		boe := types.MustUnmarshalBookOddExposure(k.cdc, iterator.Value())
		boes = append(boes, boe)
	}

	return boes
}

// IterateAllBookExposures iterates through all of the book exposures.
func (k Keeper) IterateAllBookExposures(ctx sdk.Context, cb func(be types.BookOddExposure) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.BookOddExposureKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		be := types.MustUnmarshalBookOddExposure(k.cdc, iterator.Value())
		if cb(be) {
			break
		}
	}
}

// GetAllBookExposures returns all book exposures used during genesis dump.
func (k Keeper) GetAllBookExposures(ctx sdk.Context) (bes []types.BookOddExposure) {
	k.IterateAllBookExposures(ctx, func(be types.BookOddExposure) bool {
		bes = append(bes, be)
		return false
	})

	return bes
}

// SetBookOddExposure sets a book odd exposure.
func (k Keeper) SetBookOddExposure(ctx sdk.Context, boe types.BookOddExposure) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookOddExposureKeyPrefix)
	store.Set(types.GetBookOddExposureKey(boe.BookID, boe.OddsID), types.MustMarshalBookOddExposure(k.cdc, boe))
}
