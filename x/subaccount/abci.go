package subaccount

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/keeper"
)

// EndBlocker settles the active bets of resolved markets
func EndBlocker(_ sdk.Context, _ keeper.Keeper) {}
