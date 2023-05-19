package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

func (k Keeper) RewardUser(
	ctx sdk.Context, userAddress sdk.AccAddress, payout sdk.Int, uniqueLock string,
) (err error) {
	err = k.transferFundsFromAccountToModule(ctx, userAddress, types.IncentiveReservePool, payout)
	if err != nil {
		return
	}
	return nil
}
