package v10

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/sge-network/sge/app/keepers"
	"github.com/sge-network/sge/app/params"
	housetypes "github.com/sge-network/sge/x/house/types"
	orderbooktypes "github.com/sge-network/sge/x/orderbook/types"
	rewardtypes "github.com/sge-network/sge/x/reward/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		var recAddress sdk.AccAddress
		switch strings.ToLower(ctx.ChainID()) {
		case "sgenet-1":
			recAddress = sdk.MustAccAddressFromBech32("sge10gtzxh97gvgt6s58w5gagsacrwddhmgzn5ksk6")
		case "sge-network-4":
			recAddress = sdk.MustAccAddressFromBech32("sge19dkag3dkzcmjzhctz8ysfws443yq0nyyu0r8c0")
		case "stage-sgenetwork":
			recAddress = sdk.MustAccAddressFromBech32("sge1wkzh54s97m4cr70wfq9j7e4u42mkt2a7hak84c")
		default:
			return mm.RunMigrations(ctx, configurator, fromVM)
		}

		orderBookBalance := k.BankKeeper.GetBalance(
			ctx,
			k.AccountKeeper.GetModuleAddress(orderbooktypes.OrderBookLiquidityFunder{}.GetModuleAcc()),
			params.DefaultBondDenom,
		)
		if orderBookBalance.Amount.GT(sdkmath.ZeroInt()) {
			if err := k.BankKeeper.SendCoinsFromModuleToAccount(
				ctx,
				orderbooktypes.OrderBookLiquidityFunder{}.GetModuleAcc(),
				recAddress, sdk.NewCoins(orderBookBalance),
			); err != nil {
				return nil, err
			}
		}

		houseFeeBalance := k.BankKeeper.GetBalance(
			ctx,
			k.AccountKeeper.GetModuleAddress(housetypes.HouseFeeCollectorFunder{}.GetModuleAcc()),
			params.DefaultBondDenom,
		)
		if houseFeeBalance.Amount.GT(sdkmath.ZeroInt()) {
			if err := k.BankKeeper.SendCoinsFromModuleToAccount(
				ctx,
				housetypes.HouseFeeCollectorFunder{}.GetModuleAcc(),
				recAddress, sdk.NewCoins(houseFeeBalance),
			); err != nil {
				return nil, err
			}
		}

		rewardPoolBalance := k.BankKeeper.GetBalance(
			ctx,
			k.AccountKeeper.GetModuleAddress(rewardtypes.RewardPoolFunder{}.GetModuleAcc()),
			params.DefaultBondDenom,
		)
		if rewardPoolBalance.Amount.GT(sdkmath.ZeroInt()) {
			if err := k.BankKeeper.SendCoinsFromModuleToAccount(
				ctx,
				rewardtypes.RewardPoolFunder{}.GetModuleAcc(),
				recAddress, sdk.NewCoins(rewardPoolBalance),
			); err != nil {
				return nil, err
			}
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
