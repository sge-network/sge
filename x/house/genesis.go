package house

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/house/keeper"
	"github.com/sge-network/sge/x/house/types"
)

// InitGenesis sets the deposits and parameters for the provided keeper.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	// TODO
	keeper.SetParams(ctx, data.Params)

	for _, deposit := range data.Deposits {
		keeper.SetDeposit(ctx, deposit)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the params and deposits found in the keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	// TODO
	return &types.GenesisState{
		Params:   keeper.GetParams(ctx),
		Deposits: keeper.GetAllDeposits(ctx),
	}
}

// ValidateGenesis validates the provided house genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate deposits)
func ValidateGenesis(data *types.GenesisState) error {
	// TODO

	return data.Params.Validate()
}
