package bet

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/keeper"
	"github.com/sge-network/sge/x/bet/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {

	k.SetBetStats(ctx, genState.Stats)

	// Set all the bet
	for _, elem := range genState.BetList {
		var id uint64
		for _, uid2ID := range genState.Uid2IdList {
			if uid2ID.UID == elem.UID {
				id = uid2ID.ID
			}
		}

		if id == 0 {
			// this means the imported genesis is broken because there is no corresponding
			// id mapped to the uid
			panic(fmt.Errorf(types.ErrTextInitGenesisFailedBecauseOfMissingBetID, elem.UID))
		}

		k.SetBet(ctx, elem, id)
	}

	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	var err error
	genesis.BetList, err = k.GetBets(ctx)

	if err != nil {
		panic(err)
	}

	return genesis
}
