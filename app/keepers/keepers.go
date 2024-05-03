package keepers

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
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

	// ibc-go

	ibc_hooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	wasmlckeeper "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/keeper"
	wasmlctypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	// cosmwasm

	wasmapp "github.com/CosmWasm/wasmd/app"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvm "github.com/CosmWasm/wasmvm"

	// sge

	mintkeeper "github.com/sge-network/sge/x/mint/keeper"
	minttypes "github.com/sge-network/sge/x/mint/types"

	betmodule "github.com/sge-network/sge/x/bet"
	betmodulekeeper "github.com/sge-network/sge/x/bet/keeper"
	betmoduletypes "github.com/sge-network/sge/x/bet/types"

	marketmodule "github.com/sge-network/sge/x/market"
	marketmodulekeeper "github.com/sge-network/sge/x/market/keeper"
	marketmoduletypes "github.com/sge-network/sge/x/market/types"

	ovmmodule "github.com/sge-network/sge/x/ovm"
	ovmmodulekeeper "github.com/sge-network/sge/x/ovm/keeper"
	ovmmoduletypes "github.com/sge-network/sge/x/ovm/types"

	housemodule "github.com/sge-network/sge/x/house"
	housemodulekeeper "github.com/sge-network/sge/x/house/keeper"
	housemoduletypes "github.com/sge-network/sge/x/house/types"

	orderbookmodule "github.com/sge-network/sge/x/orderbook"
	orderbookmodulekeeper "github.com/sge-network/sge/x/orderbook/keeper"
	orderbookmoduletypes "github.com/sge-network/sge/x/orderbook/types"

	subaccountmodule "github.com/sge-network/sge/x/subaccount"
	subaccountmodulekeeper "github.com/sge-network/sge/x/subaccount/keeper"
	subaccountmoduletypes "github.com/sge-network/sge/x/subaccount/types"

	rewardmodule "github.com/sge-network/sge/x/reward"
	rewardmodulekeeper "github.com/sge-network/sge/x/reward/keeper"
	rewardmoduletypes "github.com/sge-network/sge/x/reward/types"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
)

var wasmCapabilities = strings.Join(wasmapp.AllCapabilities(), ",")

type AppKeepers struct {
	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// SDK keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	//// CosmWasm keepers \\\\
	ContractKeeper   wasmtypes.ContractOpsKeeper
	WasmClientKeeper wasmlckeeper.Keeper
	WasmKeeper       wasmkeeper.Keeper
	ScopedWasmKeeper capabilitykeeper.ScopedKeeper

	//// SGE keepers \\\\
	BetKeeper        *betmodulekeeper.Keeper
	MarketKeeper     *marketmodulekeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	HouseKeeper      *housemodulekeeper.Keeper
	OrderbookKeeper  *orderbookmodulekeeper.Keeper
	OVMKeeper        *ovmmodulekeeper.Keeper
	RewardKeeper     *rewardmodulekeeper.Keeper
	SubaccountKeeper *subaccountmodulekeeper.Keeper

	//// SGE modules \\\\
	BetModule        betmodule.AppModule
	MarketModule     marketmodule.AppModule
	HouseModule      housemodule.AppModule
	OrderbookModule  orderbookmodule.AppModule
	OVMModule        ovmmodule.AppModule
	RewardModule     rewardmodule.AppModule
	SubaccountModule subaccountmodule.AppModule

	// make scoped keepers public for test purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper      capabilitykeeper.ScopedKeeper
	ScopedIBCFeeKeeper        capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper

	// IBC Keepers
	IBCKeeper           *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper        ibcfeekeeper.Keeper
	IBCHooksKeeper      *ibchookskeeper.Keeper
	ICAControllerKeeper icacontrollerkeeper.Keeper
	ICAHostKeeper       icahostkeeper.Keeper

	// IBC Modules
	IBCModule      ibc.AppModule
	ICAModule      ica.AppModule
	TransferModule transfer.AppModule
	IBCFeeModule   ibcfee.AppModule

	Ics20WasmHooks   *ibc_hooks.WasmHooks
	HooksICS4Wrapper ibc_hooks.ICS4Middleware
}

