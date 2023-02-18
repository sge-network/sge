package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/types"
)

// GetBookOddsExposure returns a specific book odds exposure.
func (k Keeper) GetBookOddsExposure(ctx sdk.Context, bookID, oddsID string) (boe types.BookOddsExposure, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookOddsExposureKeyPrefix)
	value := store.Get(types.GetBookOddsExposureKey(bookID, oddsID))
	if value == nil {
		return boe, false
	}

	boe = types.MustUnmarshalBookOddsExposure(k.cdc, value)

	return boe, true
}

// GetOddsExposuresByBook returns all exposures for a book
func (k Keeper) GetOddsExposuresByBook(ctx sdk.Context, bookID string) (boes []types.BookOddsExposure) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookOddsExposureKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.GetBookOddsExposuresKey(bookID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		boe := types.MustUnmarshalBookOddsExposure(k.cdc, iterator.Value())
		boes = append(boes, boe)
	}

	return boes
}

// IterateAllBookExposures iterates through all of the book exposures.
func (k Keeper) IterateAllBookExposures(ctx sdk.Context, cb func(be types.BookOddsExposure) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.BookOddsExposureKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		be := types.MustUnmarshalBookOddsExposure(k.cdc, iterator.Value())
		if cb(be) {
			break
		}
	}
}

// GetAllBookExposures returns all book exposures used during genesis dump.
func (k Keeper) GetAllBookExposures(ctx sdk.Context) (bes []types.BookOddsExposure) {
	k.IterateAllBookExposures(ctx, func(be types.BookOddsExposure) bool {
		bes = append(bes, be)
		return false
	})

	return bes
}

// SetBookOddsExposure sets a book odds exposure.
func (k Keeper) SetBookOddsExposure(ctx sdk.Context, boe types.BookOddsExposure) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookOddsExposureKeyPrefix)
	store.Set(types.GetBookOddsExposureKey(boe.BookID, boe.OddsID), types.MustMarshalBookOddsExposure(k.cdc, boe))
}
