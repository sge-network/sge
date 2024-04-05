package v8

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/google/uuid"
	"github.com/sge-network/sge/app/keepers"
	"github.com/sge-network/sge/x/reward/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		allCampaigns := k.RewardKeeper.GetAllCampaign(ctx)
		promoters := make(map[string]struct{})
		for _, c := range allCampaigns {
			c.CapCount = 0 // infinite cap for all campaigns
			c.Pool.Withdrawn = sdkmath.ZeroInt()
			k.RewardKeeper.SetCampaign(ctx, c)
			promoters[c.Promoter] = struct{}{}
		}

		promoterAddresses := []string{}
		for addr := range promoters {
			promoterAddresses = append(promoterAddresses, addr)
		}

		if len(promoterAddresses) > 0 {
			promoterUID := uuid.NewString()
			k.RewardKeeper.SetPromoter(ctx, types.Promoter{
				Creator:   promoterAddresses[0],
				UID:       promoterUID,
				Addresses: promoterAddresses,
				Conf: types.PromoterConf{
					CategoryCap: []types.CategoryCap{
						{Category: types.RewardCategory_REWARD_CATEGORY_SIGNUP, CapPerAcc: 1},
					},
				},
			})

			for _, addr := range promoterAddresses {
				k.RewardKeeper.SetPromoterByAddress(ctx, types.PromoterByAddress{
					PromoterUID: promoterUID,
					Address:     addr,
				})
			}
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
