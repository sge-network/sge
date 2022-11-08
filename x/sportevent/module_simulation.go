package sportevent

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/sge-network/sge/testutil/sample"
	sporteventsimulation "github.com/sge-network/sge/x/sportevent/simulation"
	"github.com/sge-network/sge/x/sportevent/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = sporteventsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

// nolint
const (
	opWeightMsgAddEvent = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddEvent int = 100

	opWeightMsgResolveEvent = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgResolveEvent int = 100

	opWeightMsgUpdateEvent = "op_weight_msg_create_chain"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateEvent int = 100
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	sporteventGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		SportEventList: []types.SportEvent{
			{
				Creator: sample.AccAddress(),
				UID:     "0",
			},
			{
				Creator: sample.AccAddress(),
				UID:     "1",
			},
		},
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&sporteventGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgAddEvent int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddEvent, &weightMsgAddEvent, nil,
		func(_ *rand.Rand) {
			weightMsgAddEvent = defaultWeightMsgAddEvent
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddEvent,
		sporteventsimulation.SimulateMsgAddEvent(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgResolveEvent int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgResolveEvent, &weightMsgResolveEvent, nil,
		func(_ *rand.Rand) {
			weightMsgResolveEvent = defaultWeightMsgResolveEvent
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgResolveEvent,
		sporteventsimulation.SimulateMsgResolveEvent(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateEvent int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateEvent, &weightMsgUpdateEvent, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateEvent = defaultWeightMsgUpdateEvent
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateEvent,
		sporteventsimulation.SimulateMsgUpdateEvent(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	return operations
}
