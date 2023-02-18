package keepers

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	mintkeeper "github.com/sge-network/sge/x/mint/keeper"
	minttypes "github.com/sge-network/sge/x/mint/types"

	betmodule "github.com/sge-network/sge/x/bet"
	betmodulekeeper "github.com/sge-network/sge/x/bet/keeper"
	betmoduletypes "github.com/sge-network/sge/x/bet/types"

	sporteventmodule "github.com/sge-network/sge/x/sportevent"
	sporteventmodulekeeper "github.com/sge-network/sge/x/sportevent/keeper"
	sporteventmoduletypes "github.com/sge-network/sge/x/sportevent/types"

	strategicreservemodule "github.com/sge-network/sge/x/strategicreserve"
	strategicreservemodulekeeper "github.com/sge-network/sge/x/strategicreserve/keeper"
	strategicreservemoduletypes "github.com/sge-network/sge/x/strategicreserve/types"

	dvmmodule "github.com/sge-network/sge/x/dvm"
	dvmmodulekeeper "github.com/sge-network/sge/x/dvm/keeper"
	dvmmoduletypes "github.com/sge-network/sge/x/dvm/types"

	housemodule "github.com/sge-network/sge/x/house"
	housemodulekeeper "github.com/sge-network/sge/x/house/keeper"
	housemoduletypes "github.com/sge-network/sge/x/house/types"

	orderbookmodule "github.com/sge-network/sge/x/orderbook"
	orderbookmodulekeeper "github.com/sge-network/sge/x/orderbook/keeper"
	orderbookmoduletypes "github.com/sge-network/sge/x/orderbook/types"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
)

type AppKeepers struct {
	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper       authkeeper.AccountKeeper
	BankKeeper          bankkeeper.Keeper
	CapabilityKeeper    *capabilitykeeper.Keeper
	StakingKeeper       stakingkeeper.Keeper
	SlashingKeeper      slashingkeeper.Keeper
	MintKeeper          mintkeeper.Keeper
	DistrKeeper         distrkeeper.Keeper
	GovKeeper           govkeeper.Keeper
	CrisisKeeper        crisiskeeper.Keeper
	UpgradeKeeper       upgradekeeper.Keeper
	ParamsKeeper        paramskeeper.Keeper
	IBCKeeper           *ibckeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper
	EvidenceKeeper      evidencekeeper.Keeper
	TransferKeeper      ibctransferkeeper.Keeper
	FeeGrantKeeper      feegrantkeeper.Keeper
	AuthzKeeper         authzkeeper.Keeper

	StrategicreserveKeeper strategicreservemodulekeeper.Keeper
	SporteventKeeper       sporteventmodulekeeper.Keeper
	BetKeeper              betmodulekeeper.Keeper
	DVMKeeper              dvmmodulekeeper.Keeper
	OrderBookKeeper        orderbookmodulekeeper.Keeper
	HouseKeeper            housemodulekeeper.Keeper
	SporteventModule       sporteventmodule.AppModule
	StrategicreserveModule strategicreservemodule.AppModule
	BetModule              betmodule.AppModule
	OrderBookModule        orderbookmodule.AppModule
	HouseModule            housemodule.AppModule

	// modules
	ICAModule      ica.AppModule
	TransferModule transfer.AppModule
	DVMModule      dvmmodule.AppModule

	// make scoped keepers public for test purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper      capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
}

