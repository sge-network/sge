package bet

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/sge-network/sge/testutil/sample"
	betsimulation "github.com/sge-network/sge/x/bet/simulation"
	"github.com/sge-network/sge/x/bet/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = betsimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	//nolint:gosec
	opWeightMsgPlaceBet = "op_weight_msg_bet"
	// TODO: Determine the simulation weight value
	defaultWeightMsgPlaceBet int = 100
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	betGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		BetList: []types.Bet{
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
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&betGenesis)
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

	var weightMsgPlaceBet int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgPlaceBet, &weightMsgPlaceBet, nil,
		func(_ *rand.Rand) {
			weightMsgPlaceBet = defaultWeightMsgPlaceBet
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgPlaceBet,
		betsimulation.SimulateMsgPlaceBet(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	return operations
}
