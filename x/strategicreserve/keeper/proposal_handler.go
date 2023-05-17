package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bettypes "github.com/sge-network/sge/x/bet/types"
	housetypes "github.com/sge-network/sge/x/house/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// HandleDataFeeCollectorFeedProposal is a handler for executing a passed data fee collector feed proposal
func HandleDataFeeCollectorFeedProposal(ctx sdk.Context, k Keeper, p *types.DataFeeCollectorFeedProposal) error {
	if p.HouseFeeSpend.GT(sdk.ZeroInt()) {
		if err := k.transferFundsFromModuleToModule(ctx, housetypes.HouseFeeCollector, types.DataFeeCollector, p.HouseFeeSpend); err != nil {
			return err
		}
	}

	if p.BetFeeSpend.GT(sdk.ZeroInt()) {
		if err := k.transferFundsFromModuleToModule(ctx, bettypes.BetFeeCollector, types.DataFeeCollector, p.BetFeeSpend); err != nil {
			return err
		}
	}

	logger := k.Logger(ctx)
	logger.Info("transferred from the House fee collector to Data fee collector", "amount", p.HouseFeeSpend.String())
	logger.Info("transferred from the Bet fee collector to Data fee collector", "amount", p.BetFeeSpend.String())

	return nil
}
