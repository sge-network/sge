package simulation

// DONTCOVER

import (
	"fmt"
	//#nosec
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sge-network/sge/x/orderbook/types"
)

const (
	keyBatchSettlementCount       = "BatchSettlementCount"
	keyMaxOrderBookParticipations = "MaxOrderBookParticipations"
	keyRequeueThreshold           = "RequeueThreshold"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(_ *rand.Rand) []simtypes.LegacyParamChange {
	return []simtypes.LegacyParamChange{
		simulation.NewSimLegacyParamChange(types.ModuleName, keyBatchSettlementCount,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenBatchSettlementCount(r))
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, keyMaxOrderBookParticipations,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxOrderBookParticipations(r))
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, keyRequeueThreshold,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenRequeueThreshold(r))
			},
		),
	}
}
