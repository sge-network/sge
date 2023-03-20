package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/market/types"
)

// SetMarket sets a specific market in the store
func (k Keeper) SetMarket(ctx sdk.Context, market types.Market) {
	store := k.getMarketsStore(ctx)
	b := k.cdc.MustMarshal(&market)
	store.Set(utils.StrBytes(market.UID), b)
}

// GetMarket returns a specific market by its UID
func (k Keeper) GetMarket(ctx sdk.Context, marketUID string) (val types.Market, found bool) {
	marketsStore := k.getMarketsStore(ctx)
	b := marketsStore.Get(utils.StrBytes(marketUID))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// MarketExists checks if a specific market exists or not
func (k Keeper) MarketExists(ctx sdk.Context, marketUID string) bool {
	_, found := k.GetMarket(ctx, marketUID)
	return found
}

// RemoveMarket removes a market from the store
func (k Keeper) RemoveMarket(ctx sdk.Context, marketUID string) {
	store := k.getMarketsStore(ctx)
	store.Delete(utils.StrBytes(marketUID))
}

// GetMarketAll returns all markets
func (k Keeper) GetMarketAll(ctx sdk.Context) (list []types.Market, err error) {
	store := k.getMarketsStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Market
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ResolveMarket updates a market with its resolution
func (k Keeper) ResolveMarket(ctx sdk.Context, resolutionEvent *types.MarketResolutionTicketPayload) (*types.Market, error) {
	storedEvent, found := k.GetMarket(ctx, resolutionEvent.UID)
	if !found {
		return nil, types.ErrNoMatchingMarket
	}

	if storedEvent.Status != types.MarketStatus_MARKET_STATUS_ACTIVE &&
		storedEvent.Status != types.MarketStatus_MARKET_STATUS_INACTIVE {
		return nil, types.ErrCanNotBeAltered
	}

	storedEvent.ResolutionTS = resolutionEvent.ResolutionTS
	storedEvent.Status = resolutionEvent.Status

	// if the result is declared for the market, we need to update the winner odds uids.
	if resolutionEvent.Status == types.MarketStatus_MARKET_STATUS_RESULT_DECLARED {
		storedEvent.WinnerOddsUIDs = resolutionEvent.WinnerOddsUIDs
	}

	// if the result is declared or the market is canceled or aborted, it should be added
	// to the unsettled resolved market list  in the state.
	if resolutionEvent.Status == types.MarketStatus_MARKET_STATUS_RESULT_DECLARED ||
		resolutionEvent.Status == types.MarketStatus_MARKET_STATUS_CANCELED ||
		resolutionEvent.Status == types.MarketStatus_MARKET_STATUS_ABORTED {
		// append market id to the unsettled resolved in statistics.
		k.appendUnsettledResolvedMarket(ctx, storedEvent.UID)
	}

	k.SetMarket(ctx, storedEvent)

	return &storedEvent, nil
}

func emitTransactionEvent(ctx sdk.Context, emitType string, uid, bookUID, creator string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			emitType,
			sdk.NewAttribute(types.AttributeKeyMarketsSuccessUID, uid),
			sdk.NewAttribute(types.AttributeKeyOrderBookUID, bookUID),
			sdk.NewAttribute(types.AttributeKeyEventsCreator, creator),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, emitType),
			sdk.NewAttribute(sdk.AttributeKeySender, creator),
		),
	})
}
