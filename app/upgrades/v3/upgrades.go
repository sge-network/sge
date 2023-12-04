package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/sge-network/sge/app/keepers"
	"time"
)

const DefaultExpeditedPeriod time.Duration = time.Hour * 24

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		depositParams := keepers.GovKeeper.GetDepositParams(ctx)
		depositParams.MinExpeditedDeposit = sdk.NewCoins(sdk.NewCoin("usge", DefaultMinExpeditedDepositTokens))
		keepers.GovKeeper.SetDepositParams(ctx, depositParams)

		tallyParams := keepers.GovKeeper.GetTallyParams(ctx)
		tallyParams.ExpeditedThreshold = DefaultExpeditedThreshold.String()
		tallyParams.ExpeditedQuorum = DefaultExpeditedQuorum.String()
		keepers.GovKeeper.SetTallyParams(ctx, tallyParams)

		votingParams := keepers.GovKeeper.GetVotingParams(ctx)
		expeditedPeriod := DefaultExpeditedPeriod
		votingParams.ExpeditedVotingPeriod = &expeditedPeriod
		keepers.GovKeeper.SetVotingParams(ctx, votingParams)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
