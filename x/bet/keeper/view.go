package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/bet/types"
)

// getBetStore returns bet store ready for iterating
func (k Keeper) getBetStore(ctx sdk.Context) prefix.Store {
	betStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BetListPrefix)
	return betStore
}
