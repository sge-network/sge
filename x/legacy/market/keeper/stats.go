package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/legacy/market/types"
)

// SetMarketStats sets market statistics in the store
func (k Keeper) SetMarketStats(ctx sdk.Context, stats types.MarketStats) {
	store := k.getMarketStatsStore(ctx)
	b := k.cdc.MustMarshal(&stats)
	store.Set(utils.StrBytes("0"), b)
}

// GetMarketStats returns market stats
func (k Keeper) GetMarketStats(ctx sdk.Context) (val types.MarketStats) {
	store := k.getMarketStatsStore(ctx)

	b := store.Get(utils.StrBytes("0"))
	if b == nil {
		return val
	}

	k.cdc.MustUnmarshal(b, &val)
	return val
}

// appendUnsettledResolvedMarket appends market to the unsettled slice
func (k Keeper) appendUnsettledResolvedMarket(ctx sdk.Context, storedMarketUID string) {
	stats := k.GetMarketStats(ctx)
	stats.ResolvedUnsettled = append(stats.ResolvedUnsettled, storedMarketUID)
	k.SetMarketStats(ctx, stats)
}

// GetFirstUnsettledResolvedMarket returns first element of resolved
// markets that have active bets
func (k Keeper) GetFirstUnsettledResolvedMarket(ctx sdk.Context) (string, bool) {
	stats := k.GetMarketStats(ctx)
	if len(stats.ResolvedUnsettled) > 0 {
		return stats.ResolvedUnsettled[0], true
	}
	return "", false
}

// RemoveUnsettledResolvedMarket removes resolved market
// from the statistics
func (k Keeper) RemoveUnsettledResolvedMarket(ctx sdk.Context, marketUID string) {
	stats := k.GetMarketStats(ctx)
	if len(stats.ResolvedUnsettled) > 0 {
		for i, e := range stats.ResolvedUnsettled {
			if e == marketUID {
				stats.ResolvedUnsettled = append(
					stats.ResolvedUnsettled[:i],
					stats.ResolvedUnsettled[i+1:]...)
			}
		}
	}
	k.SetMarketStats(ctx, stats)
}
