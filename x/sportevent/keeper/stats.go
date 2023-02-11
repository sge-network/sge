package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/sportevent/types"
)

// SetSportEventStats sets bet statistics in the store
func (k Keeper) SetSportEventStats(ctx sdk.Context, stats types.SportEventStats) {
	store := k.getSportEventStatsStore(ctx)
	b := k.cdc.MustMarshal(&stats)
	store.Set(utils.StrBytes("0"), b)
}

// GetSportEventStats returns sport-event stats
func (k Keeper) GetSportEventStats(ctx sdk.Context) (val types.SportEventStats) {
	store := k.getSportEventStatsStore(ctx)

	b := store.Get(utils.StrBytes("0"))
	if b == nil {
		return val
	}

	k.cdc.MustUnmarshal(b, &val)
	return val
}

// appendUnsettledResovedSportEvent appends sport-event to the unsettled slice
func (k Keeper) appendUnsettledResovedSportEvent(ctx sdk.Context, storedEventUID string) {
	stats := k.GetSportEventStats(ctx)
	stats.ResolvedUnsettled = append(stats.ResolvedUnsettled, storedEventUID)
	k.SetSportEventStats(ctx, stats)
}

// GetFirstUnsettledResovedSportEvent returns first element of resolved
// sport-events that have active bets
func (k Keeper) GetFirstUnsettledResovedSportEvent(ctx sdk.Context) (string, bool) {
	stats := k.GetSportEventStats(ctx)
	if len(stats.ResolvedUnsettled) > 0 {
		return stats.ResolvedUnsettled[0], true
	}
	return "", false
}

// RemoveUnsettledResolvedSportEvent removes resolved sport-event
// from the statistics
func (k Keeper) RemoveUnsettledResolvedSportEvent(ctx sdk.Context, sportEventUID string) {
	stats := k.GetSportEventStats(ctx)
	if len(stats.ResolvedUnsettled) > 0 {
		for i, e := range stats.ResolvedUnsettled {
			if e == sportEventUID {
				stats.ResolvedUnsettled = append(stats.ResolvedUnsettled[:i], stats.ResolvedUnsettled[i+1:]...)
			}
		}
	}
	k.SetSportEventStats(ctx, stats)
}
