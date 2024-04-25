package v9

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/sge-network/sge/app/keepers"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		allByCat := k.RewardKeeper.GetAllRewardsOfReceiverByPromoterAndCategory(ctx)
		for _, rc := range allByCat {
			k.RewardKeeper.RemoveRewardOfReceiverByPromoterAndCategory(ctx, "", rc)
			k.RewardKeeper.SetRewardOfReceiverByPromoterAndCategory(ctx, "f0630627-9e4e-48f3-8cd5-1422b46d2175", rc)
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
