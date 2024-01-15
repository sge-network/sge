package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/reward/types"
)

// SetReward set a specific reward in the store from its index
func (k Keeper) SetReward(ctx sdk.Context, reward types.Reward) {
	store := k.getRewardStore(ctx)
	b := k.cdc.MustMarshal(&reward)
	store.Set(types.GetRewardKey(reward.UID), b)
}

// GetReward returns a reward from its index
func (k Keeper) GetReward(
	ctx sdk.Context,
	uid string,
) (val types.Reward, found bool) {
	store := k.getRewardStore(ctx)

	b := store.Get(types.GetRewardKey(uid))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveReward removes a reward from the store
func (k Keeper) RemoveReward(
	ctx sdk.Context,
	uid string,
) {
	store := k.getRewardStore(ctx)
	store.Delete(types.GetRewardKey(uid))
}

// GetAllRewards returns all reward
func (k Keeper) GetAllRewards(ctx sdk.Context) (list []types.Reward) {
	store := k.getRewardStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Reward
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetRewardByReceiver(ctx sdk.Context, rByCategory types.RewardByCategory) {
	store := k.getRewardByReceiverAndCategoryStore(ctx)
	b := k.cdc.MustMarshal(&rByCategory)
	store.Set(types.GetRewardsByCategoryKey(rByCategory.Addr, rByCategory.RewardCategory, rByCategory.UID), b)
}

// GetRewardsByAddressAndCategory returns all rewards by address and category.
func (k Keeper) GetRewardsByAddressAndCategory(
	ctx sdk.Context,
	address string,
	category types.RewardCategory,
) (list []types.RewardByCategory, err error) {
	store := k.getRewardByReceiverAndCategoryStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.GetRewardsByCategoryPrefix(address, category))

	defer func() {
		err = iterator.Close()
	}()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RewardByCategory
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllRewardsByReceiverAndCategory returns all rewards by receiver and category
func (k Keeper) GetAllRewardsByReceiverAndCategory(ctx sdk.Context) (list []types.RewardByCategory) {
	store := k.getRewardByReceiverAndCategoryStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RewardByCategory
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) SetRewardByCampaign(ctx sdk.Context, rewByCampaign types.RewardByCampaign) {
	store := k.getRewardsByCampaignStore(ctx)
	b := k.cdc.MustMarshal(&rewByCampaign)
	store.Set(types.GetRewardsByCampaignKey(rewByCampaign.CampaignUID, rewByCampaign.UID), b)
}

// GetAllRewardsByCampaign returns all rewards by campaign
func (k Keeper) GetAllRewardsByCampaign(ctx sdk.Context) (list []types.RewardByCampaign) {
	store := k.getRewardsByCampaignStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.RewardByCampaign
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
