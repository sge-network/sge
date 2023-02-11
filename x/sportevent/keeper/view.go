package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/sportevent/types"
)

// getSportEventsStore gets the store containing all events.
func (k Keeper) getSportEventsStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.SportEventKeyPrefix)
}

// getSportEventStatsStore returns sport-event stats store ready for iterating.
func (k Keeper) getSportEventStatsStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.SportEventStatsKey)
}
