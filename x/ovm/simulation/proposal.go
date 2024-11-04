package simulation

import (
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/sge-network/sge/x/ovm/keeper"
	"github.com/sge-network/sge/x/ovm/types"
)

// SimulateMsgChangePubkeysListProposal simulates MsgChangePubkeysListProposal message registration
func SimulateMsgChangePubkeysListProposal(
	_ types.AccountKeeper,
	_ types.BankKeeper,
	_ keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, _ *baseapp.BaseApp, _ sdk.Context, accs []simtypes.Account, _ string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgSubmitPubkeysChangeProposalRequest{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the change pubkeys proposal simulation

		return simtypes.NoOpMsg(
			types.ModuleName,
			msg.Type(),
			"ChangePubkeysListProposal simulation not implemented",
		), nil, nil
	}
}
