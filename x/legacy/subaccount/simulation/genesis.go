package simulation

// DONTCOVER

import (
	"encoding/json"
	//#nosec
	"math/rand"

	"github.com/spf13/cast"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/sge-network/sge/x/legacy/subaccount/types"
)

// Simulation parameter constants
const (
	HouseParticipationFee = "HouseParticipationFee"
	MinDeposit            = "MinDeposit"
)

// GenHouseParticipationFee randomized batch settlement count
func GenHouseParticipationFee(r *rand.Rand) sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(cast.ToInt64(r.Intn(99)))
}

// GenMinDeposit randomized house by uid query count
func GenMinDeposit(r *rand.Rand) sdkmath.Int {
	return sdkmath.NewInt(cast.ToInt64(r.Intn(99)))
}

// RandomizedGenState generates a random GenesisState for house
func RandomizedGenState(simState *module.SimulationState) {
	var houseParticipationFee sdkmath.LegacyDec
	simState.AppParams.GetOrGenerate(
		HouseParticipationFee, &houseParticipationFee, simState.Rand,
		func(r *rand.Rand) { houseParticipationFee = GenHouseParticipationFee(r) },
	)

	var minDeposit sdkmath.Int
	simState.AppParams.GetOrGenerate(
		MinDeposit, &minDeposit, simState.Rand,
		func(r *rand.Rand) { minDeposit = GenMinDeposit(r) },
	)

	defaultGenesis := types.DefaultGenesis()

	_, err := json.MarshalIndent(&defaultGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(defaultGenesis)
}
