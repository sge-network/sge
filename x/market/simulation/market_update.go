package simulation

import (
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sge-network/sge/x/market/keeper"
	"github.com/sge-network/sge/x/market/types"
)

// SimulateMsgUpdateMarket simulates update market message
func SimulateMsgUpdateMarket(
	_ types.AccountKeeper,
	_ types.BankKeeper,
	_ keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgUpdateMarket{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the UpdateMarket simulation

		return simtypes.NoOpMsg(
			types.ModuleName,
			msg.Type(),
			"UpdateMarket simulation not implemented",
		), nil, nil
	}
}
