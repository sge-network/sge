package v110

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/sge-network/sge/app/keepers"
)

func setCommissionRate(stakingKeeper *stakingkeeper.Keeper, ctx sdk.Context) {
	// the minimal commission rate of 5% (0.05)
	// (default is needed to be set because of SDK store migrations that set the param)
	stakingtypes.DefaultMinCommissionRate = sdk.NewDecWithPrec(5, 2)

	stakingKeeper.IterateValidators(ctx, func(index int64, val stakingtypes.ValidatorI) (stop bool) {
		if val.GetCommission().LT(stakingtypes.DefaultMinCommissionRate) {
			validator, found := stakingKeeper.GetValidator(ctx, val.GetOperator())
			if !found {
				ctx.Logger().Error("validator not found", val)
				return true
			}
			ctx.Logger().Info("update validator's commission rate to a minimal one", val)
			validator.Commission.Rate = stakingtypes.DefaultMinCommissionRate
			if validator.Commission.MaxRate.LT(stakingtypes.DefaultMinCommissionRate) {
				validator.Commission.MaxRate = stakingtypes.DefaultMinCommissionRate
			}
			stakingKeeper.SetValidator(ctx, validator)
		}
		return false
	})
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		setCommissionRate(&keepers.StakingKeeper, ctx)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
