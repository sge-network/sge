package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/house/types"
)

// getDepositsStore gets the store containing all deposits.
func (k Keeper) getDepositsStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.DepositKeyPrefix)
}

// getWithdrawalsStore gets the store containing all withdrawals.
func (k Keeper) getWithdrawalsStore(ctx sdk.Context) prefix.Store {
	store := ctx.KVStore(k.storeKey)
	return prefix.NewStore(store, types.WithdrawalKeyPrefix)
}
