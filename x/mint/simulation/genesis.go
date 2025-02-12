package simulation

// DONTCOVER

import (
	"encoding/json"
	//#nosec
	"math/rand"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/app/params"
	"github.com/sge-network/sge/x/mint/types"
)

// Inflation Simulation parameter constants
const (
	Inflation = "inflation"
)

// GenInflation randomized Inflation
func GenInflation(r *rand.Rand) sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(int64(r.Intn(99))).QuoInt64(100)
}

// GenBlocksPerYear randomized BlocksPerYear
func GenBlocksPerYear(_ *rand.Rand) sdkmath.LegacyDec {
	return sdkmath.LegacyMustNewDecFromStr(cast.ToString(types.BlocksPerYear))
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation sdkmath.LegacyDec
	simState.AppParams.GetOrGenerate(
		Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = GenInflation(r) },
	)

	// params
	mintDenom := params.DefaultBondDenom
	params := types.NewParams(
		mintDenom,
		types.BlocksPerYear,
		types.DefaultExcludeAmount,
		types.DefaultParams().Phases,
	)

	mintGenesis := types.NewGenesisState(types.InitialMinter(inflation), params)

	_, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
