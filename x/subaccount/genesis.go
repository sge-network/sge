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
	genesis := new(types.GenesisState)

	genesis.Params = k.GetParams(ctx)
	genesis.SubaccountId = k.Peek(ctx)

	k.IterateSubaccounts(ctx, func(subAccountAddress sdk.AccAddress, ownerAddress sdk.AccAddress) (stop bool) {
		balance, exists := k.GetBalance(ctx, subAccountAddress)
		if !exists {
			panic("subaccount balance does not exist")
		}
		lockedBalances := k.GetLockedBalances(ctx, subAccountAddress)
		genesis.Subaccounts = append(genesis.Subaccounts, types.GenesisSubaccount{
			Address:        subAccountAddress.String(),
			Owner:          ownerAddress.String(),
			Balance:        balance,
			LockedBalances: lockedBalances,
		})
		return false
	})

	return genesis

}
