package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	housetypes "github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/orderbook/types"
)

// SetOrderBookParticipation sets a book participation.
func (k Keeper) SetOrderBookParticipation(ctx sdk.Context, participation types.OrderBookParticipation) {
	participationKey := types.GetOrderBookParticipationKey(
		participation.OrderBookUID,
		participation.Index,
	)

	store := k.getParticipationStore(ctx)
	b := k.cdc.MustMarshal(&participation)
	store.Set(participationKey, b)
}

// GetOrderBookParticipation returns a specific participation of an order book by index.
func (k Keeper) GetOrderBookParticipation(
	ctx sdk.Context,
	bookUID string,
	index uint64,
) (val types.OrderBookParticipation, found bool) {
	store := k.getParticipationStore(ctx)
	participationKey := types.GetOrderBookParticipationKey(bookUID, index)
	b := store.Get(participationKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetParticipationsOfOrderBook returns all participations for an order book.
func (k Keeper) GetParticipationsOfOrderBook(
	ctx sdk.Context,
	bookUID string,
) (list []types.OrderBookParticipation, err error) {
	store := k.getParticipationStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.GetOrderBookParticipationsKey(bookUID))

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OrderBookParticipation
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllOrderBookParticipations returns all book participations used during genesis dump.
func (k Keeper) GetAllOrderBookParticipations(
	ctx sdk.Context,
) (list []types.OrderBookParticipation, err error) {
	store := k.getParticipationStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OrderBookParticipation
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// InitiateOrderBookParticipation starts a participation on a book for a certain account.
func (k Keeper) InitiateOrderBookParticipation(
	ctx sdk.Context, addr sdk.AccAddress, bookUID string, depositAmount, feeAmount sdkmath.Int,
) (index uint64, err error) {
	market, found := k.marketKeeper.GetMarket(ctx, bookUID)
	if !found {
		err = sdkerrors.Wrapf(types.ErrMarketNotFound, "%s", bookUID)
		return
	}

	if market.Status != markettypes.MarketStatus_MARKET_STATUS_ACTIVE {
		err = sdkerrors.Wrapf(types.ErrParticipationOnInactiveMarket, "%s", bookUID)
		return
	}

	book, found := k.GetOrderBook(ctx, bookUID)
	if !found {
		err = sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", bookUID)
		return
	}

	if book.Status != types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE {
		err = sdkerrors.Wrapf(types.ErrOrderBookNotActive, "%s", book.Status)
		return
	}

	// check if the maximum allowed participations is met or not.
	if k.GetMaxOrderBookParticipationsAllowed(ctx) <= book.ParticipationCount {
		err = sdkerrors.Wrapf(types.ErrMaxNumberOfParticipationsReached, "%d", book.ParticipationCount)
		return
	}

	// calculate new sequential participation id
	book.ParticipationCount++
	index = book.ParticipationCount

	_, found = k.GetOrderBookParticipation(ctx, book.UID, index)
	// This should never happen, just a sanity check
	if found {
		err = sdkerrors.Wrapf(types.ErrBookParticipationAlreadyExists, "id already exists %d", index)
		return
	}

	liquidity := depositAmount.Sub(feeAmount)

	bookParticipation := types.NewOrderBookParticipation(
		index, book.UID, addr.String(),
		book.OddsCount,                  // all odds need to be filled in the next steps
		liquidity, feeAmount, liquidity, // int the start, liquidity and current round liquidity are the same
		sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdkmath.Int{}, "", sdk.ZeroInt(),
	)

	// fund order book liquidity pool from participant's account.
	if err = k.fund(types.OrderBookLiquidityFunder{}, ctx, addr, liquidity); err != nil {
		return
	}

	// fund house fee collector from participant's account.
	if err = k.fund(housetypes.HouseFeeCollectorFunder{}, ctx, addr, feeAmount); err != nil {
		return
	}

	// Make entry for book participation
	k.SetOrderBookParticipation(ctx, bookParticipation)

	// Update book odds exposures and add participant exposures
	if err = k.initParticipationExposures(ctx, book.UID, index); err != nil {
		return
	}

	// Update orderbook
	k.SetOrderBook(ctx, book)

	return index, nil
}

func (k Keeper) CalcWithdrawalAmount(
	ctx sdk.Context,
	depositorAddress string,
	marketUID string,
	participationIndex uint64,
	mode housetypes.WithdrawalMode,
	totalWithdrawnAmount sdkmath.Int,
	amount sdkmath.Int,
) (sdkmath.Int, error) {
	bp, found := k.GetOrderBookParticipation(ctx, marketUID, participationIndex)
	if !found {
		return sdkmath.Int{}, sdkerrors.Wrapf(
			types.ErrOrderBookParticipationNotFound,
			"%s, %d",
			marketUID,
			participationIndex,
		)
	}

	if bp.IsSettled {
		return sdkmath.Int{}, sdkerrors.Wrapf(
			types.ErrBookParticipationAlreadySettled,
			"%s, %d",
			bp.OrderBookUID,
			participationIndex,
		)
	}

	if bp.ParticipantAddress != depositorAddress {
		return sdkmath.Int{}, sdkerrors.Wrapf(types.ErrMismatchInDepositorAddress, "%s", bp.ParticipantAddress)
	}

	if mode == housetypes.WithdrawalMode_WITHDRAWAL_MODE_PARTIAL {
		if bp.Liquidity.Sub(totalWithdrawnAmount).LT(amount) {
			return sdkmath.Int{}, sdkerrors.Wrapf(types.ErrWithdrawalTooLarge, "%d", amount.Int64())
		}
	}

	withdrawAmount, err := bp.WithdrawableAmount(mode, amount)
	if err != nil {
		return sdkmath.Int{}, err
	}

	return withdrawAmount, nil
}

// WithdrawOrderBookParticipation withdraws the order book participation to the bettor's account
func (k Keeper) WithdrawOrderBookParticipation(
	ctx sdk.Context, marketUID string,
	participationIndex uint64,
	amount sdkmath.Int,
) error {
	bp, found := k.GetOrderBookParticipation(ctx, marketUID, participationIndex)
	if !found {
		return sdkerrors.Wrapf(
			types.ErrOrderBookParticipationNotFound,
			"%s, %d",
			marketUID,
			participationIndex,
		)
	}

	// refund participant's account from order book liquidity pool.
	if err := k.refund(types.OrderBookLiquidityFunder{}, ctx, sdk.MustAccAddressFromBech32(bp.ParticipantAddress), amount); err != nil {
		return err
	}

	bp.SetLiquidityAfterWithdrawal(amount)
	k.SetOrderBookParticipation(ctx, bp)

	return k.removeNotWithdrawableFromFulfillmentQueue(ctx, bp)
}
