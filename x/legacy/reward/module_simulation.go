package reward

import (
	//#nosec
	"math/rand"

	"github.com/google/uuid"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sge-network/sge/testutil/sample"
	rewardsimulation "github.com/sge-network/sge/x/legacy/reward/simulation"
	"github.com/sge-network/sge/x/legacy/reward/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = rewardsimulation.FindAccount
	_ = simtestutil.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	//#nosec
	opWeightMsgCreateCampaign = "op_weight_msg_campaign"
	// TODO: Determine the simulation weight value
	//#nosec
	defaultWeightMsgCreateCampaign int = 100

	//#nosec
	opWeightMsgUpdateCampaign = "op_weight_msg_campaign"
	// TODO: Determine the simulation weight value
	defaultWeightMsgUpdateCampaign int = 100

	//#nosec
	opWeightMsgDeleteCampaign = "op_weight_msg_campaign"
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeleteCampaign int = 100
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	rewardGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		CampaignList: []types.Campaign{
			{
				Creator: sample.AccAddress(),
				UID:     uuid.NewString(),
			},
			{
				Creator: sample.AccAddress(),
				UID:     uuid.NewString(),
			},
		},
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&rewardGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgCreateCampaign int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateCampaign, &weightMsgCreateCampaign, nil,
		func(_ *rand.Rand) {
			weightMsgCreateCampaign = defaultWeightMsgCreateCampaign
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateCampaign,
		rewardsimulation.SimulateMsgCreateCampaign(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgUpdateCampaign int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateCampaign, &weightMsgUpdateCampaign, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateCampaign = defaultWeightMsgUpdateCampaign
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateCampaign,
		rewardsimulation.SimulateMsgUpdateCampaign(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeleteCampaign int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteCampaign, &weightMsgDeleteCampaign, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteCampaign = defaultWeightMsgDeleteCampaign
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteCampaign,
		rewardsimulation.SimulateMsgDeleteCampaign(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	return operations
}
