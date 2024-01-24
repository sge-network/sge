package market

import (
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sge-network/sge/testutil/sample"
	marketsimulation "github.com/sge-network/sge/x/market/simulation"
	"github.com/sge-network/sge/x/market/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = marketsimulation.FindAccount
	_ = simtestutil.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	//#nosec
	opWeightMsgAdd          = "op_weight_msg_create_chain"
	defaultWeightMsgAdd int = 100

	//#nosec
	opWeightMsgResolve          = "op_weight_msg_create_chain"
	defaultWeightMsgResolve int = 100

	//#nosec
	opWeightMsgUpdate          = "op_weight_msg_create_chain"
	defaultWeightMsgUpdate int = 100
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	marketsimulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[types.StoreKey] = marketsimulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgAdd int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAdd, &weightMsgAdd, nil,
		func(_ *rand.Rand) {
			weightMsgAdd = defaultWeightMsgAdd
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAdd,
		marketsimulation.SimulateMsgAdd(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgResolve int
	simState.AppParams.GetOrGenerate(
		simState.Cdc,
		opWeightMsgResolve,
		&weightMsgResolve,
		nil,
		func(_ *rand.Rand) {
			weightMsgResolve = defaultWeightMsgResolve
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgResolve,
		marketsimulation.SimulateMsgResolve(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdate int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdate, &weightMsgUpdate, nil,
		func(_ *rand.Rand) {
			weightMsgUpdate = defaultWeightMsgUpdate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdate,
		marketsimulation.SimulateMsgUpdate(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	return operations
}
