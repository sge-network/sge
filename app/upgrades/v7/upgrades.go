package v7

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/sge-network/sge/app/keepers"
	"github.com/sge-network/sge/app/params"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		allSubaccounts := k.SubaccountKeeper.GetAllSubaccounts(ctx)

		for _, sa := range allSubaccounts {
			subAccAddr := sdk.MustAccAddressFromBech32(sa.Address)
			accSummary, found := k.SubaccountKeeper.GetAccountSummary(ctx, subAccAddr)
			if !found {
				panic(fmt.Errorf("account summary for the subaccount not found %s", subAccAddr))
			}

			_, totalBalances := k.SubaccountKeeper.GetBalances(ctx, subAccAddr, subaccounttypes.LockedBalanceStatus_LOCKED_BALANCE_STATUS_UNSPECIFIED)
			bankBalance := k.BankKeeper.GetBalance(ctx, subAccAddr, params.DefaultBondDenom).Amount

			totalBalanceDiff := accSummary.DepositedAmount.
				Sub(accSummary.SpentAmount).
				Sub(accSummary.LostAmount).
				Sub(accSummary.WithdrawnAmount)
			missingBalance := sdkmath.MinInt(bankBalance, totalBalanceDiff).Sub(totalBalances)
			if missingBalance.GT(sdkmath.ZeroInt()) {
				k.SubaccountKeeper.SetLockedBalances(ctx, subAccAddr, []subaccounttypes.LockedBalance{
					{
						Amount:   missingBalance,
						UnlockTS: 1710830000,
					},
				})
			}
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
