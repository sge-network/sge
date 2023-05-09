package simulation

// DONTCOVER

import (
	"fmt"
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sge-network/sge/x/bet/types"
)

const (
	keyMaxBetByUIDQueryCount = "MaxBetByUidQueryCount"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyMaxBetByUIDQueryCount,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxBetByUIDQueryCount(r))
			},
		),
	}
}
