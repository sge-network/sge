package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ OrderBookHooks = MultiOrderBookHooks{}

// MultiOrderBookHooks combine multiple orderbook hooks, all hook functions are run in array sequence
type MultiOrderBookHooks []OrderBookHooks

// NewMultiOrderBookHooks returns list of hooks defined for the orderbook module.
func NewMultiOrderBookHooks(hooks ...OrderBookHooks) MultiOrderBookHooks {
	return hooks
}

// AfterHouseWin registers all of hooks for this method.
func (h MultiOrderBookHooks) AfterHouseWin(ctx sdk.Context, house sdk.AccAddress, originalAmount, profit sdkmath.Int) {
	for i := range h {
		h[i].AfterHouseWin(ctx, house, originalAmount, profit)
	}
}

// AfterHouseLoss registers all of hooks for this method.
func (h MultiOrderBookHooks) AfterHouseLoss(ctx sdk.Context, house sdk.AccAddress, originalAmount, lostAmt sdkmath.Int) {
	for i := range h {
		h[i].AfterHouseLoss(ctx, house, originalAmount, lostAmt)
	}
}

// AfterHouseRefund registers all of hooks for this method.
func (h MultiOrderBookHooks) AfterHouseRefund(ctx sdk.Context, house sdk.AccAddress, originalAmount sdkmath.Int) {
	for i := range h {
		h[i].AfterHouseRefund(ctx, house, originalAmount)
	}
}

// AfterHouseFeeRefund registers all of hooks for this method.
func (h MultiOrderBookHooks) AfterHouseFeeRefund(ctx sdk.Context, house sdk.AccAddress, fee sdkmath.Int) {
	for i := range h {
		h[i].AfterHouseRefund(ctx, house, fee)
	}
}
