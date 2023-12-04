package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	v1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/sge-network/sge/app/keepers"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		depositParams := keepers.GovKeeper.GetDepositParams(ctx)
		depositParams.MinExpeditedDeposit = sdk.NewCoins(sdk.NewCoin("usge", v1.DefaultMinExpeditedDepositTokens))
		keepers.GovKeeper.SetDepositParams(ctx, depositParams)

		tallyParams := keepers.GovKeeper.GetTallyParams(ctx)
		tallyParams.ExpeditedThreshold = v1.DefaultExpeditedThreshold.String()
		tallyParams.ExpeditedQuorum = v1.DefaultExpeditedQuorum.String()
		keepers.GovKeeper.SetTallyParams(ctx, tallyParams)

		votingParams := keepers.GovKeeper.GetVotingParams(ctx)
		expeditedPeriod := v1.DefaultExpeditedPeriod
		votingParams.ExpeditedVotingPeriod = &expeditedPeriod
		keepers.GovKeeper.SetVotingParams(ctx, votingParams)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
