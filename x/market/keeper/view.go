package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/market/types"
)

// getMarketsStore gets the store containing all markets.
func (k Keeper) getMarketsStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.MarketKeyPrefix)
}

// getMarketStatsStore returns market stats store ready for iterating.
func (k Keeper) getMarketStatsStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.MarketStatsKey)
}
