package keeper

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	orderbookmodulekeeper "github.com/sge-network/sge/x/orderbook/keeper"
)

var _ orderbookmodulekeeper.Hook = Keeper{}

func (k Keeper) AfterBettorWin(ctx sdk.Context, bettor sdk.AccAddress, originalAmount sdk.Int, profit sdk.Int) {
	balance, exists := k.GetBalance(ctx, bettor)
	if !exists {
		return
	}
	err := balance.Unspend(originalAmount)
	if err != nil {
		panic(err)
	}
	// send profits to subaccount owner
	owner, exists := k.GetSubAccountOwner(ctx, bettor)
	if !exists {
		panic("subaccount owner not found")
	}
	err = k.bankKeeper.SendCoins(ctx, bettor, owner, sdk.NewCoins(sdk.NewCoin(k.GetParams(ctx).LockedBalanceDenom, profit)))
	if err != nil {
		panic(err)
	}
	k.SetBalance(ctx, bettor, balance)
}

func (k Keeper) AfterBettorLoss(ctx sdk.Context, bettor sdk.AccAddress, originalAmount sdk.Int) {
	balance, exists := k.GetBalance(ctx, bettor)
	if !exists {
		return
	}
	err := balance.Unspend(originalAmount)
	if err != nil {
		panic(err)
	}
	err = balance.AddLoss(originalAmount)
	if err != nil {
		panic(err)
	}
	k.SetBalance(ctx, bettor, balance)
}

func (k Keeper) AfterBettorRefund(ctx sdk.Context, bettor sdk.AccAddress, originalAmount, fee sdk.Int) {
	balance, exists := k.GetBalance(ctx, bettor)
	if !exists {
		return
	}
	totalUnspent := originalAmount.Add(fee)
	err := balance.Unspend(totalUnspent)
	if err != nil {
		panic(err)
	}
	k.SetBalance(ctx, bettor, balance)
	log.Printf("bettor refunded, yay!")
}
