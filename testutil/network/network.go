package network

import (
	"fmt"
	"testing"
	"time"

	tmdb "github.com/cometbft/cometbft-db"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/sge-network/sge/app"
	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/testutil/simapp"
	ovmtypes "github.com/sge-network/sge/x/ovm/types"
)

type (
	// Network represents network for simulations
	Network = network.Network
	// Config represents configurations of simulations
	Config = network.Config
)

// New creates instance with fully configured cosmos network.
// Accepts optional config, that will be used in place of the DefaultConfig() if provided.
func New(t *testing.T, configs ...network.Config) *network.Network {
	if len(configs) > 1 {
		panic("at most one config should be provided")
	}
	var cfg network.Config
	if len(configs) == 0 {
		cfg = DefaultConfig()
	} else {
		cfg = configs[0]
	}
	// baseDir := t.TempDir()

	// nodeDirName := fmt.Sprintf("node%d", 0)
	// node0Dir := filepath.Join(baseDir, nodeDirName, "simcli")

	// buf := bufio.NewReader(os.Stdin)
	// kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, node0Dir, buf, cfg.Codec, cfg.KeyringOptions...)
	// if err != nil {
	// 	panic(err)
	// }

	// keyringAlgos, _ := kb.SupportedAlgorithms()
	// algo, err := keyring.NewSigningAlgoFromString(cfg.SigningAlgo, keyringAlgos)
	// if err != nil {
	// 	panic(err)
	// }

	// addr, secret, err := testutil.GenerateSaveCoinKey(kb, nodeDirName, "", true, algo)
	// if err != nil {
	// 	panic(err)
	// }

	// info := map[string]string{"secret": secret}
	// infoBz, err := json.Marshal(info)
	// if err != nil {
	// 	panic(err)
	// }

	// // save private key seed words
	// err = simapp.WriteKeyringFile(fmt.Sprintf("%v.json", "key_seed"), node0Dir, infoBz)
	// if err != nil {
	// 	panic(err)
	// }

	net, err := network.New(t, t.TempDir(), cfg)
	if err != nil {
		panic(err)
	}

	t.Cleanup(net.Cleanup)
	return net
}

// DefaultConfig will initialize config for the network with custom application,
// genesis and single validator. All other parameters are inherited from cosmos-sdk/testutil/network.DefaultConfig
func DefaultConfig() network.Config {
	var (
		encoding = app.MakeEncodingConfig()
		chainID  = "chain-" + tmrand.NewRand().Str(6)
	)

	// simapp.GenerateSimappUsers()

	// Initialize test app by genesis account
	// genAccs := simapp.GenerateSimappGenesisAccounts()

	// Create testapp instance
	// balances := simapp.GenerateSimappUserBalances()

	defGen := app.ModuleBasics.DefaultGenesis(encoding.Marshaler)
	{
		// modify the staking denom in the genesis
		stakingGenState := defGen[stakingtypes.ModuleName]
		var newStakingGenState stakingtypes.GenesisState

		if err := encoding.Marshaler.UnmarshalJSON(stakingGenState, &newStakingGenState); err != nil {
			panic(err)
		}

		// change to default bond denom
		newStakingGenState.Params.BondDenom = params.DefaultBondDenom

		var err error
		defGen[stakingtypes.ModuleName], err = encoding.Marshaler.MarshalJSON(&newStakingGenState)
		if err != nil {
			panic(err)
		}
	}
	// {
	// 	authGenState := defGen[authtypes.ModuleName]
	// 	var newAuthGenState authtypes.GenesisState

	// 	if err := encoding.Marshaler.UnmarshalJSON(authGenState, &newAuthGenState); err != nil {
	// 		panic(err)
	// 	}

	// 	for _, v := range genAccs {
	// 		anyVal, err := codectypes.NewAnyWithValue(v)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		newAuthGenState.Accounts = append(newAuthGenState.Accounts, anyVal)
	// 	}

	// 	var err error
	// 	defGen[authtypes.ModuleName], err = encoding.Marshaler.MarshalJSON(&newAuthGenState)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// {
	// 	bankGenState := defGen[banktypes.ModuleName]
	// 	var newBankGenState banktypes.GenesisState

	// 	if err := encoding.Marshaler.UnmarshalJSON(bankGenState, &newBankGenState); err != nil {
	// 		panic(err)
	// 	}

	// 	newBankGenState.Balances = append(newBankGenState.Balances, balances...)

	// 	// totalSupply := sdk.NewCoins()
	// 	// for _, b := range balances {
	// 	// 	totalSupply = totalSupply.Add(b.Coins...)
	// 	// }
	// 	// newBankGenState.Supply = totalSupply

	// 	var err error
	// 	defGen[banktypes.ModuleName], err = encoding.Marshaler.MarshalJSON(&newBankGenState)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// }

	{
		ovmGenesisState := &ovmtypes.GenesisState{
			KeyVault: ovmtypes.KeyVault{
				PublicKeys: simapp.GenerateOvmPublicKeys(ovmtypes.MinPubKeysCount),
			},
		}
		defGen[ovmtypes.ModuleName] = encoding.Marshaler.MustMarshalJSON(ovmGenesisState)
	}

	return network.Config{
		Codec:             encoding.Marshaler,
		TxConfig:          encoding.TxConfig,
		LegacyAmino:       encoding.Amino,
		InterfaceRegistry: encoding.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(val network.ValidatorI) servertypes.Application {
			return app.NewSgeApp(
				val.GetCtx().Logger,
				tmdb.NewMemDB(),
				nil,
				true,
				map[int64]bool{},
				val.GetCtx().Config.RootDir,
				0,
				encoding,
				simtestutil.EmptyAppOptions{},
				baseapp.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
				baseapp.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
				baseapp.SetChainID(chainID),
			)
		},
		GenesisState:    defGen,
		TimeoutCommit:   2 * time.Second,
		ChainID:         chainID,
		NumValidators:   1,
		BondDenom:       params.DefaultBondDenom,
		MinGasPrices:    fmt.Sprintf("0.000006%s", params.DefaultBondDenom),
		AccountTokens:   sdk.TokensFromConsensusPower(1000, sdk.DefaultPowerReduction),
		StakingTokens:   sdk.TokensFromConsensusPower(500, sdk.DefaultPowerReduction),
		BondedTokens:    sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction),
		PruningStrategy: pruningtypes.PruningOptionNothing,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions:  []keyring.Option{},
	}
}