func NewAppKeeper(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	maccPerms map[string][]string,
	moduleAccAddress map[string]bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	appOpts servertypes.AppOptions,
) AppKeepers {
	appKeepers := AppKeepers{}
	// Set keys KVStoreKey, TransientStoreKey, MemoryStoreKey
	appKeepers.GenerateKeys()

	appKeepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		cdc,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)

	// set the BaseApp's parameter store
	bApp.SetParamStore(
		appKeepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()),
	)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.ScopedICAControllerKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	appKeepers.ScopedICAHostKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)

	appKeepers.CapabilityKeeper.Seal()

	// add keepers
	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appKeepers.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
	)

	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		appKeepers.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
	)
	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.GetSubspace(banktypes.ModuleName),
		moduleAccAddress,
	)
	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[feegrant.StoreKey],
		appKeepers.AccountKeeper,
	)

	appKeepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.GetSubspace(stakingtypes.ModuleName),
	)

	appKeepers.StrategicreserveKeeper = *strategicreservemodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[strategicreservemoduletypes.StoreKey],
		appKeepers.keys[strategicreservemoduletypes.MemStoreKey],
		appKeepers.GetSubspace(strategicreservemoduletypes.ModuleName),
		strategicreservemodulekeeper.ExpectedKeepers{
			BankKeeper:    appKeepers.BankKeeper,
			AccountKeeper: appKeepers.AccountKeeper,
		},
	)

	appKeepers.MintKeeper = *mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.GetSubspace(minttypes.ModuleName),
		appKeepers.AccountKeeper,
		mintkeeper.ExpectedKeepers{
			StakingKeeper: appKeepers.StakingKeeper,
			BankKeeper:    appKeepers.BankKeeper,
		},
		authtypes.FeeCollectorName,
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[distrtypes.StoreKey],
		appKeepers.GetSubspace(distrtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		moduleAccAddress,
	)
	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[slashingtypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.GetSubspace(slashingtypes.ModuleName),
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	appKeepers.StakingKeeper = *appKeepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
		),
	)

	// UpgradeKeeper must be created before IBCKeeper
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
	)

	// Create IBC Keeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibchost.StoreKey],
		appKeepers.GetSubspace(ibchost.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.ScopedIBCKeeper,
	)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(appKeepers.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper))

	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.GetSubspace(govtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		govRouter,
	)

	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
	)
	appKeepers.TransferModule = transfer.NewAppModule(appKeepers.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(appKeepers.TransferKeeper)

	appKeepers.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec, appKeepers.keys[icacontrollertypes.StoreKey], appKeepers.GetSubspace(icacontrollertypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		appKeepers.IBCKeeper.ChannelKeeper, &appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedICAControllerKeeper, bApp.MsgServiceRouter(),
	)

	appKeepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.ScopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)

	appKeepers.ICAModule = ica.NewAppModule(&appKeepers.ICAControllerKeeper, &appKeepers.ICAHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(appKeepers.ICAHostKeeper)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	appKeepers.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evidencetypes.StoreKey],
		&appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)

	appKeepers.DVMKeeper = *dvmmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[dvmmoduletypes.StoreKey],
		appKeepers.keys[dvmmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(dvmmoduletypes.ModuleName),
	)
	appKeepers.DVMModule = dvmmodule.NewAppModule(appCodec, appKeepers.DVMKeeper, appKeepers.AccountKeeper, appKeepers.BankKeeper)

	appKeepers.OrderBookKeeper = *orderbookmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[orderbookmoduletypes.StoreKey],
		appKeepers.GetSubspace(orderbookmoduletypes.ModuleName),
		appKeepers.BankKeeper,
		appKeepers.AccountKeeper,
	)
	appKeepers.OrderBookModule = orderbookmodule.NewAppModule(appCodec, appKeepers.OrderBookKeeper)

	appKeepers.SporteventKeeper = *sporteventmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[sporteventmoduletypes.StoreKey],
		appKeepers.keys[sporteventmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(sporteventmoduletypes.ModuleName),
		sporteventmodulekeeper.ExpectedKeepers{
			DVMKeeper:  appKeepers.DVMKeeper,
			BookKeeper: appKeepers.OrderBookKeeper,
		},
	)
	appKeepers.SporteventModule = sporteventmodule.NewAppModule(appCodec, appKeepers.SporteventKeeper, appKeepers.AccountKeeper, appKeepers.BankKeeper, appKeepers.DVMKeeper)

	appKeepers.StrategicreserveModule = strategicreservemodule.NewAppModule(appCodec, appKeepers.StrategicreserveKeeper, appKeepers.AccountKeeper, appKeepers.BankKeeper)

	appKeepers.BetKeeper = *betmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[betmoduletypes.StoreKey],
		appKeepers.keys[betmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(betmoduletypes.ModuleName),
		betmodulekeeper.ExpectedKeepers{
			SporteventKeeper: appKeepers.SporteventKeeper,
			OrderBookKeeper:  appKeepers.OrderBookKeeper,
			DVMKeeper:        appKeepers.DVMKeeper,
		},
	)
	appKeepers.BetModule = betmodule.NewAppModule(appCodec, appKeepers.BetKeeper, appKeepers.AccountKeeper, appKeepers.BankKeeper, appKeepers.SporteventKeeper, appKeepers.OrderBookKeeper, appKeepers.DVMKeeper)

	appKeepers.HouseKeeper = *housemodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[housemoduletypes.StoreKey],
		appKeepers.OrderBookKeeper,
		appKeepers.GetSubspace(housemoduletypes.ModuleName),
	)
	appKeepers.HouseModule = housemodule.NewAppModule(appCodec, appKeepers.HouseKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibctransfertypes.ModuleName, transferIBCModule)

	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	/****  Module Options ****/
	return appKeepers
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec,
	legacyAmino *codec.LegacyAmino,
	key, tkey sdk.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(betmoduletypes.ModuleName)
	paramsKeeper.Subspace(sporteventmoduletypes.ModuleName)
	paramsKeeper.Subspace(strategicreservemoduletypes.ModuleName)
	paramsKeeper.Subspace(dvmmoduletypes.ModuleName)
	paramsKeeper.Subspace(orderbookmoduletypes.ModuleName)
	paramsKeeper.Subspace(housemoduletypes.ModuleName)

	return paramsKeeper
}
