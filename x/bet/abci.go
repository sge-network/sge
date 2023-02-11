package bet

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/keeper"
)

// EndBlocker settles the active bets of resolved sport events
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	err := k.BatchSportEventSettlements(ctx)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("end block no %d failed : %s", ctx.BlockHeight(), err.Error()))
	}
}
