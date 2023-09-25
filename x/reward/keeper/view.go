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
