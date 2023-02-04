package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OrderBookKeeper defines the expected orderbook keeper.
type OrderBookKeeper interface {
	AddBookParticipant(ctx sdk.Context, addr sdk.AccAddress, bookId string, liquidity, fee sdk.Int, feeAccountName string) (uint64, error)
}
