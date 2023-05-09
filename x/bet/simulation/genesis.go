package simulation

// DONTCOVER

import (
	"encoding/json"
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/spf13/cast"
)

// Simulation parameter constants
const (
	BatchSettlementCount  = "BatchSettlementCount"
	MaxBetByUIDQueryCount = "MaxBetByUidQueryCount"
)

// GenBatchSettlementCount randomized bathc settlement count
func GenBatchSettlementCount(r *rand.Rand) uint32 {
	return cast.ToUint32(r.Intn(99))
}

// GenMaxBetByUIDQueryCount randomized bet by uid query count
func GenMaxBetByUIDQueryCount(r *rand.Rand) uint32 {
	return cast.ToUint32(r.Intn(99))
}

// RandomizedGenState generates a random GenesisState for bet
func RandomizedGenState(simState *module.SimulationState) {
	var batchSettlementCount uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, BatchSettlementCount, &batchSettlementCount, simState.Rand,
		func(r *rand.Rand) { batchSettlementCount = GenBatchSettlementCount(r) },
	)

	var maxBetByUIDQueryCount uint32
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxBetByUIDQueryCount, &maxBetByUIDQueryCount, simState.Rand,
		func(r *rand.Rand) { maxBetByUIDQueryCount = GenMaxBetByUIDQueryCount(r) },
	)

	defaultGenesis := types.DefaultGenesis()

	_, err := json.MarshalIndent(&defaultGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(defaultGenesis)
}
