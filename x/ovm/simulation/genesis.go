package simulation

// DONTCOVER

import (
	"encoding/json"
	//#nosec

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/sge-network/sge/x/ovm/types"
)

// RandomizedGenState generates a random GenesisState for ovm
func RandomizedGenState(simState *module.SimulationState) {
	defaultGenesis := types.DefaultGenesis()

	_, err := json.MarshalIndent(&defaultGenesis, "", " ")
	if err != nil {
		panic(err)
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(defaultGenesis)
}
