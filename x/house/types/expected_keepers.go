package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OrderBookKeeper defines the expected orderbook keeper.
type OrderBookKeeper interface {
	AddBookParticipant(ctx sdk.Context, addr sdk.AccAddress, bookID string, liquidity, fee sdk.Int, feeAccountName string) (uint64, error)
	LiquidateBookParticipant(ctx sdk.Context, depAddr, bookID string, bpNumber uint64, mode WithdrawalMode, amount sdk.Int) (sdk.Int, error)
}
