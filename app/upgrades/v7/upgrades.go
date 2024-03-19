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
		allSubAccounts := k.SubaccountKeeper.GetAllSubaccounts(ctx)

		for _, sa := range allSubAccounts {
			subAccAddr := sdk.MustAccAddressFromBech32(sa.Address)
			accSumm, found := k.SubaccountKeeper.GetAccountSummary(ctx, subAccAddr)
			if !found {
				panic(fmt.Errorf("account summary for the subaccount not found %s", subAccAddr))
			}

			_, totalBalances := k.SubaccountKeeper.GetBalances(ctx, subAccAddr, subaccounttypes.BalanceType_BALANCE_TYPE_UNSPECIFIED)
			bankBalance := k.BankKeeper.GetBalance(ctx, subAccAddr, params.DefaultBondDenom).Amount

			totalBalanceDiff := accSumm.DepositedAmount.Sub(totalBalances)
			missingBalance := sdkmath.MinInt(bankBalance, totalBalanceDiff)
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
