package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
)

// getBetStore returns bet store ready for iterating
func (k Keeper) getBetStore(ctx sdk.Context) prefix.Store {
	betStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BetListPrefix)
	return betStore
}

// getBetIDStore returns bet id store ready for iterating
func (k Keeper) getBetIDStore(ctx sdk.Context) prefix.Store {
	betStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BetIDListPrefix)
	return betStore
}

// getBetStatsStore returns bet stats store ready for iterating
func (k Keeper) getBetStatsStore(ctx sdk.Context) prefix.Store {
	betStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BetStatsKey)
	return betStore
}

// getActiveStore returns active bet store ready for iterating
func (k Keeper) getActiveStore(ctx sdk.Context) prefix.Store {
	betStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActiveBetListPrefix)
	return betStore
}

// getSettledStore returns settled bet store ready for iterating
func (k Keeper) getSettledStore(ctx sdk.Context) prefix.Store {
	betStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.SettledBetListPrefix)
	return betStore
}
