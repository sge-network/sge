package cmd

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	dbm "github.com/cometbft/cometbft-db"
	tmcfg "github.com/cometbft/cometbft/config"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cometbft/cometbft/libs/log"
	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	snapshottypes "github.com/cosmos/cosmos-sdk/snapshots/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	"github.com/sge-network/sge/app"
	"github.com/sge-network/sge/app/params"
)

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	// Set config for prefixes
	app.SetConfig()

	encodingConfig := app.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithHomeDir(app.DefaultNodeHome).
		WithViper("")

	rootCmd := &cobra.Command{
		Use:   "sged",
		Short: "SGE-Network App",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd, "", nil, tmcfg.DefaultConfig())
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	cfg := sdk.GetConfig()
	cfg.Seal()

	gentxModule, ok := app.ModuleBasics[genutiltypes.ModuleName].(genutil.AppModuleBasic)
	if !ok {
		panic("genutil AppModuleBasic assertion failed")
	}

	rootCmd.AddCommand(
		genutilcli.InitCmd(app.ModuleBasics, app.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, app.DefaultNodeHome, gentxModule.GenTxValidator),
		genutilcli.GenTxCmd(
			app.ModuleBasics,
			encodingConfig.TxConfig,
			banktypes.GenesisBalancesIterator{},
			app.DefaultNodeHome,
		),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(app.DefaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true),
		debug.Cmd(),
		config.Cmd(),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, createSimappAndExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(app.DefaultNodeHome),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	app.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
	)

	app.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// newApp is an AppCreator
func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	chainID := cast.ToString(appOpts.Get(flags.FlagChainID))
	if chainID == "" {
		// fallback to genesis chain-id
		appGenesis, err := tmtypes.GenesisDocFromFile(filepath.Join(homeDir, "config", "genesis.json"))
		if err != nil {
			panic(err)
		}

		chainID = appGenesis.ChainID
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := dbm.NewGoLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	return app.NewSgeApp(
		logger,
		db,
		traceStore,
		true,
		skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		app.MakeEncodingConfig(), // Ideally, we would reuse the one created by NewRootCmd.
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshot(snapshotStore, snapshottypes.SnapshotOptions{
			Interval:   cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval)),
			KeepRecent: cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent)),
		}),
		baseapp.SetIAVLCacheSize(cast.ToInt(appOpts.Get(server.FlagIAVLCacheSize))),
		baseapp.SetIAVLDisableFastNode(cast.ToBool(appOpts.Get(server.FlagDisableIAVLFastNode))),
		baseapp.SetChainID(chainID),
	)
}

func createSimappAndExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	encCfg := app.MakeEncodingConfig() // Ideally, we would reuse the one created by NewRootCmd.
	encCfg.Marshaler = codec.NewProtoCodec(encCfg.InterfaceRegistry)
	var sgeApp *app.SgeApp
	if height != -1 {
		sgeApp = app.NewSgeApp(
			logger,
			db,
			traceStore,
			false,
			map[int64]bool{},
			"",
			cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
			encCfg,
			appOpts,
		)

		if err := sgeApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		sgeApp = app.NewSgeApp(logger, db, traceStore, true, map[int64]bool{}, "", cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)), encCfg, appOpts)
	}

	return sgeApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