func NewAppKeeper(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	maccPerms map[string][]string,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	wasmOpts []wasmkeeper.Option,
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

	govModAddress := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// set the BaseApp's parameter store
	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec, appKeepers.keys[consensusparamtypes.StoreKey], govModAddress)
	bApp.SetParamStore(&appKeepers.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[capabilitytypes.StoreKey],
		appKeepers.memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.ScopedICAControllerKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	appKeepers.ScopedICAHostKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	appKeepers.ScopedWasmKeeper = appKeepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	// add keepers
	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		appKeepers.keys[crisistypes.StoreKey],
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		govModAddress,
	)

	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		AccountAddressPrefix,
		govModAddress,
	)
	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		BlockedAddresses(maccPerms),
		govModAddress,
	)
	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	groupConfig := group.DefaultConfig()
	appKeepers.GroupKeeper = groupkeeper.NewKeeper(
		appKeepers.keys[group.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper, groupConfig,
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
		govModAddress,
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
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		govModAddress,
	)
	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		cdc,
		appKeepers.keys[slashingtypes.StoreKey],
		appKeepers.StakingKeeper,
		govModAddress,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	appKeepers.StakingKeeper.SetHooks(
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
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create IBC Keeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcexported.StoreKey],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.ScopedIBCKeeper,
	)

	// register the proposal types
	govRouter := govv1beta1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper))

	govConfig := govtypes.DefaultConfig()

	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		govModAddress,
	)
	appKeepers.GovKeeper.SetLegacyRouter(govRouter)

	// Configure the hooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(
		appKeepers.keys[ibchookstypes.StoreKey],
	)
	appKeepers.IBCHooksKeeper = &hooksKeeper

	junoPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	wasmHooks := ibc_hooks.NewWasmHooks(appKeepers.IBCHooksKeeper, &appKeepers.WasmKeeper, junoPrefix) // The contract keeper needs to be set later
	appKeepers.Ics20WasmHooks = &wasmHooks
	appKeepers.HooksICS4Wrapper = ibc_hooks.NewICS4Middleware(
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.Ics20WasmHooks,
	)

	// IBC Fee Module keeper
	appKeepers.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcfeetypes.StoreKey],
		appKeepers.IBCKeeper.ChannelKeeper, // more middlewares can be added in future
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)

	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCFeeKeeper, // ISC4 Wrapper: fee IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
	)

	appKeepers.ICAControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icacontrollertypes.StoreKey],
		appKeepers.GetSubspace(icacontrollertypes.SubModuleName),
		appKeepers.IBCFeeKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedICAControllerKeeper,
		bApp.MsgServiceRouter(),
	)

	appKeepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCFeeKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.ScopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	appKeepers.EvidenceKeeper = *evidencekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evidencetypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)

	dataDir := filepath.Join(homePath, "data")

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}

	// Stargate Queries
	acceptedStargateQueries := wasmkeeper.AcceptedStargateQueries{
		// ibc
		"/ibc.core.client.v1.Query/ClientState":    &ibcclienttypes.QueryClientStateResponse{},
		"/ibc.core.client.v1.Query/ConsensusState": &ibcclienttypes.QueryConsensusStateResponse{},
		"/ibc.core.connection.v1.Query/Connection": &ibcconnectiontypes.QueryConnectionResponse{},

		// governance
		"/cosmos.gov.v1beta1.Query/Vote": &govv1.QueryVoteResponse{},

		// distribution
		"/cosmos.distribution.v1beta1.Query/DelegationRewards": &distrtypes.QueryDelegationRewardsResponse{},

		// staking
		"/cosmos.staking.v1beta1.Query/Delegation":          &stakingtypes.QueryDelegationResponse{},
		"/cosmos.staking.v1beta1.Query/Redelegations":       &stakingtypes.QueryRedelegationsResponse{},
		"/cosmos.staking.v1beta1.Query/UnbondingDelegation": &stakingtypes.QueryUnbondingDelegationResponse{},
		"/cosmos.staking.v1beta1.Query/Validator":           &stakingtypes.QueryValidatorResponse{},
		"/cosmos.staking.v1beta1.Query/Params":              &stakingtypes.QueryParamsResponse{},
		"/cosmos.staking.v1beta1.Query/Pool":                &stakingtypes.QueryPoolResponse{},
	}

	querierOpts := wasmkeeper.WithQueryPlugins(
		&wasmkeeper.QueryPlugins{
			Stargate: wasmkeeper.AcceptListStargateQuerier(acceptedStargateQueries, bApp.GRPCQueryRouter(), appCodec),
		})
	wasmOpts = append(wasmOpts, querierOpts)

	mainWasmer, err := wasmvm.NewVM(path.Join(dataDir, "wasm"), wasmCapabilities, 32, wasmConfig.ContractDebugMode, wasmConfig.MemoryCacheSize)
	if err != nil {
		panic(fmt.Sprintf("failed to create juno wasm vm: %s", err))
	}

	lcWasmer, err := wasmvm.NewVM(filepath.Join(dataDir, "light-client-wasm"), wasmCapabilities, 32, wasmConfig.ContractDebugMode, wasmConfig.MemoryCacheSize)
	if err != nil {
		panic(fmt.Sprintf("failed to create juno wasm vm for 08-wasm: %s", err))
	}

	appKeepers.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[wasmtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distrkeeper.NewQuerier(appKeepers.DistrKeeper),
		appKeepers.IBCFeeKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		dataDir,
		wasmConfig,
		wasmCapabilities,
		govModAddress,
		append(wasmOpts, wasmkeeper.WithWasmEngine(mainWasmer))...,
	)

	// 08-wasm light client
	accepted := make([]string, 0)
	for k := range acceptedStargateQueries {
		accepted = append(accepted, k)
	}

	wasmLightClientQuerier := wasmlctypes.QueryPlugins{
		// Custom: MyCustomQueryPlugin(),
		// `myAcceptList` is a `[]string` containing the list of gRPC query paths that the chain wants to allow for the `08-wasm` module to query.
		// These queries must be registered in the chain's gRPC query router, be deterministic, and track their gas usage.
		// The `AcceptListStargateQuerier` function will return a query plugin that will only allow queries for the paths in the `myAcceptList`.
		// The query responses are encoded in protobuf unlike the implementation in `x/wasm`.
		Stargate: wasmlctypes.AcceptListStargateQuerier(accepted),
	}

	appKeepers.WasmClientKeeper = wasmlckeeper.NewKeeperWithVM(
		appCodec,
		appKeepers.keys[wasmlctypes.StoreKey],
		appKeepers.IBCKeeper.ClientKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		lcWasmer,
		bApp.GRPCQueryRouter(),
		wasmlckeeper.WithQueryPlugins(&wasmLightClientQuerier),
	)

	// set the contract keeper for the Ics20WasmHooks
	appKeepers.ContractKeeper = wasmkeeper.NewDefaultPermissionKeeper(&appKeepers.WasmKeeper)
	appKeepers.Ics20WasmHooks.ContractKeeper = &appKeepers.WasmKeeper

	// // SGE keepers \\\\

	appKeepers.OrderbookKeeper = orderbookmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[orderbookmoduletypes.StoreKey],
		appKeepers.GetSubspace(orderbookmoduletypes.ModuleName),
		orderbookmodulekeeper.SdkExpectedKeepers{
			BankKeeper:     appKeepers.BankKeeper,
			AccountKeeper:  appKeepers.AccountKeeper,
			FeeGrantKeeper: appKeepers.FeeGrantKeeper,
		},
	)

	appKeepers.OVMKeeper = ovmmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ovmmoduletypes.StoreKey],
		appKeepers.keys[ovmmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(ovmmoduletypes.ModuleName),
	)

	appKeepers.MarketKeeper = marketmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[marketmoduletypes.StoreKey],
		appKeepers.keys[marketmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(marketmoduletypes.ModuleName),
	)
	appKeepers.MarketKeeper.SetOVMKeeper(appKeepers.OVMKeeper)
	appKeepers.MarketKeeper.SetOrderbookKeeper(appKeepers.OrderbookKeeper)

	appKeepers.BetKeeper = betmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[betmoduletypes.StoreKey],
		appKeepers.keys[betmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(betmoduletypes.ModuleName),
	)
	appKeepers.BetKeeper.SetMarketKeeper(appKeepers.MarketKeeper)
	appKeepers.BetKeeper.SetOrderbookKeeper(appKeepers.OrderbookKeeper)
	appKeepers.BetKeeper.SetOVMKeeper(appKeepers.OVMKeeper)

	appKeepers.OrderbookKeeper.SetBetKeeper(appKeepers.BetKeeper)
	appKeepers.OrderbookKeeper.SetMarketKeeper(appKeepers.MarketKeeper)
	appKeepers.OrderbookKeeper.SetOVMKeeper(appKeepers.OVMKeeper)

	appKeepers.HouseKeeper = housemodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[housemoduletypes.StoreKey],
		appKeepers.OrderbookKeeper,
		appKeepers.OVMKeeper,
		appKeepers.GetSubspace(housemoduletypes.ModuleName),
		housemodulekeeper.SdkExpectedKeepers{
			AuthzKeeper: appKeepers.AuthzKeeper,
		},
	)
	appKeepers.OrderbookKeeper.SetHouseKeeper(appKeepers.HouseKeeper)

	appKeepers.SubaccountKeeper = subaccountmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[subaccountmoduletypes.StoreKey],
		appKeepers.GetSubspace(subaccountmoduletypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.OVMKeeper,
		appKeepers.BetKeeper,
		appKeepers.OrderbookKeeper,
		appKeepers.HouseKeeper,
	)

	appKeepers.RewardKeeper = rewardmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[rewardmoduletypes.StoreKey],
		appKeepers.keys[rewardmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(rewardmoduletypes.ModuleName),
		appKeepers.BetKeeper,
		appKeepers.OVMKeeper,
		appKeepers.SubaccountKeeper,
		rewardmodulekeeper.SdkExpectedKeepers{
			AuthzKeeper:   appKeepers.AuthzKeeper,
			BankKeeper:    appKeepers.BankKeeper,
			AccountKeeper: appKeepers.AccountKeeper,
		},
	)

	// ** Hooks ** \\

	appKeepers.OrderbookKeeper.SetHooks(
		orderbookmoduletypes.NewMultiOrderBookHooks(
			appKeepers.SubaccountKeeper.Hooks(),
		),
	)

	// // SGE modules \\\\

	appKeepers.BetModule = betmodule.NewAppModule(
		appCodec,
		*appKeepers.BetKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.MarketKeeper,
		appKeepers.OrderbookKeeper,
		appKeepers.OVMKeeper,
	)
	appKeepers.MarketModule = marketmodule.NewAppModule(
		appCodec,
		*appKeepers.MarketKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.OVMKeeper,
	)
	appKeepers.HouseModule = housemodule.NewAppModule(
		appCodec,
		*appKeepers.HouseKeeper,
	)
	appKeepers.OrderbookModule = orderbookmodule.NewAppModule(
		appCodec,
		*appKeepers.OrderbookKeeper,
	)
	appKeepers.OVMModule = ovmmodule.NewAppModule(
		appCodec,
		*appKeepers.OVMKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)
	appKeepers.RewardModule = rewardmodule.NewAppModule(
		appCodec,
		*appKeepers.RewardKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)

	appKeepers.SubaccountModule = subaccountmodule.NewAppModule(*appKeepers.SubaccountKeeper)

	//// IBC modules \\\\
	appKeepers.IBCModule = ibc.NewAppModule(appKeepers.IBCKeeper)
	appKeepers.IBCFeeModule = ibcfee.NewAppModule(appKeepers.IBCFeeKeeper)
	appKeepers.TransferModule = transfer.NewAppModule(appKeepers.TransferKeeper)
	appKeepers.ICAModule = ica.NewAppModule(&appKeepers.ICAControllerKeeper, &appKeepers.ICAHostKeeper)

	// IBC stacks \\\
	var transferStack ibcporttypes.IBCModule
	transferStack = transfer.NewIBCModule(appKeepers.TransferKeeper)
	transferStack = ibcfee.NewIBCMiddleware(transferStack, appKeepers.IBCFeeKeeper)

	var icaControllerStack ibcporttypes.IBCModule
	icaControllerStack = icacontroller.NewIBCMiddleware(icaControllerStack, appKeepers.ICAControllerKeeper)
	icaControllerStack = ibcfee.NewIBCMiddleware(icaControllerStack, appKeepers.IBCFeeKeeper)

	var icaHostStack ibcporttypes.IBCModule
	icaHostStack = icahost.NewIBCModule(appKeepers.ICAHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, appKeepers.IBCFeeKeeper)

	// Create fee enabled wasm ibc Stack
	var wasmStack porttypes.IBCModule
	wasmStack = wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, appKeepers.IBCFeeKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(icacontrollertypes.SubModuleName, icaControllerStack)
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostStack)
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)
	ibcRouter.AddRoute(wasmtypes.ModuleName, wasmStack)

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
	key, tkey storetypes.StoreKey,
) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(betmoduletypes.ModuleName)
	paramsKeeper.Subspace(marketmoduletypes.ModuleName)
	paramsKeeper.Subspace(orderbookmoduletypes.ModuleName)
	paramsKeeper.Subspace(ovmmoduletypes.ModuleName)
	paramsKeeper.Subspace(housemoduletypes.ModuleName)
	paramsKeeper.Subspace(rewardmoduletypes.ModuleName)
	paramsKeeper.Subspace(subaccountmoduletypes.ModuleName)
	paramsKeeper.Subspace(wasmtypes.ModuleName)

	return paramsKeeper
}

// BlockedAddresses returns all the app's blocked account addresses.
func BlockedAddresses(maccPerms map[string][]string) map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range GetMaccPerms(maccPerms) {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	// allow the following addresses to receive funds
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return modAccAddrs
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms(maccPerms map[string][]string) map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}

	return dupMaccPerms
}
