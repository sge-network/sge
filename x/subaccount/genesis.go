package subaccount

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/keeper"
	"github.com/sge-network/sge/x/subaccount/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)
	if genState.SubaccountId != 0 {
		k.SetID(ctx, genState.SubaccountId)
	}
	for _, acc := range genState.Subaccounts {
		owner := sdk.MustAccAddressFromBech32(acc.Owner)
		addr := sdk.MustAccAddressFromBech32(acc.Address)
		k.SetSubAccountOwner(ctx, addr, owner)
		k.SetLockedBalances(ctx, addr, acc.LockedBalances)
		k.SetBalance(ctx, addr, acc.Balance)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:       k.GetParams(ctx),
		SubaccountId: k.Peek(ctx),
		Subaccounts:  k.GetAllSubaccounts(ctx),
	}

}
