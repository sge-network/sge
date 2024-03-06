package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/utils"
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

func (k Keeper) SetRewardOfReceiverByPromoterAndCategory(ctx sdk.Context, promoterUID string, rByCategory types.RewardByCategory) {
	store := k.getRewardByReceiverAndCategoryStore(ctx)
	b := k.cdc.MustMarshal(&rByCategory)
	store.Set(types.GetRewardsOfReceiverByPromoterAndCategoryKey(promoterUID, rByCategory.Addr, rByCategory.RewardCategory, rByCategory.UID), b)
}

// GetRewardsOfReceiverByPromoterAndCategory returns all rewards by address and category.
func (k Keeper) GetRewardsOfReceiverByPromoterAndCategory(
	ctx sdk.Context,
	promoterUID, addr string,
	category types.RewardCategory,
) (list []types.RewardByCategory, err error) {
	store := k.getRewardByReceiverAndCategoryStore(ctx)
	iterator := sdk.KVStorePrefixIterator(store, types.GetRewardsOfReceiverByPromoterAndCategoryPrefix(promoterUID, addr, category))

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

// HasRewardByReceiver returns true if there is a record for the category and reward receiver
func (k Keeper) HasRewardOfReceiverByPromoter(ctx sdk.Context, promoterUID, addr string, category types.RewardCategory) bool {
	rewardsByCat, err := k.GetRewardsOfReceiverByPromoterAndCategory(ctx, promoterUID, addr, category)
	if err != nil || len(rewardsByCat) > 0 {
		return true
	}
	return false
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

func (k Keeper) SetRewardGrantsStats(ctx sdk.Context, campaignUID, accAddr string, count uint64) {
	store := k.getRewardGrantsStatStore(ctx)
	b := utils.Uint64ToBytes(count)
	store.Set(types.GetRewardGrantStatKey(campaignUID, accAddr), b)
}

func (k Keeper) GetRewardGrantsStats(ctx sdk.Context, campaignUID, accAddr string) (val uint64, found bool) {
	store := k.getRewardGrantsStatStore(ctx)
	b := store.Get(types.GetRewardGrantStatKey(campaignUID, accAddr))
	if b == nil {
		return val, false
	}

	val = utils.Uint64FromBytes(b)
	return val, true
}
