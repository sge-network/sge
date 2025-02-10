package orderbook

import (
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sge-network/sge/testutil/sample"
	orderbooksimulation "github.com/sge-network/sge/x/legacy/orderbook/simulation"
	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = orderbooksimulation.FindAccount
	_ = simtestutil.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	//#nosec
	opWeightMsgAdd          = "op_weight_msg_create_chain"
	defaultWeightMsgAdd int = 100
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	orderbooksimulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(sdr simtypes.StoreDecoderRegistry) {
	sdr[types.StoreKey] = orderbooksimulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgAdd int
	simState.AppParams.GetOrGenerate(opWeightMsgAdd, &weightMsgAdd, nil,
		func(_ *rand.Rand) {
			weightMsgAdd = defaultWeightMsgAdd
		},
	)

	return operations
}
