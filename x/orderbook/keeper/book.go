package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/orderbook/types"
	srtypes "github.com/sge-network/sge/x/strategicreserve/types"
)

// GetBook returns a specific book.
func (k Keeper) GetBook(ctx sdk.Context, bookID string) (book types.OrderBook, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookKeyPrefix)
	value := store.Get(types.GetBookKey(bookID))
	if value == nil {
		return book, false
	}

	book = types.MustUnmarshalBook(k.cdc, value)

	return book, true
}

// IterateAllBooks iterates through all of the books.
func (k Keeper) IterateAllBooks(ctx sdk.Context, cb func(book types.OrderBook) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.BookKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		book := types.MustUnmarshalBook(k.cdc, iterator.Value())
		if cb(book) {
			break
		}
	}
}

// GetAllBooks returns all books used during genesis dump.
func (k Keeper) GetAllBooks(ctx sdk.Context) (books []types.OrderBook) {
	k.IterateAllBooks(ctx, func(book types.OrderBook) bool {
		books = append(books, book)
		return false
	})

	return books
}

// SetBook sets a book.
func (k Keeper) SetBook(ctx sdk.Context, book types.OrderBook) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookKeyPrefix)
	store.Set(types.GetBookKey(book.Id), types.MustMarshalBook(k.cdc, book))
}

// InitiateBook initiates a book for a given sport event
func (k Keeper) InitiateBook(ctx sdk.Context, sportEventUid string, srContribution sdk.Int, oddIds []string) (string, error) {
	book, found := k.GetBook(ctx, sportEventUid)
	if found {
		return "", sdkerrors.Wrapf(types.ErrOrderBookAlreadyPresent, "%s", book.Id)
	} else {
		book = types.NewBook(sportEventUid, 0, uint64(len(oddIds)), types.OrderBookStatus_STATUS_ACTIVE)
	}

	// Transfer sr contribution from sr to `sr_book_liquidity_pool` Account
	err := k.transferFundsFromModuleToModule(ctx, srtypes.SRPoolName, types.SRBookLiquidityName, srContribution)
	if err != nil {
		return "", err
	}

	// Add book participant
	bp := types.NewBookParticipant(
		book.Id, k.accountKeeper.GetModuleAddress(srtypes.SRPoolName), 1, book.NumberOfOdds, true, srContribution, srContribution,
		sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.Int{}, "", []string{}, sdk.ZeroInt(),
	)
	_, found = k.GetBookParticipant(ctx, book.Id, bp.ParticipantNumber)
	if found {
		return "", sdkerrors.Wrapf(types.ErrOrderBookAlreadyPresent, "%s", bp.ParticipantNumber)
	}
	k.SetBookParticipant(ctx, bp)

	// Add book exposures
	fullfillmentQueue := []uint64{bp.ParticipantNumber}
	for _, oddId := range oddIds {
		boe := types.NewBookOddExposure(book.Id, oddId, fullfillmentQueue)
		k.SetBookOddExposure(ctx, boe)

		pe := types.NewParticipantExposure(bp.BookId, oddId, sdk.ZeroInt(), sdk.ZeroInt(), bp.ParticipantNumber, 1, false)
		k.SetParticipantExposure(ctx, pe)
	}

	// Make entry for book
	book.Participants = 1
	k.SetBook(ctx, book)

	return book.Id, nil
}
