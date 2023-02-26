package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/types"
)

// getActivePubKeysChangeProposalStore returns pubkeys list change store ready for iterating
func (k Keeper) getActivePubKeysChangeProposalStore(ctx sdk.Context) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PubKeysChangeProposalListActivePrefix)
	return store
}

// getFinishedPubKeysChangeProposalStore returns approved pubkeys list change store ready for iterating
func (k Keeper) getFinishedPubKeysChangeProposalStore(ctx sdk.Context) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.FinishedPubKeysChangeProposalListPrefix)
	return store
}

// getPubKeysStore returns pubkeys list store ready for iterating
func (k Keeper) getPubKeysStore(ctx sdk.Context) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PubKeysListKey)
	return store
}

// getProposalStatStore returns proposal stats store ready for iterating
func (k Keeper) getProposalStatStore(ctx sdk.Context) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalStatsKey)
	return store
}
