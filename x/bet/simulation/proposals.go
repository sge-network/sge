package simulation

import (
	"math/rand"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/x/bet/types"
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
func SimulateMsgUpdateParams(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) sdk.Msg {
	// use the default gov module account address as authority
	var authority sdk.AccAddress = address.Module("gov")

	params := types.DefaultParams()
	params.BatchSettlementCount = cast.ToUint32(r.Intn(10000))
	params.MaxBetByUidQueryCount = cast.ToUint32(r.Intn(1000))
	params.Constraints = types.Constraints{
		MinAmount: sdkmath.NewInt(cast.ToInt64(r.Intn(1000))),
		Fee:       sdkmath.NewInt(cast.ToInt64(r.Intn(99))),
	}

	return &types.MsgUpdateParams{
		Authority: authority.String(),
		Params:    params,
	}
}
