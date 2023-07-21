package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

type Hook interface {
	AfterBettorWin(ctx sdk.Context, bettor sdk.AccAddress, originalAmount, profit sdk.Int)
	AfterBettorLoss(ctx sdk.Context, bettor sdk.AccAddress, originalAmount sdk.Int)
	AfterBettorRefund(ctx sdk.Context, bettor sdk.AccAddress, originalAmount, fee sdk.Int)
}

func (k *Keeper) RegisterHook(hook Hook) {
	k.hooks = append(k.hooks, hook)
}
