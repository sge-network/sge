package bet

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
	betsimulation "github.com/sge-network/sge/x/bet/simulation"
	"github.com/sge-network/sge/x/bet/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = betsimulation.FindAccount
	_ = simtestutil.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	//#nosec
	opWeightMsgWager = "op_weight_msg_wager"
	// TODO: Determine the simulation weight value
	defaultWeightMsgWager int = 100
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	betsimulation.RandomizedGenState(simState)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {
	sdr[types.StoreKey] = betsimulation.NewDecodeStore(am.cdc)
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgWager int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgWager, &weightMsgWager, nil,
		func(_ *rand.Rand) {
			weightMsgWager = defaultWeightMsgWager
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgWager,
		betsimulation.SimulateMsgWager(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	return operations
}
