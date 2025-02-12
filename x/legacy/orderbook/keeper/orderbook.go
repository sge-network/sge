package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// SetOrderBook sets an order book.
func (k Keeper) SetOrderBook(ctx sdk.Context, book types.OrderBook) {
	bookKey := types.GetOrderBookKey(book.UID)

	store := k.getOrderBookStore(ctx)
	b := k.cdc.MustMarshal(&book)
	store.Set(bookKey, b)
}

// GetOrderBook returns a specific order book by its uid.
func (k Keeper) GetOrderBook(ctx sdk.Context, orderBookUID string) (val types.OrderBook, found bool) {
	marketsStore := k.getOrderBookStore(ctx)
	bookKey := types.GetOrderBookKey(orderBookUID)
	b := marketsStore.Get(bookKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetAllOrderBooks returns all order books used during genesis dump.
func (k Keeper) GetAllOrderBooks(ctx sdk.Context) (list []types.OrderBook, err error) {
	store := k.getOrderBookStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OrderBook
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// InitiateOrderBook initiates an order book for a given market.
func (k Keeper) InitiateOrderBook(ctx sdk.Context, marketUID string, oddsUIDs []string) (err error) {
	// book and market have one-to-one relationship
	orderBookUID := marketUID

	// check for existing orderBook with uid
	orderBook, found := k.GetOrderBook(ctx, orderBookUID)
	if found {
		return sdkerrors.Wrapf(types.ErrOrderBookAlreadyPresent, "%s", orderBook.UID)
	}

	// create new active book object
	orderBook = types.NewOrderBook(
		orderBookUID,
		uint64(len(oddsUIDs)),
		types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE,
	)

	// Add book exposures
	for _, oddsUID := range oddsUIDs {
		boe := types.NewOrderBookOddsExposure(orderBook.UID, oddsUID, []uint64{})
		k.SetOrderBookOddsExposure(ctx, boe)
	}

	k.SetOrderBook(ctx, orderBook)

	return nil
}
