package dvm

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/dvm/keeper"
)

// EndBlocker settles the active bets of resolved sport events
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	err := k.FinishProposals(ctx)
	k.Logger(ctx).Error(fmt.Sprintf("end block number %d error: %s", ctx.BlockHeight(), err))
}
