package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/orderbook/types"
)

// SetOrderBookOddsExposure sets a book odds exposure.
func (k Keeper) SetOrderBookOddsExposure(ctx sdk.Context, boe types.OrderBookOddsExposure) {
	bookKey := types.GetOrderBookOddsExposureKey(boe.OrderBookUID, boe.OddsUID)

	store := k.getOrderBookOddsExposureStore(ctx)
	b := k.cdc.MustMarshal(&boe)
	store.Set(bookKey, b)
}

// GetOrderBookOddsExposure returns a specific book odds exposure.
func (k Keeper) GetOrderBookOddsExposure(
	ctx sdk.Context,
	orderBookUID, oddsUID string,
) (val types.OrderBookOddsExposure, found bool) {
	marketsStore := k.getOrderBookOddsExposureStore(ctx)
	exposureKey := types.GetOrderBookOddsExposureKey(orderBookUID, oddsUID)
	b := marketsStore.Get(exposureKey)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetOddsExposuresByOrderBook returns all exposures for an order book.
func (k Keeper) GetOddsExposuresByOrderBook(
	ctx sdk.Context,
	orderBookUID string,
) (list []types.OrderBookOddsExposure, err error) {
	store := k.getOrderBookOddsExposureStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.GetOrderBookOddsExposuresKey(orderBookUID))

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OrderBookOddsExposure
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllOrderBookExposures returns all order book exposures used during genesis dump.
func (k Keeper) GetAllOrderBookExposures(
	ctx sdk.Context,
) (list []types.OrderBookOddsExposure, err error) {
	store := k.getOrderBookOddsExposureStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OrderBookOddsExposure
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// initParticipationExposures initialize the odds and participation exposures for the
// participation at index.
func (k Keeper) initParticipationExposures(
	ctx sdk.Context,
	orderBookUID string,
	participationIndex uint64,
) error {
	// Update book odds exposures and add participant exposures
	boes, err := k.GetOddsExposuresByOrderBook(ctx, orderBookUID)
	bookEvent := types.NewOrderBookEvent()
	if err != nil {
		return err
	}
	for _, boe := range boes {
		boe.FulfillmentQueue = append(boe.FulfillmentQueue, participationIndex)
		k.SetOrderBookOddsExposure(ctx, boe)
		bookEvent.AddOrderBookOddsExposure(boe)

		pe := types.NewParticipationExposure(
			orderBookUID,
			boe.OddsUID,
			sdk.ZeroInt(),
			sdk.ZeroInt(),
			participationIndex,
			1,
			false,
		)
		k.SetParticipationExposure(ctx, pe)
		bookEvent.AddParticipationExposure(pe)
	}
	bookEvent.Emit(ctx)

	return nil
}

func (k Keeper) removeNotWithdrawableFromFulfillmentQueue(
	ctx sdk.Context,
	bp types.OrderBookParticipation,
) error {
	if !bp.IsLiquidityInCurrentRound() {
		boes, err := k.GetOddsExposuresByOrderBook(ctx, bp.OrderBookUID)
		if err != nil {
			return err
		}
		for _, boe := range boes {
			for i, pn := range boe.FulfillmentQueue {
				if pn == bp.Index {
					boe.FulfillmentQueue = append(
						boe.FulfillmentQueue[:i],
						boe.FulfillmentQueue[i+1:]...)
				}
			}
			k.SetOrderBookOddsExposure(ctx, boe)
		}
	}

	return nil
}
