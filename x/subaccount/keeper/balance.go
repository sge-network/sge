package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/subaccount/types"
)

// getBalances returns the balance, unlocked balance and bank balance of a subaccount
func (k Keeper) getBalances(sdkContext sdk.Context, subaccountAddr sdk.AccAddress, params types.Params) (types.Balance, sdk.Int, sdk.Coin) {
	balance, exists := k.GetBalance(sdkContext, subaccountAddr)
	if !exists {
		panic("data corruption: subaccount exists but balance does not")
	}
	unlockedBalance := k.GetUnlockedBalance(sdkContext, subaccountAddr)
	bankBalance := k.bankKeeper.GetBalance(sdkContext, subaccountAddr, params.LockedBalanceDenom)

	return balance, unlockedBalance, bankBalance
}
