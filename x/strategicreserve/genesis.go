package strategicreserve

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/strategicreserve/keeper"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetReserver(ctx, genState.Reserver)
	k.SetParams(ctx, genState.Params)

	// check if the sr and bet_reserve pools accounts exists
	srPool := k.GetSRPoolAcc(ctx)
	if srPool == nil {
		panic(fmt.Sprintf(consts.ErrModuleAccountHasNotBeenSet, types.SRPoolName))
	}

	if k.GetSRPoolBalance(ctx).Amount.ToDec() == sdk.ZeroDec() {
		panic(fmt.Sprintf("%s module account has no balance", types.SRPoolName))
	}

	betReserveAcc := k.GetBetReserveAcc(ctx)
	if betReserveAcc == nil {
		panic(fmt.Sprintf(consts.ErrModuleAccountHasNotBeenSet, types.BetReserveName))
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.Reserver = k.GetReserver(ctx)

	return genesis
}
