package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/strategicreserve/types"
)

// SetOrderBookOddsExposure sets a book odds exposure.
func (k Keeper) SetOrderBookOddsExposure(ctx sdk.Context, boe types.OrderBookOddsExposure) {
	bookKey := types.GetOrderBookOddsExposureKey(boe.OrderBookUID, boe.OddsUID)

	store := k.getOrderBookOddsExposureStore(ctx)
	b := k.cdc.MustMarshal(&boe)
	store.Set(bookKey, b)
}

// GetOrderBookOddsExposure returns a specific book odds exposure.
func (k Keeper) GetOrderBookOddsExposure(
	ctx sdk.Context,
	bookUID, oddsUID string,
) (val types.OrderBookOddsExposure, found bool) {
	marketsStore := k.getOrderBookOddsExposureStore(ctx)
	exposureKey := types.GetOrderBookOddsExposureKey(bookUID, oddsUID)
	b := marketsStore.Get(exposureKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetOddsExposuresByOrderBook returns all exposures for an order book.
func (k Keeper) GetOddsExposuresByOrderBook(
	ctx sdk.Context,
	bookUID string,
) (list []types.OrderBookOddsExposure, err error) {
	store := k.getOrderBookOddsExposureStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.GetOrderBookOddsExposuresKey(bookUID))

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OrderBookOddsExposure
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllOrderBookExposures returns all order book exposures used during genesis dump.
func (k Keeper) GetAllOrderBookExposures(
	ctx sdk.Context,
) (list []types.OrderBookOddsExposure, err error) {
	store := k.getOrderBookOddsExposureStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OrderBookOddsExposure
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
