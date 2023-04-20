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
		2, // sr participation is made in two tranches to solve the reordering issue
		uint64(len(oddsUIDs)),
		types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE,
	)

	// Transfer sr contribution from sr to `sr_book_liquidity_pool` Account
	err = k.transferFundsFromModuleToModule(ctx, types.SRPoolName, types.OrderBookLiquidityName, srContribution)
	if err != nil {
		return
	}

	// Add book participations
	tranch1SRContribution := srContribution.Quo(sdk.NewInt(2))
	tranch2SRContribution := srContribution.Sub(tranch1SRContribution)
	fulfillmentQueue := []uint64{}
	for i := 0; i <= 1; i++ {
		var tranchSRContribution sdk.Int
		if i == 0 {
			tranchSRContribution = tranch1SRContribution
		} else {
			tranchSRContribution = tranch2SRContribution
		}
		srParticipation := types.NewOrderBookParticipation(
			types.SrparticipationIndex+uint64(i), orderBook.UID, k.accountKeeper.GetModuleAddress(types.SRPoolName).String(), orderBook.OddsCount, true, tranchSRContribution, tranchSRContribution,
			sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.Int{}, "", sdk.ZeroInt(),
		)
		_, found = k.GetOrderBookParticipation(ctx, orderBook.UID, srParticipation.Index)
		if found {
			err = sdkerrors.Wrapf(types.ErrOrderBookAlreadyPresent, "%d", srParticipation.Index)
			return
		}
		k.SetOrderBookParticipation(ctx, srParticipation)
		fulfillmentQueue = append(fulfillmentQueue, srParticipation.Index)

		// Add participation exposures
		for _, oddsUID := range oddsUIDs {
			pe := types.NewParticipationExposure(srParticipation.OrderBookUID, oddsUID, sdk.ZeroInt(), sdk.ZeroInt(), srParticipation.Index, types.RoundStart, false)
			k.SetParticipationExposure(ctx, pe)
		}
	}

	// Add book exposures
	for _, oddsUID := range oddsUIDs {
		boe := types.NewOrderBookOddsExposure(orderBook.UID, oddsUID, fulfillmentQueue)
		k.SetOrderBookOddsExposure(ctx, boe)
	}

	k.SetOrderBook(ctx, orderBook)

	return nil
}
