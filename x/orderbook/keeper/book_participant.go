package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	htypes "github.com/sge-network/sge/x/house/types"
	"github.com/sge-network/sge/x/orderbook/types"
)

// GetBookParticipant returns a specific book participant.
func (k Keeper) GetBookParticipant(ctx sdk.Context, bookID string, bpNumber uint64) (bp types.BookParticipant, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookParticipantKeyPrefix)
	value := store.Get(types.GetBookParticipantKey(bookID, bpNumber))
	if value == nil {
		return bp, false
	}

	bp = types.MustUnmarshalBookParticipant(k.cdc, value)

	return bp, true
}

// GetParticipantsByBook returns all participants for a book
func (k Keeper) GetParticipantsByBook(ctx sdk.Context, bookID string) (bps []types.BookParticipant) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookParticipantKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, types.GetBookParticipantsKey(bookID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		bp := types.MustUnmarshalBookParticipant(k.cdc, iterator.Value())
		bps = append(bps, bp)
	}

	return bps
}

// IterateAllBookParticipants iterates through all of the book participants.
func (k Keeper) IterateAllBookParticipants(ctx sdk.Context, cb func(bp types.BookParticipant) (stop bool)) {
	iterator := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), types.BookParticipantKeyPrefix)
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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookParticipantKeyPrefix)
	store.Set(types.GetBookParticipantKey(bp.BookID, bp.ParticipantNumber), types.MustMarshalBookParticipant(k.cdc, bp))
}

func (k Keeper) AddBookParticipant(
	ctx sdk.Context, addr sdk.AccAddress, bookID string, liquidity, fee sdk.Int, feeAccountName string,
) (uint64, error) {
	var bookParticipantNumber uint64
	book, found := k.GetBook(ctx, bookID)
	if !found {
		return bookParticipantNumber, sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", bookID)
	}

	if book.Status != types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE {
		return bookParticipantNumber, sdkerrors.Wrapf(types.ErrOrderBookNotActive, "%s", book.Status)
	}

	if k.MaxBookParticipants(ctx) <= book.Participants {
		return bookParticipantNumber, sdkerrors.Wrapf(types.ErrMaxNumberOfParticipantsReached, "%d", book.Participants)
	}

	bookParticipant, found := k.GetBookParticipant(ctx, book.ID, book.Participants+1)
	// This should never happen, just a sanity check
	if found {
		return bookParticipantNumber, sdkerrors.Wrapf(
			types.ErrBookParticipantAlreadyExists, "id already exists %d", bookParticipant.ParticipantNumber,
		)
	}

	bookParticipant = types.NewBookParticipant(
		book.ID, addr, book.Participants+1, book.NumberOfOdds, false, liquidity, liquidity, sdk.ZeroInt(), sdk.ZeroInt(),
		sdk.ZeroInt(), sdk.Int{}, "", sdk.ZeroInt(),
	)

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

	// Make entry for book participant
	k.SetBookParticipant(ctx, bookParticipant)

	// Update book odd exposures and add particiapnt exposures
	boes := k.GetOddExposuresByBook(ctx, bookParticipant.BookID)
	for _, boe := range boes {
		boe.FullfillmentQueue = append(boe.FullfillmentQueue, bookParticipant.ParticipantNumber)
		k.SetBookOddExposure(ctx, boe)

		pe := types.NewParticipantExposure(book.ID, boe.OddsID, sdk.ZeroInt(), sdk.ZeroInt(), bookParticipant.ParticipantNumber, 1, false)
		k.SetParticipantExposure(ctx, pe)
	}

	// Update orderbook
	book.Participants++
	k.SetBook(ctx, book)

	return bookParticipant.ParticipantNumber, nil
}

func (k Keeper) LiquidateBookParticipant(
	ctx sdk.Context, depAddr, bookID string, bpNumber uint64, mode htypes.WithdrawalMode, amount sdk.Int,
) (sdk.Int, error) {
	var withdrawalAmt sdk.Int
	depositorAddress, err := sdk.AccAddressFromBech32(depAddr)
	if err != nil {
		return withdrawalAmt, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidDesositor, err)
	}

	bp, found := k.GetBookParticipant(ctx, bookID, bpNumber)
	if !found {
		return withdrawalAmt, sdkerrors.Wrapf(types.ErrBookParticipantNotFound, "%s, %d", bookID, bpNumber)
	}

	if bp.IsSettled {
		return withdrawalAmt, sdkerrors.Wrapf(types.ErrBookParticipantAlreadySettled, "%s, %d", bookID, bpNumber)
	}

	if bp.ParticipantAddress != depAddr {
		return withdrawalAmt, sdkerrors.Wrapf(types.ErrMismatchInDepositorAddress, "%s", bp.ParticipantAddress)
	}

	if bp.IsModuleAccount {
		return withdrawalAmt, sdkerrors.Wrapf(types.ErrDepositorIsModuleAccount, "%s", bp.ParticipantAddress)
	}

	// Calculate max amount that can be transferred
	maxTransferableAmount := bp.CurrentRoundLiquidity.Sub(bp.CurrentRoundMaxLoss)

	switch mode {
	case htypes.WithdrawalMode_WITHDRAWAL_MODE_FULL:
		if maxTransferableAmount.LTE(sdk.ZeroInt()) {
			return withdrawalAmt, sdkerrors.Wrapf(types.ErrMaxWithdrawableAmountIsZero, "%d, %d", bp.CurrentRoundLiquidity.Int64(), bp.CurrentRoundMaxLoss.Int64())
		}
		err := k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, depositorAddress, maxTransferableAmount)
		if err != nil {
			return withdrawalAmt, err
		}
		withdrawalAmt = maxTransferableAmount
	case htypes.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL:
		if maxTransferableAmount.LT(amount) {
			return withdrawalAmt, sdkerrors.Wrapf(types.ErrWithdrawalAmountIsTooLarge, ": got %d, max %d", amount, maxTransferableAmount)
		}
		err := k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, depositorAddress, amount)
		if err != nil {
			return withdrawalAmt, err
		}
		withdrawalAmt = amount
	default:
		return withdrawalAmt, sdkerrors.Wrapf(htypes.ErrInvalidMode, "%s", mode.String())
	}

	bp.CurrentRoundLiquidity = bp.CurrentRoundLiquidity.Sub(withdrawalAmt)
	bp.Liquidity = bp.Liquidity.Sub(withdrawalAmt)
	k.SetBookParticipant(ctx, bp)

	if bp.CurrentRoundLiquidity.LTE(sdk.ZeroInt()) {
		boes := k.GetOddExposuresByBook(ctx, bookID)
		for _, boe := range boes {
			for i, pn := range boe.FullfillmentQueue {
				if pn == bp.ParticipantNumber {
					boe.FullfillmentQueue = append(boe.FullfillmentQueue[:i], boe.FullfillmentQueue[i+1:]...)
				}
			}
			k.SetBookOddExposure(ctx, boe)
		}
	}

	return withdrawalAmt, nil
}
