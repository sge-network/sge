package reward

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/reward/keeper"
	"github.com/sge-network/sge/x/reward/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the campaign
	for _, elem := range genState.CampaignList {
		k.SetCampaign(ctx, elem)
	}
	for _, elem := range genState.RewardList {
		k.SetReward(ctx, elem)
	}
	for _, elem := range genState.RewardByCategoryList {
		k.SetRewardByReceiver(ctx, elem)
	}
	for _, elem := range genState.RewardByCampaignList {
		k.SetRewardByCampaign(ctx, elem)
	}
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.CampaignList = k.GetAllCampaign(ctx)
	genesis.RewardList = k.GetAllRewards(ctx)
	genesis.RewardByCategoryList = k.GetAllRewardsByReceiverAndCategory(ctx)
	genesis.RewardByCampaignList = k.GetAllRewardsByCampaign(ctx)

	return genesis
}
