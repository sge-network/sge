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

// getRewardByReceiverAndTypeStore gets the store containing all rewards by receiver.
func (k Keeper) getRewardByReceiverAndCategoryStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.RewardByReceiverAndCategoryKeyPrefix)
}

// getRewardsByCampaignStore gets the store containing all rewards by campaign uid.
func (k Keeper) getRewardsByCampaignStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.RewardByCampaignKeyPrefix)
}

// getPromoterByAddressStore gets the store containing all promoter by address.
func (k Keeper) getPromoterByAddressStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.PromoterAddressKeyPrefix)
}

// getPromoterStore gets the store containing all promoters.
func (k Keeper) getPromoterStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.PromoterKeyPrefix)
}

// getRewardGrantsStatStore gets the store containing all reward grant stats.
func (k Keeper) getRewardGrantsStatStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.RewardGrantStatKeyPrefix)
}
