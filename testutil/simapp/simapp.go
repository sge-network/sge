package simapp

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"testing"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/sge-network/sge/app"
	"github.com/sge-network/sge/app/params"
	dvmtypes "github.com/sge-network/sge/x/dvm/types"
	mintmoduletypes "github.com/sge-network/sge/x/mint/types"
	strategicreservetypes "github.com/sge-network/sge/x/strategicreserve/types"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

// TestApp is used as a container of the sge app
type TestApp struct {
	app.SgeApp
}

// SimappOptions defines options related to simapp initialization
type SimappOptions struct {
	CreateGenesisValidators bool
}

// setup initializes new test app instance
func setup(withGenesis bool, invCheckPeriod uint) (*TestApp, app.GenesisState) {
	db := tmdb.NewMemDB()
	encCdc := app.MakeEncodingConfig()
	appInstance := app.NewSgeApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, "", invCheckPeriod, encCdc,
		// this line is used by starport scaffolding # stargate/testutil/appArgument
		simapp.EmptyAppOptions{})
	if withGenesis {
		return &TestApp{SgeApp: *appInstance}, app.NewDefaultGenesisState()
	}
	return &TestApp{SgeApp: *appInstance}, app.GenesisState{}
}

// Setup initializes genesis the same as simapp
func Setup(isCheckTx bool) *TestApp {
	app, genesisState := setup(!isCheckTx, 5)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

// SetupWithGenesisAccounts sets up the genesis accounts for testing
func SetupWithGenesisAccounts(genAccs []authtypes.GenesisAccount, options SimappOptions, balances ...banktypes.Balance) *TestApp {
	appInstance, genesisState := setup(true, 0)

	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = appInstance.AppCodec().MustMarshalJSON(authGenesis)

	var validatorUpdates []abci.ValidatorUpdate
	if options.CreateGenesisValidators {
		var moduleBalance banktypes.Balance
		var stakingGenesis *stakingtypes.GenesisState

		stakingGenesis, validatorUpdates, moduleBalance = stakingDefaultTestGenesis(appInstance)
		genesisState[stakingtypes.ModuleName] = appInstance.AppCodec().MustMarshalJSON(stakingGenesis)

		balances = append(balances, moduleBalance)
	}

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		totalSupply = totalSupply.Add(b.Coins...)
	}

	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = appInstance.AppCodec().MustMarshalJSON(bankGenesis)

	{
		bs, err := x509.MarshalPKIXPublicKey(TestDVMPublicKey)
		if err != nil {
			panic(err)
		}

		dvmGenesisState := &dvmtypes.GenesisState{
			PublicKeys: &dvmtypes.PublicKeys{
				List: []string{string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: bs}))},
			},
		}
		genesisState[dvmtypes.ModuleName] = appInstance.AppCodec().MustMarshalJSON(dvmGenesisState)
	}

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	appInstance.InitChain(
		abci.RequestInitChain{
			Validators:      validatorUpdates,
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	appInstance.Commit()
	appInstance.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: appInstance.LastBlockHeight() + 1}})

	return appInstance
}

// GetTestObjects gets the test objects and ingredients for testing phase start with default options
func GetTestObjects() (*TestApp, sdk.Context, error) {
	// return
	return GetTestObjectsWithOptions(SimappOptions{
		CreateGenesisValidators: true,
	})
}

