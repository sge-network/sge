package v2

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/sge-network/sge/app/keepers"
)

// betOddsTypeRemoval removes bet odds type from current state
func betOddsTypeRemoval(keepers *keepers.AppKeepers, ctx sdk.Context) {
	oldBets, err := keepers.BetKeeper.GetBets(ctx)
	if err != nil {
		panic(err)
	}

	for _, old := range oldBets {
		betID, found := keepers.BetKeeper.GetBetID(ctx, old.UID)
		if !found {
			panic(fmt.Errorf("bet id for the uid not found %s", old.UID))
		}

		keepers.BetKeeper.SetBet(ctx, old, betID.ID)
	}
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		betOddsTypeRemoval(keepers, ctx)

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
