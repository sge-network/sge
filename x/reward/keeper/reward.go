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

// GetAllReward returns all reward
func (k Keeper) GetAllReward(ctx sdk.Context) (list []types.Reward) {
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

// SetOneTimeReward set a specific record to track if a reward is already given or not.
func (k Keeper) SetOneTimeReward(ctx sdk.Context, oneTimeReward types.OneTimeReward) {
	store := k.getRewardOneTimeStore(ctx)
	b := k.cdc.MustMarshal(&oneTimeReward)
	store.Set(types.GetOneTimeRewardsPrefix(oneTimeReward.RewardType, oneTimeReward.OneTimeKey), b)
}

// HasOneTimeReward returns if the reward for the reward type is already given.
func (k Keeper) HasOneTimeReward(
	ctx sdk.Context,
	oneTimeReward types.OneTimeReward,
) bool {
	store := k.getRewardOneTimeStore(ctx)
	return store.Has(types.GetOneTimeRewardsPrefix(oneTimeReward.RewardType, oneTimeReward.OneTimeKey))
}

// SetOneTimeReward set a specific record to track if a reward is already given or not.
func (k Keeper) SetRewardByReceiver(ctx sdk.Context, rewByType types.RewardByType) {
	store := k.getRewardByReceiverAndTypeStore(ctx)
	b := k.cdc.MustMarshal(&rewByType)
	store.Set(types.GetRewardsByTypeKey(rewByType.Addr, rewByType.RewardType, rewByType.UID), b)
}

// SetOneTimeReward set a specific record to track if a reward is already given or not.
func (k Keeper) SetRewardByCampaign(ctx sdk.Context, rewByCampaign types.RewardByCampaign) {
	store := k.getRewardsByCampaignStore(ctx)
	b := k.cdc.MustMarshal(&rewByCampaign)
	store.Set(types.GetRewardsByCampaignKey(rewByCampaign.CampaignUID, rewByCampaign.UID), b)
}
