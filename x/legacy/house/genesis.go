package house

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/house/keeper"
	"github.com/sge-network/sge/x/legacy/house/types"
)

// InitGenesis sets the deposits and parameters for the provided keeper.
func InitGenesis(ctx sdk.Context, keeper *keeper.Keeper, data *types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, deposit := range data.DepositList {
		keeper.SetDeposit(ctx, deposit)
	}

	for _, withdrawal := range data.WithdrawalList {
		keeper.SetWithdrawal(ctx, withdrawal)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the params and deposits found in the keeper.
func ExportGenesis(ctx sdk.Context, k *keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	var err error
	genesis.DepositList, err = k.GetAllDeposits(ctx)
	if err != nil {
		panic(err)
	}

	genesis.WithdrawalList, err = k.GetAllWithdrawals(ctx)
	if err != nil {
		panic(err)
	}

	return genesis
}
