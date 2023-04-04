package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/strategicreserve/types"
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
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

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
func (k Keeper) InitiateOrderBook(ctx sdk.Context, marketUID string, srContribution sdk.Int, oddsUIDs []string) (err error) {
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
		1, // sr participation is the only participation of the new book
		uint64(len(oddsUIDs)),
		types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE,
	)

	// Transfer sr contribution from sr to `sr_book_liquidity_pool` Account
	err = k.transferFundsFromModuleToModule(ctx, types.SRPoolName, types.OrderBookLiquidityName, srContribution)
	if err != nil {
		return
	}

	// Add book participation
	srParticipation := types.NewOrderBookParticipation(
		types.SrparticipationIndex, orderBook.UID, k.accountKeeper.GetModuleAddress(types.SRPoolName).String(), orderBook.OddsCount, true, srContribution, srContribution,
		sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.Int{}, "", sdk.ZeroInt(),
	)

	_, found = k.GetOrderBookParticipation(ctx, orderBook.UID, srParticipation.Index)
	if found {
		err = sdkerrors.Wrapf(types.ErrOrderBookAlreadyPresent, "%d", srParticipation.Index)
		return
	}
	k.SetOrderBookParticipation(ctx, srParticipation)

	// Add book exposures
	fulfillmentQueue := []uint64{srParticipation.Index}
	for _, oddsUID := range oddsUIDs {
		boe := types.NewOrderBookOddsExposure(orderBook.UID, oddsUID, fulfillmentQueue)
		k.SetOrderBookOddsExposure(ctx, boe)

		pe := types.NewParticipationExposure(srParticipation.OrderBookUID, oddsUID, sdk.ZeroInt(), sdk.ZeroInt(), srParticipation.Index, types.RoundStart, false)
		k.SetParticipationExposure(ctx, pe)
	}

	k.SetOrderBook(ctx, orderBook)

	return nil
}
