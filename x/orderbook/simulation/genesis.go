package simulation

// DONTCOVER

import (
	"encoding/json"
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/sge-network/sge/x/orderbook/types"
	"github.com/spf13/cast"
)

// Simulation parameter constants
const (
	BatchSettlementCount       = "BatchSettlementCount"
	MaxOrderBookParticipations = "MaxOrderBookParticipations"
	RequeueThreshold           = "RequeueThreshold"
)

// GenBatchSettlementCount randomized min bet amount
func GenBatchSettlementCount(r *rand.Rand) uint64 {
	return cast.ToUint64(r.Intn(99))
}

// GenMaxOrderBookParticipations randomized min bet fee
func GenMaxOrderBookParticipations(r *rand.Rand) uint64 {
	return cast.ToUint64(r.Intn(99))
}

// GenRequeueThreshold randomized max bet fee
func GenRequeueThreshold(r *rand.Rand) uint64 {
	return cast.ToUint64(r.Intn(99))
}

// RandomizedGenState generates a random GenesisState for market
func RandomizedGenState(simState *module.SimulationState) {
	var batchSettlementCount uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, BatchSettlementCount, &batchSettlementCount, simState.Rand,
		func(r *rand.Rand) { batchSettlementCount = GenBatchSettlementCount(r) },
	)

	var maxOrderBookParticipations uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxOrderBookParticipations, &maxOrderBookParticipations, simState.Rand,
		func(r *rand.Rand) { maxOrderBookParticipations = GenMaxOrderBookParticipations(r) },
	)

	var requeueThreshold uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, RequeueThreshold, &requeueThreshold, simState.Rand,
		func(r *rand.Rand) { requeueThreshold = GenRequeueThreshold(r) },
	)

	defaultGenesis := types.DefaultGenesis()

	_, err := json.MarshalIndent(&defaultGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(defaultGenesis)
}
