package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

func (k Keeper) RewardUser(
	ctx sdk.Context, userAddress sdk.AccAddress, payout sdk.Int, uniqueLock string,
) (err error) {
	fmt.Println("Calling transferFundsFromModuleToAccount")
	accAddr := k.accountKeeper.GetModuleAddress(types.IncentiveReservePool)
	bech, _ := sdk.AccAddressFromHex(accAddr.String())
	fmt.Println("addres: ", accAddr, bech)
	err = k.transferFundsFromModuleToAccount(ctx, types.IncentiveReservePool, userAddress, payout)
	if err != nil {
		return
	}
	return nil
}
