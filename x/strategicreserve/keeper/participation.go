package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	housetypes "github.com/sge-network/sge/x/house/types"
	markettypes "github.com/sge-network/sge/x/market/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// SetOrderBookParticipation sets a book participation.
func (k Keeper) SetOrderBookParticipation(ctx sdk.Context, participation types.OrderBookParticipation) {
	participationKey := types.GetOrderBookParticipationKey(participation.OrderBookUID, participation.Index)

	store := k.getParticipationStore(ctx)
	b := k.cdc.MustMarshal(&participation)
	store.Set(participationKey, b)
}

// GetOrderBookParticipation returns a specific participation of an order book by index.
func (k Keeper) GetOrderBookParticipation(ctx sdk.Context, bookUID string, index uint64) (val types.OrderBookParticipation, found bool) {
	store := k.getParticipationStore(ctx)
	aprticipationKey := types.GetOrderBookParticipationKey(bookUID, index)
	b := store.Get(aprticipationKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetParticipationsOfOrderBook returns all participations for an order book.
func (k Keeper) GetParticipationsOfOrderBook(ctx sdk.Context, bookUID string) (list []types.OrderBookParticipation, err error) {
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
func (k Keeper) GetAllOrderBookParticipations(ctx sdk.Context) (list []types.OrderBookParticipation, err error) {
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
	ctx sdk.Context, addr sdk.AccAddress, bookUID string, liquidity, fee sdk.Int,
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

	bookParticipation := types.NewOrderBookParticipation(
		index, book.UID, addr.String(),
		book.OddsCount,       // all of odds need to be filled in the next steps
		liquidity, liquidity, // int the start, liquidity and current round liquidity are the same
		sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt(), sdk.Int{}, "", sdk.ZeroInt(),
	)

	// Transfer liquidity amount from book participation  to `house_deposit_collector` Account
	err = k.transferFundsFromAccountToModule(ctx, addr, types.OrderBookLiquidityPool, liquidity)
	if err != nil {
		return
	}

	// Make entry for book participation
	k.SetOrderBookParticipation(ctx, bookParticipation)

	// Update book odds exposures and add particiapnt exposures
	k.initParticipationExposures(ctx, book.UID, index)

	// Update strategicreserve
	k.SetOrderBook(ctx, book)

	return index, nil
}

// initParticipationExposures initialize the odds and participation exposures for the
// participation at index.
func (k Keeper) initParticipationExposures(ctx sdk.Context, orderBookUID string, participationIndex uint64) error {
	// Update book odds exposures and add particiapnt exposures
	boes, err := k.GetOddsExposuresByOrderBook(ctx, orderBookUID)
	if err != nil {
		return err
	}
	for _, boe := range boes {
		boe.FulfillmentQueue = append(boe.FulfillmentQueue, participationIndex)
		k.SetOrderBookOddsExposure(ctx, boe)

		pe := types.NewParticipationExposure(orderBookUID, boe.OddsUID, sdk.ZeroInt(), sdk.ZeroInt(), participationIndex, 1, false)
		k.SetParticipationExposure(ctx, pe)
	}

	return nil
}

// WithdrawOrderBookParticipation withdraws the order book participation to the bettor's account
func (k Keeper) WithdrawOrderBookParticipation(
	ctx sdk.Context, depositorAddr, bookUID string, participationIndex uint64, mode housetypes.WithdrawalMode, amount sdk.Int,
) (sdk.Int, error) {
	depositorAddress, err := sdk.AccAddressFromBech32(depositorAddr)
	if err != nil {
		return sdk.Int{}, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidDesositor, err)
	}

	bp, found := k.GetOrderBookParticipation(ctx, bookUID, participationIndex)
	if !found {
		return sdk.Int{}, sdkerrors.Wrapf(types.ErrOrderBookParticipationNotFound, "%s, %d", bookUID, participationIndex)
	}

	if err = bp.ValidateWithdraw(depositorAddr, participationIndex); err != nil {
		return sdk.Int{}, err
	}

	withdrawalAmt, err := bp.WithdrawableAmount(depositorAddr, mode, amount)
	if err != nil {
		return sdk.Int{}, err
	}

	if err = k.transferFundsFromModuleToAccount(ctx, types.OrderBookLiquidityPool, depositorAddress, withdrawalAmt); err != nil {
		return sdk.Int{}, err
	}

	bp.SetLiquidityAfterWithdrawal(withdrawalAmt)
	k.SetOrderBookParticipation(ctx, bp)

	if err = k.removeNotWithdrawableFromFulfillmentQueue(ctx, bp); err != nil {
		return sdk.Int{}, err
	}

	return withdrawalAmt, nil
}

func (k Keeper) removeNotWithdrawableFromFulfillmentQueue(ctx sdk.Context, bp types.OrderBookParticipation) error {
	if !bp.IsWithdrawable() {
		boes, err := k.GetOddsExposuresByOrderBook(ctx, bp.OrderBookUID)
		if err != nil {
			return err
		}
		for _, boe := range boes {
			for i, pn := range boe.FulfillmentQueue {
				if pn == bp.Index {
					boe.FulfillmentQueue = append(boe.FulfillmentQueue[:i], boe.FulfillmentQueue[i+1:]...)
				}
			}
			k.SetOrderBookOddsExposure(ctx, boe)
		}
	}

	return nil
}
