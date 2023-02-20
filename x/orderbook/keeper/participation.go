package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	housetypes "github.com/sge-network/sge/x/house/types"
	"github.com/sge-network/sge/x/orderbook/types"
)

// SetBookParticipation sets a book participation.
func (k Keeper) SetBookParticipation(ctx sdk.Context, participation types.BookParticipation) {
	participationKey := types.GetBookParticipationKey(participation.BookID, participation.Index)

	store := k.getParticipationStore(ctx)
	b := k.cdc.MustMarshal(&participation)
	store.Set(participationKey, b)
}

// GetBook GetBookParticipation a specific participation.
func (k Keeper) GetBookParticipation(ctx sdk.Context, bookID string, index uint64) (val types.BookParticipation, found bool) {
	store := k.getParticipationStore(ctx)
	aprticipationKey := types.GetBookParticipationKey(bookID, index)
	b := store.Get(aprticipationKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetParticipationsOfBook returns all participations for a book
func (k Keeper) GetParticipationsOfBook(ctx sdk.Context, bookID string) (list []types.BookParticipation, err error) {
	store := k.getParticipationStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.GetBookParticipationsKey(bookID))

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.BookParticipation
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllBookParticipations returns all book participations used during genesis dump.
func (k Keeper) GetAllBookParticipations(ctx sdk.Context) (list []types.BookParticipation, err error) {
	store := k.getParticipationStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.BookParticipation
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// InitiateBookParticipation starts a participation on a book for a certain account.
func (k Keeper) InitiateBookParticipation(
	ctx sdk.Context, addr sdk.AccAddress, bookID string, liquidity, fee sdk.Int,
) (index uint64, err error) {
	book, found := k.GetBook(ctx, bookID)
	if !found {
		err = sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", bookID)
		return
	}

	if book.Status != types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE {
		err = sdkerrors.Wrapf(types.ErrOrderBookNotActive, "%s", book.Status)
		return
	}

	// check if the maximum allowed participations is met or not.
	if k.GetMaxBookParticipationsAllowed(ctx) <= book.ParticipationCount {
		err = sdkerrors.Wrapf(types.ErrMaxNumberOfParticipationsReached, "%d", book.ParticipationCount)
		return
	}

	// calculate new sequential participation id
	book.ParticipationCount++
	index = book.ParticipationCount

	_, found = k.GetBookParticipation(ctx, book.ID, index)
	// This should never happen, just a sanity check
	if found {
		err = sdkerrors.Wrapf(types.ErrBookParticipationAlreadyExists, "id already exists %d", index)
		return
	}

	bookParticipation := types.NewBookParticipation(
		index, book.ID, addr.String(),
		book.OddsCount, // all of odds need to be filled in the next steps
		false,
		liquidity, liquidity, // int the start, liquidity and current round liquidity are the same
		sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.Int{}, "", sdk.ZeroInt(),
	)

	// Transfer fee from book participation to the feeAccountName
	err = k.transferFundsFromUserToModule(ctx, addr, housetypes.HouseParticipationFeeName, fee)
	if err != nil {
		return
	}

	// Transfer liquidity amount from book participation  to `book_liquidity_pool` Account
	err = k.transferFundsFromUserToModule(ctx, addr, types.BookLiquidityName, liquidity)
	if err != nil {
		return
	}

	// Make entry for book participation
	k.SetBookParticipation(ctx, bookParticipation)

	// Update book odds exposures and add particiapnt exposures
	boes, err := k.GetOddsExposuresByBook(ctx, bookParticipation.BookID)
	if err != nil {
		return
	}
	for _, boe := range boes {
		boe.FulfillmentQueue = append(boe.FulfillmentQueue, index)
		k.SetBookOddsExposure(ctx, boe)

		pe := types.NewParticipationExposure(book.ID, boe.OddsID, sdk.ZeroInt(), sdk.ZeroInt(), index, 1, false)
		k.SetParticipationExposure(ctx, pe)
	}

	// Update orderbook
	k.SetBook(ctx, book)

	return
}

func (k Keeper) LiquidateBookParticipation(
	ctx sdk.Context, depositorAddr, bookID string, participationIndex uint64, mode housetypes.WithdrawalMode, amount sdk.Int,
) (sdk.Int, error) {
	var withdrawalAmt sdk.Int
	depositorAddress, err := sdk.AccAddressFromBech32(depositorAddr)
	if err != nil {
		return withdrawalAmt, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidDesositor, err)
	}

	bp, found := k.GetBookParticipation(ctx, bookID, participationIndex)
	if !found {
		return withdrawalAmt, sdkerrors.Wrapf(types.ErrBookParticipationNotFound, "%s, %d", bookID, participationIndex)
	}

	if bp.IsSettled {
		return withdrawalAmt, sdkerrors.Wrapf(types.ErrBookParticipationAlreadySettled, "%s, %d", bookID, participationIndex)
	}

	if bp.ParticipantAddress != depositorAddr {
		return withdrawalAmt, sdkerrors.Wrapf(types.ErrMismatchInDepositorAddress, "%s", bp.ParticipantAddress)
	}

	if bp.IsModuleAccount {
		return withdrawalAmt, sdkerrors.Wrapf(types.ErrDepositorIsModuleAccount, "%s", bp.ParticipantAddress)
	}

	// Calculate max amount that can be transferred
	maxTransferableAmount := bp.CurrentRoundLiquidity.Sub(bp.CurrentRoundMaxLoss)

	switch mode {
	case housetypes.WithdrawalMode_WITHDRAWAL_MODE_FULL:
		if maxTransferableAmount.LTE(sdk.ZeroInt()) {
			return withdrawalAmt, sdkerrors.Wrapf(types.ErrMaxWithdrawableAmountIsZero, "%d, %d", bp.CurrentRoundLiquidity.Int64(), bp.CurrentRoundMaxLoss.Int64())
		}
		err := k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, depositorAddress, maxTransferableAmount)
		if err != nil {
			return withdrawalAmt, err
		}
		withdrawalAmt = maxTransferableAmount
	case housetypes.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL:
		if maxTransferableAmount.LT(amount) {
			return withdrawalAmt, sdkerrors.Wrapf(types.ErrWithdrawalAmountIsTooLarge, ": got %d, max %d", amount, maxTransferableAmount)
		}
		err := k.transferFundsFromModuleToUser(ctx, types.BookLiquidityName, depositorAddress, amount)
		if err != nil {
			return withdrawalAmt, err
		}
		withdrawalAmt = amount
	default:
		return withdrawalAmt, sdkerrors.Wrapf(housetypes.ErrInvalidMode, "%s", mode.String())
	}

	bp.CurrentRoundLiquidity = bp.CurrentRoundLiquidity.Sub(withdrawalAmt)
	bp.Liquidity = bp.Liquidity.Sub(withdrawalAmt)
	k.SetBookParticipation(ctx, bp)

	if bp.CurrentRoundLiquidity.LTE(sdk.ZeroInt()) {
		boes, err := k.GetOddsExposuresByBook(ctx, bookID)
		if err != nil {
			return withdrawalAmt, err
		}
		for _, boe := range boes {
			for i, pn := range boe.FulfillmentQueue {
				if pn == bp.Index {
					boe.FulfillmentQueue = append(boe.FulfillmentQueue[:i], boe.FulfillmentQueue[i+1:]...)
				}
			}
			k.SetBookOddsExposure(ctx, boe)
		}
	}

	return withdrawalAmt, nil
}
