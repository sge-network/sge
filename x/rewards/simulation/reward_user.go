package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/sge-network/sge/x/rewards/keeper"
	"github.com/sge-network/sge/x/rewards/types"
)

func SimulateMsgRewardUser(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgRewardUser{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the RewardUser simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "RewardUser simulation not implemented"), nil, nil
	}
}
