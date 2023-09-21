package simulation

// DONTCOVER

import (
	"encoding/json"
	//#nosec
	"math/rand"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/sge-network/sge/x/bet/types"
	"github.com/spf13/cast"
)

// Simulation parameter constants
const (
	BatchSettlementCount  = "BatchSettlementCount"
	MaxBetByUIDQueryCount = "MaxBetByUidQueryCount"
	MinAmount             = "MinAmount"
	BetFee                = "BetFee"
)

// GenBatchSettlementCount randomized batch settlement count
func GenBatchSettlementCount(r *rand.Rand) uint32 {
	return cast.ToUint32(r.Intn(99))
}

// GenMaxBetByUIDQueryCount randomized bet by uid query count
func GenMaxBetByUIDQueryCount(r *rand.Rand) uint32 {
	return cast.ToUint32(r.Intn(99))
}

// GenMinAmount randomized min bet amount
func GenMinAmount(r *rand.Rand) sdkmath.Int {
	return sdkmath.NewInt(int64(r.Intn(99)))
}

// GenFee randomized min bet fee
func GenFee(r *rand.Rand) sdkmath.Int {
	return sdkmath.NewInt(int64(r.Intn(99)))
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

	var minAmount sdkmath.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinAmount, &minAmount, simState.Rand,
		func(r *rand.Rand) { minAmount = GenMinAmount(r) },
	)

	var minBetFee sdkmath.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, BetFee, &minBetFee, simState.Rand,
		func(r *rand.Rand) { minBetFee = GenFee(r) },
	)

	defaultGenesis := types.DefaultGenesis()

	_, err := json.MarshalIndent(&defaultGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(defaultGenesis)
}
