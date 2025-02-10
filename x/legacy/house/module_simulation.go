package house

import (
	//#nosec

	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sge-network/sge/testutil/sample"
	"github.com/sge-network/sge/x/legacy/bet/types"
	housesimulation "github.com/sge-network/sge/x/legacy/house/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = housesimulation.FindAccount
	_ = simtestutil.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	housesimulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	sdr[types.StoreKey] = housesimulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}
