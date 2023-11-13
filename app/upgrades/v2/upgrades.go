package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	sdkmath "cosmossdk.io/math"
	"github.com/sge-network/sge/app/keepers"
	housetypes "github.com/sge-network/sge/x/house/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		housePS := keepers.GetSubspace(housetypes.ModuleName)

		if !housePS.Has(ctx, []byte("MaxWithdrawalCount")) {
			var minDeposit sdkmath.Int
			housePS.Get(ctx, []byte("MinDeposit"), &minDeposit)

			var houseParticipationFee sdk.Dec
			housePS.Get(ctx, []byte("HouseParticipationFee"), &houseParticipationFee)

			p := housetypes.NewParams(
				minDeposit,
				houseParticipationFee,
				housetypes.DefaultMaxWithdrawalCount,
			)

			keepers.HouseKeeper.SetParams(ctx, p)
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
