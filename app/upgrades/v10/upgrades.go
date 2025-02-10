package v10

import (
	"context"
	"time"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/sge-network/sge/app/keepers"
	"github.com/sge-network/sge/app/params"
	minttypes "github.com/sge-network/sge/x/mint/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// https://github.com/cosmos/cosmos-sdk/pull/12363/files
		// Set param key table for params module migration
		for _, subspace := range k.ParamsKeeper.GetSubspaces() {
			var keyTable paramstypes.KeyTable
			switch subspace.Name() {
			// sdk
			case minttypes.ModuleName:
				keyTable = minttypes.ParamKeyTable()
			default:
				continue
			}

			if !subspace.HasKeyTable() {
				subspace.WithKeyTable(keyTable)
			}
		}

		migrations, err := mm.RunMigrations(ctx, configurator, fromVM)
		if err != nil {
			return nil, err
		}

		// Set expedited proposal param:
		govParams, err := k.GovKeeper.Params.Get(ctx)
		if err != nil {
			return nil, err
		}
		govParams.ExpeditedMinDeposit = sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdkmath.NewInt(50000000000)))
		govParams.MinInitialDepositRatio = "0.000000000000000000"
		govParams.ExpeditedThreshold = "0.750000000000000000"
		expediteVotingPeriod := 86400 * time.Second
		govParams.ExpeditedVotingPeriod = &expediteVotingPeriod
		err = k.GovKeeper.Params.Set(ctx, govParams)
		if err != nil {
			return nil, err
		}

		return migrations, nil
	}
}
