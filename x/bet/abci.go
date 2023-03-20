package bet

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/keeper"
)

// EndBlocker settles the active bets of resolved markets
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	err := k.BatchMarketSettlements(ctx)
	if err != nil {
		panic(fmt.Sprintf("end block no %d failed : %s", ctx.BlockHeight(), err.Error()))
	}
}
