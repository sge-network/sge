package v6

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/sge-network/sge/app/keepers"
	"github.com/sge-network/sge/x/market/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {

		participationList, err := k.OrderbookKeeper.GetAllOrderBookParticipations(ctx)
		if err != nil {
			panic(err)
		}

		for _, p := range participationList {
			if !p.IsSettled {
				continue
			}

			market, found := k.MarketKeeper.GetMarket(ctx, p.OrderBookUID)
			if !found {
				panic(fmt.Errorf("market not found %s", p.OrderBookUID))
			}

			reimburseFee := false
			if market.Status == types.MarketStatus_MARKET_STATUS_CANCELED ||
				market.Status == types.MarketStatus_MARKET_STATUS_ABORTED {
				reimburseFee = true
				p.ReimbursedLiquidity = p.Liquidity
			}

			if reimburseFee || p.NotParticipatedInBetFulfillment() {
				p.ReimbursedFee = p.Fee
			}

			k.OrderbookKeeper.SetOrderBookParticipation(ctx, p)
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
