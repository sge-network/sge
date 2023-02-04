package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/orderbook/types"
)

// GetBookParticipant returns a specific book participant.
func (k Keeper) GetBookParticipant(ctx sdk.Context, bookID string, bpNumber uint64) (bp types.BookParticipant, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookParticiapntKeyPrefix)
	value := store.Get(types.GetBookParticipantKey(bookID, bpNumber))
	if value == nil {
		return bp, false
	}

	bp = types.MustUnmarshalBookParticipant(k.cdc, value)

	return bp, true
}

// IterateAllBookParticipants iterates through all of the book participants.
func (k Keeper) IterateAllBookParticipants(ctx sdk.Context, cb func(bp types.BookParticipant) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.BookParticiapntKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		bp := types.MustUnmarshalBookParticipant(k.cdc, iterator.Value())
		if cb(bp) {
			break
		}
	}
}

// GetAllBookParticipants returns all book participants used during genesis dump.
func (k Keeper) GetAllBookParticipants(ctx sdk.Context) (bps []types.BookParticipant) {
	k.IterateAllBookParticipants(ctx, func(bp types.BookParticipant) bool {
		bps = append(bps, bp)
		return false
	})

	return bps
}

// SetBookParticipant sets a book participant.
func (k Keeper) SetBookParticipant(ctx sdk.Context, bp types.BookParticipant) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookParticiapntKeyPrefix)
	store.Set(types.GetBookParticipantKey(bp.BookId, bp.ParticipantNumber), types.MustMarshalBookParticipant(k.cdc, bp))
}

func (k Keeper) AddBookParticipant(
	ctx sdk.Context, addr sdk.AccAddress, bookID string, liquidity, fee sdk.Int, feeAccountName string,
) (uint64, error) {
	var bookParticipantNumber uint64
	book, found := k.GetBook(ctx, bookID)
	if !found {
		return bookParticipantNumber, sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", bookID)
	}

	if book.Status != types.OrderBookStatus_STATUS_ACTIVE {
		return bookParticipantNumber, sdkerrors.Wrapf(types.ErrOrderBookNotActive, "%s", book.Status)
	}

	if k.MaxBookParticipants(ctx) <= book.Participants {
		return bookParticipantNumber, sdkerrors.Wrapf(types.ErrMaxNumberOfParticipantsReached, "%d", book.Participants)
	}

	bookParticipant, found := k.GetBookParticipant(ctx, book.Id, book.Participants+1)
	// This should never happen, just a sanity check
	if found {
		return bookParticipantNumber, sdkerrors.Wrapf(
			types.ErrBookParticipantAlreadyExists, "id already exists %d", bookParticipant.ParticipantNumber,
		)
	} else {
		bookParticipant = types.NewBookParticipant(book.Id, addr, book.Participants+1, false)
	}

	// Transfer fee from book participant to the feeAccountName
	err := k.transferFundsFromUserToModule(ctx, addr, feeAccountName, fee)
	if err != nil {
		return bookParticipantNumber, err
	}

	// Transfer liquidity amount from book participant  to `book_liquidity_pool` Account
	err = k.transferFundsFromUserToModule(ctx, addr, types.BookLiquidityName, liquidity)
	if err != nil {
		return bookParticipantNumber, err
	}

	// Update orderbook
	book.Participants += 1
	k.SetBook(ctx, book)

	// Make entry for book participant
	k.SetBookParticipant(ctx, bookParticipant)
	return bookParticipant.ParticipantNumber, nil
}