// GetTestObjectsWithOptions gets the test objects and ingredients for testing phase start with custom options
func GetTestObjectsWithOptions(options SimappOptions) (*TestApp, sdk.Context, error) {
	generateSimappUsers()

	// Initialize test app by genesis account
	genAccs := generateSimappGenesisAccounts()

	// Create testapp instance
	balances := generateSimappUserBalances()

	tApp := SetupWithGenesisAccounts(genAccs, options, balances...)

	// Create the context
	ctx := tApp.NewContext(true, tmproto.Header{Height: tApp.LastBlockHeight()})

	//set minter params
	setMinterParams(tApp, ctx)

	if err := generateSimappAccountCoins(&ctx, tApp); err != nil {
		return &TestApp{}, sdk.Context{}, err
	}

	err := SetModuleAccountCoins(&ctx, tApp.BankKeeper, strategicreservetypes.SRPoolName, 150000000000000)
	if err != nil {
		return &TestApp{}, sdk.Context{}, err
	}

	err = SetModuleAccountCoins(&ctx, tApp.BankKeeper, strategicreservetypes.BetReserveName, 2000)
	if err != nil {
		return &TestApp{}, sdk.Context{}, err
	}

	err = SetModuleAccountCoins(&ctx, tApp.BankKeeper, strategicreservetypes.WinningsCollectorName, 20000)
	if err != nil {
		return &TestApp{}, sdk.Context{}, err
	}

	return tApp, ctx, nil
}

func setMinterParams(tApp *TestApp, ctx sdk.Context) {
	tApp.MintKeeper.SetParams(ctx, mintmoduletypes.DefaultParams())
	tApp.MintKeeper.SetMinter(ctx, mintmoduletypes.DefaultInitialMinter())
}

func generateSimappUsers() {
	createIncrementalAccounts(10)
	for i := 0; i < 10; i++ {
		prvKey := secp256k1.GenPrivKey()
		TestParamUsers[usernamePrefix+cast.ToString(i)] = TestUser{
			PrvKey:  prvKey,
			Address: sdk.AccAddress(prvKey.PubKey().Address()),
			Balance: 1000000000,
		}
	}
}

func generateSimappUserBalances() (balances []banktypes.Balance) {
	genTokens := sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)
	genCoin := sdk.NewCoin(params.DefaultBondDenom, genTokens)
	sdkgenCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000))
	for _, v := range TestParamUsers {
		balances = append(balances, banktypes.Balance{
			Address: v.Address.String(),
			Coins:   sdk.Coins{sdkgenCoin, genCoin},
		})
	}
	return balances
}

func generateSimappGenesisAccounts() (genAccs []authtypes.GenesisAccount) {
	for _, v := range TestParamUsers {
		genAccs = append(genAccs, &authtypes.BaseAccount{Address: v.Address.String()})
	}
	return genAccs
}

func generateSimappAccountCoins(ctx *sdk.Context, tApp *TestApp) error {
	for _, v := range TestParamUsers {
		if err := SetAccountCoins(ctx, tApp.BankKeeper, v.Address, v.Balance); err != nil {
			return err
		}
	}
	return nil
}

// SetAccountCoins sets the balance of accounts for testing
func SetAccountCoins(ctx *sdk.Context, k bankkeeper.Keeper, addr sdk.AccAddress, amount int64) error {
	coin := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(amount)))
	err := k.MintCoins(*ctx, mintmoduletypes.ModuleName, coin)
	if err != nil {
		return err
	}
	err = k.SendCoinsFromModuleToAccount(*ctx, mintmoduletypes.ModuleName, addr, coin)
	if err != nil {
		return err
	}
	return nil
}

// SetModuleAccountCoins sets the balance of accounts for testing
func SetModuleAccountCoins(ctx *sdk.Context, k bankkeeper.Keeper, moduleName string, amount int64) error {
	coin := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, sdk.NewInt(amount)))
	err := k.MintCoins(*ctx, mintmoduletypes.ModuleName, coin)
	if err != nil {
		return err
	}
	err = k.SendCoinsFromModuleToModule(*ctx, mintmoduletypes.ModuleName, moduleName, coin)
	if err != nil {
		return err
	}
	return nil
}

// DefaultConsensusParams parameters for tendermint consensus
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
		MaxBytes:        10000,
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

