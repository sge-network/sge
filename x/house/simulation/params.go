package simulation

// DONTCOVER

import (
	"fmt"
	//#nosec
	"math/rand"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sge-network/sge/x/house/types"
)

const (
	keyMinDeposit = "MinDeposit"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(_ *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyMinDeposit,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenMinDeposit(r))
			},
		),
	}
}
