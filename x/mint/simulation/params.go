package simulation

// DONTCOVER

import (
	"fmt"
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sge-network/sge/x/mint/types"
)

const (
	keyBlocksPerYear = "BlocksPerYear"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(_ *rand.Rand) []simtypes.LegacyParamChange {
	return []simtypes.LegacyParamChange{
		simulation.NewSimLegacyParamChange(types.ModuleName, keyBlocksPerYear,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenBlocksPerYear(r))
			},
		),
	}
}
