package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

type Hook interface {
	AfterBettorWin(ctx sdk.Context, bettor sdk.AccAddress, originalAmount, profit sdk.Int)
	AfterBettorLoss(ctx sdk.Context, bettor sdk.AccAddress, originalAmount sdk.Int)
	AfterBettorRefund(ctx sdk.Context, bettor sdk.AccAddress, originalAmount, fee sdk.Int)
	AfterHouseWin(ctx sdk.Context, house sdk.AccAddress, originalAmount, profit sdk.Int, fee *sdk.Int)
	AfterHouseLoss(ctx sdk.Context, house sdk.AccAddress, originalAmount sdk.Int, lostAmt sdk.Int, fee *sdk.Int)
	AfterHouseRefund(ctx sdk.Context, house sdk.AccAddress, originalAmount, fee sdk.Int)
}

func (k *Keeper) RegisterHook(hook Hook) {
	k.hooks = append(k.hooks, hook)
}
