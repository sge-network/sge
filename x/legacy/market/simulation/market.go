package simulation

import (
	//#nosec
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/sge-network/sge/x/legacy/market/keeper"
	"github.com/sge-network/sge/x/legacy/market/types"
)

// SimulateMsgAdd simulates the add market flow
func SimulateMsgAdd(
	_ types.AccountKeeper,
	_ types.BankKeeper,
	_ keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, _ *baseapp.BaseApp, _ sdk.Context, accs []simtypes.Account, _ string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgAdd{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the AddMarket simulation

		return simtypes.NoOpMsg(
			types.ModuleName,
			msg.Type(),
			"AddMarket simulation not implemented",
		), nil, nil
	}
}

// SimulateMsgResolve simulates the resolve market flow
func SimulateMsgResolve(
	_ types.AccountKeeper,
	_ types.BankKeeper,
	_ keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, _ *baseapp.BaseApp, _ sdk.Context, accs []simtypes.Account, _ string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgResolve{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the ResolveMarket simulation

		return simtypes.NoOpMsg(
			types.ModuleName,
			msg.Type(),
			"ResolveMarket simulation not implemented",
		), nil, nil
	}
}
