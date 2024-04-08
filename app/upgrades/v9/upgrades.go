package v9

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"

	"github.com/sge-network/sge/app/keepers"
	betmoduletypes "github.com/sge-network/sge/x/bet/types"
	housemoduletypes "github.com/sge-network/sge/x/house/types"
	marketmoduletypes "github.com/sge-network/sge/x/market/types"
	minttypes "github.com/sge-network/sge/x/mint/types"
	orderbookmoduletypes "github.com/sge-network/sge/x/orderbook/types"
	ovmmoduletypes "github.com/sge-network/sge/x/ovm/types"
	rewardmoduletypes "github.com/sge-network/sge/x/reward/types"
	subaccountmoduletypes "github.com/sge-network/sge/x/subaccount/types"
)

//nolint:staticcheck
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	k *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// https://github.com/cosmos/cosmos-sdk/pull/12363/files
		// Set param key table for params module migration
		for _, subspace := range k.ParamsKeeper.GetSubspaces() {
			subspace := subspace

			var keyTable paramstypes.KeyTable
			switch subspace.Name() {
			case authtypes.ModuleName:
				keyTable = authtypes.ParamKeyTable()
			case banktypes.ModuleName:
				keyTable = banktypes.ParamKeyTable()
			case stakingtypes.ModuleName:
				keyTable = stakingtypes.ParamKeyTable()

			// case minttypes.ModuleName:
			// 	keyTable = minttypes.ParamKeyTable()
			case distrtypes.ModuleName:
				keyTable = distrtypes.ParamKeyTable()
			case slashingtypes.ModuleName:
				keyTable = slashingtypes.ParamKeyTable()
			case govtypes.ModuleName:
				keyTable = govv1.ParamKeyTable()
			case crisistypes.ModuleName:
				keyTable = crisistypes.ParamKeyTable()

			// ibc types
			case ibctransfertypes.ModuleName:
				keyTable = ibctransfertypes.ParamKeyTable()
			case icahosttypes.SubModuleName:
				keyTable = icahosttypes.ParamKeyTable()
			case icacontrollertypes.SubModuleName:
				keyTable = icacontrollertypes.ParamKeyTable()

			// sge modules
			case betmoduletypes.ModuleName:
				keyTable = betmoduletypes.ParamKeyTable()
			case housemoduletypes.ModuleName:
				keyTable = housemoduletypes.ParamKeyTable()
			case marketmoduletypes.ModuleName:
				keyTable = marketmoduletypes.ParamKeyTable()
			case minttypes.ModuleName:
				keyTable = minttypes.ParamKeyTable()
			case orderbookmoduletypes.ModuleName:
				keyTable = orderbookmoduletypes.ParamKeyTable()
			case ovmmoduletypes.ModuleName:
				keyTable = ovmmoduletypes.ParamKeyTable()
			case subaccountmoduletypes.ModuleName:
				keyTable = subaccountmoduletypes.ParamKeyTable()
			case rewardmoduletypes.ModuleName:
				keyTable = rewardmoduletypes.ParamKeyTable()
			}

			if !subspace.HasKeyTable() {
				subspace.WithKeyTable(keyTable)
			}
		}

		// Migrate Tendermint consensus parameters from x/params module to a
		// dedicated x/consensus module.
		baseAppLegacySS := k.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable())
		baseapp.MigrateParams(ctx, baseAppLegacySS, &k.ConsensusParamsKeeper)

		// https://github.com/cosmos/ibc-go/blob/v7.1.0/docs/migrations/v7-to-v7_1.md
		// explicitly update the IBC 02-client params, adding the localhost client type
		params := k.IBCKeeper.ClientKeeper.GetParams(ctx)
		params.AllowedClients = append(params.AllowedClients, exported.Localhost)
		k.IBCKeeper.ClientKeeper.SetParams(ctx, params)

		// update gov params to use a 20% initial deposit ratio, allowing us to remote the ante handler
		govParams := k.GovKeeper.GetParams(ctx)
		govParams.MinInitialDepositRatio = sdkmath.LegacyNewDec(20).Quo(sdkmath.LegacyNewDec(100)).String()
		if err := k.GovKeeper.SetParams(ctx, govParams); err != nil {
			return nil, err
		}

		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
