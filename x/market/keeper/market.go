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

// GetMarkets returns all markets
func (k Keeper) GetMarkets(ctx sdk.Context) (list []types.Market, err error) {
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

// Resolve updates a market with its resolution
func (k Keeper) Resolve(
	ctx sdk.Context,
	storedMarket types.Market,
	resolutionMarket *types.MarketResolutionTicketPayload,
) *types.Market {
	storedMarket.ResolutionTS = resolutionMarket.ResolutionTS
	storedMarket.Status = resolutionMarket.Status

	// if the result is declared for the market, we need to update the winner odds uids.
	if resolutionMarket.Status == types.MarketStatus_MARKET_STATUS_RESULT_DECLARED {
		storedMarket.WinnerOddsUIDs = resolutionMarket.WinnerOddsUIDs
	}

	// if the result is declared or the market is canceled or aborted, it should be added
	// to the unsettled resolved market list  in the state.
	if storedMarket.IsResolved() {
		// append market id to the unsettled resolved in statistics.
		k.appendUnsettledResolvedMarket(ctx, storedMarket.UID)
	}

	k.SetMarket(ctx, storedMarket)

	return &storedMarket
}
