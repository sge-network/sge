package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/reward/types"
)

// getCampaignStore gets the store containing all campaigns.
func (k Keeper) getCampaignStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.CampaignKeyPrefix)
}

// getRewardStore gets the store containing all rewards.
func (k Keeper) getRewardStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.RewardKeyPrefix)
}

// getRewardOneTimeStore gets the store containing all onetime rewards.
func (k Keeper) getRewardOneTimeStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.RewardOneTimeKeyPrefix)
}

// getRewardByReceiverAndTypeStore gets the store containing all rewards by receiver.
func (k Keeper) getRewardByReceiverAndTypeStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.RewardByReceiverAndTypeKeyPrefix)
}

// getRewardsByCampaignStore gets the store containing all rewards by campaign uid.
func (k Keeper) getRewardsByCampaignStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.RewardByCampaignKeyPrefix)
}
