package simulation

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/sge-network/sge/x/legacy/reward/types"
)

// Simulation operation weights constants
const (
	DefaultWeightMsgUpdateParams int = 100

	OpWeightMsgUpdateParams = "op_weight_msg_update_params" //#nosec
)

// ProposalMsgs defines the module weighted proposals' contents
func ProposalMsgs() []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			OpWeightMsgUpdateParams,
			DefaultWeightMsgUpdateParams,
			SimulateMsgUpdateParams,
		),
	}
}

// SimulateMsgUpdateParams returns a random MsgUpdateParams
func SimulateMsgUpdateParams(_ *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	// use the default gov module account address as authority
	var authority sdk.AccAddress = address.Module("gov")

	params := types.DefaultParams()

	return &types.MsgUpdateParams{
		Authority: authority.String(),
		Params:    params,
	}
}
