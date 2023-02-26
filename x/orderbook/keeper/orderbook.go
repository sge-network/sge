package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/orderbook/types"
	srtypes "github.com/sge-network/sge/x/strategicreserve/types"
)

// SetBook sets a book.
func (k Keeper) SetBook(ctx sdk.Context, book types.OrderBook) {
	bookKey := types.GetBookKey(book.ID)

	store := k.getBookStore(ctx)
	b := k.cdc.MustMarshal(&book)
	store.Set(bookKey, b)
}

// GetBook returns a specific order book.
func (k Keeper) GetBook(ctx sdk.Context, bookID string) (val types.OrderBook, found bool) {
	sportEventsStore := k.getBookStore(ctx)
	bookKey := types.GetBookKey(bookID)
	b := sportEventsStore.Get(bookKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetAllBooks returns all order books used during genesis dump.
func (k Keeper) GetAllBooks(ctx sdk.Context) (list []types.OrderBook, err error) {
	store := k.getBookStore(ctx)
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

// InitiateBook initiates a book for a given sport event
func (k Keeper) InitiateBook(ctx sdk.Context, sportEventUID string, srContribution sdk.Int, oddsUIDs []string) (err error) {
	// book and sport event have one-to-one relationship
	bookID := sportEventUID

	// check for existing book with id
	book, found := k.GetBook(ctx, bookID)
	if found {
		return sdkerrors.Wrapf(types.ErrOrderBookAlreadyPresent, "%s", book.ID)
	}

	// create new active book object
	book = types.NewOrderBook(
		bookID,
		1, // sr participation is the only participation of the new book
		uint64(len(oddsUIDs)),
		types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE,
	)

	// Transfer sr contribution from sr to `sr_book_liquidity_pool` Account
	err = k.transferFundsFromModuleToModule(ctx, srtypes.SRPoolName, types.BookLiquidityName, srContribution)
	if err != nil {
		return
	}

	// Add book participation
	srParticipation := types.NewBookParticipation(
		types.SrparticipationIndex, book.ID, k.accountKeeper.GetModuleAddress(srtypes.SRPoolName).String(), book.OddsCount, true, srContribution, srContribution,
		sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.Int{}, "", sdk.ZeroInt(),
	)

	_, found = k.GetBookParticipation(ctx, book.ID, srParticipation.Index)
	if found {
		err = sdkerrors.Wrapf(types.ErrOrderBookAlreadyPresent, "%d", srParticipation.Index)
		return
	}
	k.SetBookParticipation(ctx, srParticipation)

	// Add book exposures
	fulfillmentQueue := []uint64{srParticipation.Index}
	for _, oddsUID := range oddsUIDs {
		boe := types.NewBookOddsExposure(book.ID, oddsUID, fulfillmentQueue)
		k.SetBookOddsExposure(ctx, boe)

		pe := types.NewParticipationExposure(srParticipation.BookID, oddsUID, sdk.ZeroInt(), sdk.ZeroInt(), srParticipation.Index, types.RoundStart, false)
		k.SetParticipationExposure(ctx, pe)
	}

	k.SetBook(ctx, book)

	return nil
}
