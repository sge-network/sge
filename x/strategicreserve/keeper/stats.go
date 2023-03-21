package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// SetOrderBookStats sets bet statistics in the store
func (k Keeper) SetOrderBookStats(ctx sdk.Context, stats types.OrderBookStats) {
	store := k.getBookStatsStore(ctx)
	b := k.cdc.MustMarshal(&stats)
	store.Set(utils.StrBytes("0"), b)
}

// GetOrderBookStats returns order-book stats
func (k Keeper) GetOrderBookStats(ctx sdk.Context) (val types.OrderBookStats) {
	store := k.getBookStatsStore(ctx)
	b := store.Get(utils.StrBytes("0"))
	if b == nil {
		return val
	}

	k.cdc.MustUnmarshal(b, &val)
	return val
}

func (k Keeper) AddBookSettlement(ctx sdk.Context, orderBookUID string) error {
	book, found := k.GetBook(ctx, orderBookUID)
	if !found {
		return sdkerrors.Wrapf(types.ErrOrderBookNotFound, "%s", orderBookUID)
	}
	if book.Status != types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE {
		return sdkerrors.Wrapf(types.ErrOrderBookNotActive, "%s", orderBookUID)
	}
	book.Status = types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_RESOLVED
	k.SetBook(ctx, book)

	stats := k.GetOrderBookStats(ctx)
	stats.ResolvedUnsettled = append(stats.ResolvedUnsettled, orderBookUID)
	k.SetOrderBookStats(ctx, stats)
	return nil
}

// GetFirstUnsettledResolvedOrderBook returns first element of resolved strategicreserve that have active deposits
func (k Keeper) GetFirstUnsettledResolvedOrderBook(ctx sdk.Context) (string, bool) {
	stats := k.GetOrderBookStats(ctx)
	if len(stats.ResolvedUnsettled) > 0 {
		return stats.ResolvedUnsettled[0], true
	}
	return "", false
}

// RemoveUnsettledResolvedOrderBook removes resolved order-book from the statistics
func (k Keeper) RemoveUnsettledResolvedOrderBook(ctx sdk.Context, orderBookUID string) {
	stats := k.GetOrderBookStats(ctx)
	if len(stats.ResolvedUnsettled) > 0 {
		for i, e := range stats.ResolvedUnsettled {
			if e == orderBookUID {
				stats.ResolvedUnsettled = append(stats.ResolvedUnsettled[:i], stats.ResolvedUnsettled[i+1:]...)
			}
		}
	}
	k.SetOrderBookStats(ctx, stats)
}
