package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/orderbook/types"
)

// SetOrderBookStats sets bet statistics in the store
func (k Keeper) SetOrderBookStats(ctx sdk.Context, stats types.OrderBookStats) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookStatsKeyPrefix)
	b := k.cdc.MustMarshal(&stats)
	store.Set(utils.StrBytes("0"), b)
}

// GetOrderBookStats returns order-book stats
func (k Keeper) GetOrderBookStats(ctx sdk.Context) (val types.OrderBookStats) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookStatsKeyPrefix)
	b := store.Get(utils.StrBytes("0"))
	if b == nil {
		return val
	}

	k.cdc.MustUnmarshal(b, &val)
	return val
}

func (k Keeper) AddBookSettlement(ctx sdk.Context, orderBookID string) error {
	book, found := k.GetBook(ctx, orderBookID)
	if !found {
		return sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", orderBookID)
	}
	if book.Status != types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE {
		return sdkerrors.Wrapf(types.ErrOrderBookNotActive, "%s", orderBookID)
	}
	book.Status = types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_RESOLVED
	k.SetBook(ctx, book)

	stats := k.GetOrderBookStats(ctx)
	stats.ResolvedUnsettled = append(stats.ResolvedUnsettled, orderBookID)
	k.SetOrderBookStats(ctx, stats)
	return nil
}

// GetFirstUnsettledResolvedOrderBook returns first element of resolved orderbook that have active deposits
func (k Keeper) GetFirstUnsettledResolvedOrderBook(ctx sdk.Context) (string, bool) {
	stats := k.GetOrderBookStats(ctx)
	if len(stats.ResolvedUnsettled) > 0 {
		return stats.ResolvedUnsettled[0], true
	}
	return "", false
}

// RemoveUnsettledResolvedOrderBook removes resolved order-book from the statistics
func (k Keeper) RemoveUnsettledResolvedOrderBook(ctx sdk.Context, orderBookID string) {
	stats := k.GetOrderBookStats(ctx)
	if len(stats.ResolvedUnsettled) > 0 {
		for i, e := range stats.ResolvedUnsettled {
			if e == orderBookID {
				stats.ResolvedUnsettled = append(stats.ResolvedUnsettled[:i], stats.ResolvedUnsettled[i+1:]...)
			}
		}
	}
	k.SetOrderBookStats(ctx, stats)
}