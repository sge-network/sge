package simulation

// DONTCOVER

import (
	"fmt"
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sge-network/sge/x/market/types"
)

const (
	keyMinBetAmount      = "MinBetAmount"
	keyMinBetFee         = "MinBetFee"
	keyMaxBetFee         = "MaxBetFee"
	keyMaxSrContribution = "MaxSrContribution"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyMinBetAmount,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMinBetAmount(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyMinBetFee,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMinBetFee(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyMaxBetFee,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxBetFee(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, keyMaxSrContribution,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GenMaxSrContribution(r))
			},
		),
	}
}
