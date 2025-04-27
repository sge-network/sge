package v10

import (
	"context"

	sdkmath "cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

	"github.com/sge-network/sge/app/keepers"
	"github.com/sge-network/sge/app/params"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		k.AccountKeeper.IterateAccounts(ctx, func(acc sdk.AccountI) bool {
			vestingAcc, ok := acc.(*vestingtypes.PeriodicVestingAccount)
			if !ok {
				return false
			}

			currentTime := sdkCtx.BlockTime().Unix()

			originalVesting := vestingAcc.OriginalVesting.AmountOf(params.DefaultBondDenom)

			var refinedPeriods vestingtypes.Periods
			endTime := vestingAcc.StartTime
			vestedBalance := sdkmath.ZeroInt()
			for _, period := range vestingAcc.VestingPeriods {
				endTime += period.Length
				if currentTime > endTime {
					// past periods
					vestedBalance = vestedBalance.Add(period.Amount.AmountOf(params.DefaultBondDenom))
					refinedPeriods = append(refinedPeriods, period)
				} else {
					break
				}
			}

			if originalVesting.GT(vestedBalance) {
				// add the remaining balance to the last period
				refinedPeriods = append(refinedPeriods, vestingtypes.Period{
					Length: 2592000 * 24, // 24 months
					Amount: sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, originalVesting.Sub(vestedBalance))),
				})

				refinedEndTime := vestingAcc.StartTime
				for _, period := range refinedPeriods {
					refinedEndTime += period.Length
				}

				vestingAcc.VestingPeriods = refinedPeriods
				vestingAcc.EndTime = refinedEndTime

				k.AccountKeeper.SetAccount(ctx, vestingAcc)
			}

			return false
		})

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
