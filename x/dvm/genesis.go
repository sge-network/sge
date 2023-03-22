package dvm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/keeper"
	"github.com/sge-network/sge/x/dvm/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set the key vault
	k.SetKeyVault(ctx, genState.KeyVault)

	// set pubkeys change proposals
	for i := range genState.PubkeysChangeProposals {
		proposal := genState.PubkeysChangeProposals[i]
		k.SetPubkeysChangeProposal(ctx, proposal)
	}

	// set proposal statistics
	k.SetProposalStats(ctx, genState.ProposalStats)

	// set parameters
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// load the key vault
	keys, found := k.GetKeyVault(ctx)
	if found {
		genesis.KeyVault = keys
	}

	// load pubkeys change proposals
	proposals, err := k.GetAllPubkeysChangeProposals(ctx)
	if err != nil {
		panic(err)
	}
	genesis.PubkeysChangeProposals = append(genesis.PubkeysChangeProposals, proposals...)

	// load proposal statistics
	genesis.ProposalStats = k.GetProposalStats(ctx)

	// load the default Params
	genesis.Params = k.GetParams(ctx)

	return genesis
}
