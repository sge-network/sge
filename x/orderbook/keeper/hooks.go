package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Hook interface {
	AfterBettorWin(ctx sdk.Context, bettor sdk.AccAddress, originalAmount, profit math.Int)
	AfterBettorLoss(ctx sdk.Context, bettor sdk.AccAddress, originalAmount math.Int)
	AfterBettorRefund(ctx sdk.Context, bettor sdk.AccAddress, originalAmount, fee math.Int)

	AfterHouseWin(ctx sdk.Context, house sdk.AccAddress, originalAmount, profit math.Int)
	AfterHouseLoss(ctx sdk.Context, house sdk.AccAddress, originalAmount, lostAmt math.Int)
	AfterHouseRefund(ctx sdk.Context, house sdk.AccAddress, originalAmount math.Int)
	AfterHouseFeeRefund(ctx sdk.Context, house sdk.AccAddress, fee math.Int)
}

func (k *Keeper) RegisterHook(hook Hook) {
	k.hooks = append(k.hooks, hook)
}
