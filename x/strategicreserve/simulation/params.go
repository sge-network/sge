package simulation

// DONTCOVER

import (
	"fmt"
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

const (
	keyBatchSettlementCount       = "BatchSettlementCount"
	keyMaxOrderBookParticipations = "MaxOrderBookParticipations"
	keyRequeueThreshold           = "RequeueThreshold"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyBatchSettlementCount,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenBatchSettlementCount(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyMaxOrderBookParticipations,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxOrderBookParticipations(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyRequeueThreshold,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenRequeueThreshold(r))
			},
		),
	}
}
