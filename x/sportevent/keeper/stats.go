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