func stakingDefaultTestGenesis(tApp *TestApp) (*stakingtypes.GenesisState, []abci.ValidatorUpdate, banktypes.Balance) {
	defaultParams := stakingtypes.DefaultParams()
	defaultParams.BondDenom = params.DefaultBondDenom

	addr1 := TestParamUsers["user1"].Address
	addr2 := TestParamUsers["user2"].Address

	p1 := int64(8)
	p2 := int64(2)

	pks := simapp.CreateTestPubKeys(2)
	valConsPk1 := pks[0]
	valConsPk2 := pks[1]

	valPower1 := sdk.TokensFromConsensusPower(p1, sdk.DefaultPowerReduction)
	valPower2 := sdk.TokensFromConsensusPower(p2, sdk.DefaultPowerReduction)

	var validators []stakingtypes.Validator
	var delegations []stakingtypes.Delegation

	pk0, err := codectypes.NewAnyWithValue(valConsPk1)
	if err != nil {
		panic(err)
	}
	pk1, err := codectypes.NewAnyWithValue(valConsPk2)

	if err != nil {
		panic(err)
	}

	// initialize the validators
	bondedVal1 := stakingtypes.Validator{
		OperatorAddress: sdk.ValAddress(addr1).String(),
		ConsensusPubkey: pk0,
		Status:          stakingtypes.Bonded,
		Tokens:          valPower1,
		DelegatorShares: valPower1.ToDec(),
		Description:     stakingtypes.NewDescription("hoop", "", "", "", ""),
		Commission:      stakingtypes.NewCommission(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0)),
	}
	bondedVal2 := stakingtypes.Validator{
		OperatorAddress: sdk.ValAddress(addr2).String(),
		ConsensusPubkey: pk1,
		Status:          stakingtypes.Bonded,
		Tokens:          valPower2,
		DelegatorShares: valPower2.ToDec(),
		Description:     stakingtypes.NewDescription("bloop", "", "", "", ""),
		Commission:      stakingtypes.NewCommission(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0)),
	}

	// append new bonded validators to the list
	validators = append(validators, bondedVal1, bondedVal2)
	// mint coins in the bonded pool representing the validators coins

	var valudatorUpdates []abci.ValidatorUpdate
	valudatorUpdates = append(valudatorUpdates, bondedVal1.ABCIValidatorUpdate(sdk.DefaultPowerReduction))
	delegations = append(delegations, stakingtypes.Delegation{
		DelegatorAddress: addr1.String(),
		ValidatorAddress: bondedVal1.OperatorAddress,
		Shares:           bondedVal1.DelegatorShares,
	})
	valudatorUpdates = append(valudatorUpdates, bondedVal2.ABCIValidatorUpdate(sdk.DefaultPowerReduction))
	delegations = append(delegations, stakingtypes.Delegation{
		DelegatorAddress: addr2.String(),
		ValidatorAddress: bondedVal2.OperatorAddress,
		Shares:           bondedVal2.DelegatorShares,
	})

	moduleAddress := tApp.AccountKeeper.GetModuleAddress(stakingtypes.BondedPoolName)
	moduleBalance := banktypes.Balance{
		Address: moduleAddress.String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, valPower1.Add(valPower2))),
	}

	TestParamValidatorAddresses["val1"] = TestValidator{
		PubKey:      valConsPk1,
		Address:     bondedVal1.GetOperator(),
		ConsAddress: sdk.ConsAddress(valConsPk1.Address()),
		Power:       valPower1,
	}
	TestParamValidatorAddresses["val2"] = TestValidator{
		PubKey:      valConsPk2,
		Address:     bondedVal2.GetOperator(),
		ConsAddress: sdk.ConsAddress(valConsPk2.Address()),
		Power:       valPower2,
	}

	genesisState := stakingtypes.NewGenesisState(defaultParams, validators, delegations)
	return genesisState, valudatorUpdates, moduleBalance
}

// NewStakingHelper creates staking Handler wrapper for tests
func NewStakingHelper(t *testing.T, ctx sdk.Context, k stakingKeeper.Keeper) *teststaking.Helper {
	helper := teststaking.NewHelper(t, ctx, k)
	helper.Commission = validatorDefaultCommission()
	helper.Denom = params.DefaultBondDenom
	return helper
}

func validatorDefaultCommission() stakingtypes.CommissionRates {
	return stakingtypes.NewCommissionRates(sdk.MustNewDecFromStr("0.1"), sdk.MustNewDecFromStr("0.2"), sdk.MustNewDecFromStr("0.01"))
}
