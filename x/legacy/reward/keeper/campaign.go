package keeper

import (
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/reward/types"
)

// SetCampaign set a specific campaign in the store from its index
func (k Keeper) SetCampaign(ctx sdk.Context, campaign types.Campaign) {
	store := k.getCampaignStore(ctx)
	b := k.cdc.MustMarshal(&campaign)
	store.Set(types.GetCampaignKey(campaign.UID), b)
}

// GetCampaign returns a campaign from its index
func (k Keeper) GetCampaign(
	ctx sdk.Context,
	uid string,
) (val types.Campaign, found bool) {
	store := k.getCampaignStore(ctx)

	b := store.Get(types.GetCampaignKey(uid))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCampaign removes a campaign from the store
func (k Keeper) RemoveCampaign(
	ctx sdk.Context,
	uid string,
) {
	store := k.getCampaignStore(ctx)
	store.Delete(types.GetCampaignKey(uid))
}

// GetAllCampaign returns all campaign
func (k Keeper) GetAllCampaign(ctx sdk.Context) (list []types.Campaign) {
	store := k.getCampaignStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Campaign
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// UpdateCampaignPool updates campaign pool according to the total receiving amount of receiver
func (k Keeper) UpdateCampaignPool(ctx sdk.Context, campaign types.Campaign, receiver types.Receiver) {
	campaign.Pool.Spend(receiver.TotalAmount())
	k.SetCampaign(ctx, campaign)
}
