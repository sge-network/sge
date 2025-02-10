package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/app/params"
	orderbooktypes "github.com/sge-network/sge/x/legacy/orderbook/types"
)

// Hooks wrapper struct for slashing keeper
type Hooks struct {
	k Keeper
}

// subaccount module shoule implement the orderbook module hooks.
var _ orderbooktypes.OrderBookHooks = Hooks{}

// Create new distribution hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// AfterHouseWin is subaccount module hook for house winning over subbacount.
func (h Hooks) AfterHouseWin(ctx sdk.Context, house sdk.AccAddress, originalAmount, profit sdkmath.Int) {
	// update balance
	balance, exists := h.k.GetAccountSummary(ctx, house)
	if !exists {
		return
	}

	err := balance.Unspend(originalAmount)
	if err != nil {
		panic(err)
	}
	h.k.SetAccountSummary(ctx, house, balance)

	// send profits
	subAccountOwner, exists := h.k.GetSubaccountOwner(ctx, house)
	if !exists {
		panic("data corruption: subaccount owner not found")
	}
	err = h.k.bankKeeper.SendCoins(ctx, house, subAccountOwner, sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, profit)))
	if err != nil {
		panic(err)
	}
}

// AfterHouseLoss is subaccount module hook for house loss for subbacount.
func (h Hooks) AfterHouseLoss(ctx sdk.Context, house sdk.AccAddress, originalAmount, lostAmt sdkmath.Int) {
	balance, exists := h.k.GetAccountSummary(ctx, house)
	if !exists {
		return
	}

	err := balance.Unspend(originalAmount)
	if err != nil {
		panic(err)
	}
	err = balance.AddLoss(lostAmt)
	if err != nil {
		panic(err)
	}

	h.k.SetAccountSummary(ctx, house, balance)
}

// AfterHouseRefund is subaccount module hook for house refund in subaccount deposit.
func (h Hooks) AfterHouseRefund(ctx sdk.Context, house sdk.AccAddress, originalAmount sdkmath.Int) {
	balance, exists := h.k.GetAccountSummary(ctx, house)
	if !exists {
		return
	}

	err := balance.Unspend(originalAmount)
	if err != nil {
		panic(err)
	}

	h.k.SetAccountSummary(ctx, house, balance)
}

// AfterHouseFeeRefund is subaccount module hook for house fee refund in subaccount deposit.
func (h Hooks) AfterHouseFeeRefund(ctx sdk.Context, house sdk.AccAddress, fee sdkmath.Int) {
	balance, exists := h.k.GetAccountSummary(ctx, house)
	if !exists {
		return
	}

	err := balance.Unspend(fee)
	if err != nil {
		panic(err)
	}

	h.k.SetAccountSummary(ctx, house, balance)
}
