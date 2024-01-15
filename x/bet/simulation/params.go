package simulation

// DONTCOVER

import (
	"fmt"
	//#nosec
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sge-network/sge/x/bet/types"
)

const (
	keyBatchSettlementCount  = "BatchSettlementCount"
	keyMaxBetByUIDQueryCount = "MaxBetByUidQueryCount"
	keyMinAmount             = "MinAmount"
	keyMinFee                = "Fee"
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
		simulation.NewSimLegacyParamChange(types.ModuleName, keyMaxBetByUIDQueryCount,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxBetByUIDQueryCount(r))
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, keyMinAmount,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMinAmount(r))
			},
		),
		simulation.NewSimLegacyParamChange(types.ModuleName, keyMinFee,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenFee(r))
			},
		),
	}
}
