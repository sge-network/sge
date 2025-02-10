package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/ovm/types"
)

// getPubKeysChangeProposalStore returns pubkeys list change store ready for iterating.
func (k Keeper) getPubKeysChangeProposalStore(ctx sdk.Context) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PubKeysChangeProposalListPrefix)
	return store
}

// getKeyVaultStore returns key vault store ready for iterating.
func (k Keeper) getKeyVaultStore(ctx sdk.Context) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyVaultKey)
	return store
}

// getProposalStatStore returns proposal stats store ready for iterating.
func (k Keeper) getProposalStatStore(ctx sdk.Context) prefix.Store {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProposalStatsKey)
	return store
}
