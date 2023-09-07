package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/reward/types"
)

// SetCampaign set a specific campaign in the store from its index
func (k Keeper) SetCampaign(ctx sdk.Context, campaign types.Campaign) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignKeyPrefix)
	b := k.cdc.MustMarshal(&campaign)
	store.Set(types.GetCampaignKey(
		campaign.UID,
	), b)
}

// GetCampaign returns a campaign from its index
func (k Keeper) GetCampaign(
	ctx sdk.Context,
	index string,

) (val types.Campaign, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignKeyPrefix)

	b := store.Get(types.GetCampaignKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCampaign removes a campaign from the store
func (k Keeper) RemoveCampaign(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignKeyPrefix)
	store.Delete(types.GetCampaignKey(
		index,
	))
}

// GetAllCampaign returns all campaign
func (k Keeper) GetAllCampaign(ctx sdk.Context) (list []types.Campaign) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.CampaignKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Campaign
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
