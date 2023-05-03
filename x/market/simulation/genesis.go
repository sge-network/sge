package simulation

// DONTCOVER

import (
	"encoding/json"
	//#nosec
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/sge-network/sge/x/market/types"
)

// Simulation parameter constants
const (
	MinBetAmount      = "MinBetAmount"
	MinBetFee         = "MinBetFee"
	MaxBetFee         = "MaxBetFee"
	MaxSrContribution = "MaxSrContribution"
)

// GenMinBetAmount randomized min bet amount
func GenMinBetAmount(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(r.Intn(99)))
}

// GenMinBetFee randomized min bet fee
func GenMinBetFee(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(r.Intn(99)))
}

// GenMaxBetFee randomized max bet fee
func GenMaxBetFee(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(r.Intn(99)))
}

// GenMaxSrContribution randomized st contribution
func GenMaxSrContribution(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(r.Intn(99)))
}

// RandomizedGenState generates a random GenesisState for market
func RandomizedGenState(simState *module.SimulationState) {
	var minBetAmount sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinBetAmount, &minBetAmount, simState.Rand,
		func(r *rand.Rand) { minBetAmount = GenMinBetAmount(r) },
	)

	var minBetFee sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinBetFee, &minBetFee, simState.Rand,
		func(r *rand.Rand) { minBetFee = GenMinBetFee(r) },
	)

	var maxBetFee sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxBetFee, &maxBetFee, simState.Rand,
		func(r *rand.Rand) { maxBetFee = GenMaxBetFee(r) },
	)

	var maxSrContribution sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxSrContribution, &maxSrContribution, simState.Rand,
		func(r *rand.Rand) { maxSrContribution = GenMaxSrContribution(r) },
	)

	defaultGenesis := types.DefaultGenesis()

	_, err := json.MarshalIndent(&defaultGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(defaultGenesis)
}
