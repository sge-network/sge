package v7

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/sge-network/sge/app/keepers"
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

			allBalances, totalAmount := k.SubaccountKeeper.GetBalances(ctx, subAccAddr, subaccounttypes.BalanceType_BALANCE_TYPE_UNSPECIFIED)

			if totalAmount.LT(accSumm.DepositedAmount) {
				for _, b := range allBalances {
					if b.Amount.IsZero() {
						b.Amount = accSumm.DepositedAmount.Sub(totalAmount)
						break
					}
				}
			}
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
