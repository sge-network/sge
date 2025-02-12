package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/legacy/house/types"
)

// getDepositStore gets the store containing all deposits.
func (k Keeper) getDepositStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.DepositKeyPrefix)
}

// getWithdrawalStore gets the store containing all withdrawals.
func (k Keeper) getWithdrawalStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.WithdrawalKeyPrefix)
}
