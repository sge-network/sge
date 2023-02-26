package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OrderBookKeeper defines the expected orderbook keeper.
type OrderBookKeeper interface {
	InitiateBookParticipation(ctx sdk.Context, addr sdk.AccAddress, bookUID string, liquidity, fee sdk.Int) (uint64, error)
	LiquidateBookParticipation(ctx sdk.Context, depAddr, bookUID string, bpNumber uint64, mode WithdrawalMode, amount sdk.Int) (sdk.Int, error)
}
