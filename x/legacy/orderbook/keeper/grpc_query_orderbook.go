package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/legacy/orderbook/types"
)

// OrderBook queries orderbook info for given order book id
func (k Keeper) OrderBook(
	c context.Context,
	req *types.QueryOrderBookRequest,
) (*types.QueryOrderBookResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "order book id can not be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	orderBook, found := k.GetOrderBook(ctx, req.OrderBookUid)
	if !found {
		return nil, status.Errorf(codes.NotFound, "order book %s not found", req.OrderBookUid)
	}

	return &types.QueryOrderBookResponse{OrderBook: orderBook}, nil
}

// OrderBooks queries all order books that match the given status
func (k Keeper) OrderBooks(
	c context.Context,
	req *types.QueryOrderBooksRequest,
) (*types.QueryOrderBooksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	// validate the provided status, return all the orderbooks if the status is empty
	if req.Status != "" &&
		!(req.Status == types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE.String() ||
			req.Status == types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_SETTLED.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order book status %s", req.Status)
	}

	var orderBooks []types.OrderBook
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	bookStore := prefix.NewStore(store, types.OrderBookKeyPrefix)

	pageRes, err := query.FilteredPaginate(
		bookStore,
		req.Pagination,
		func(_, value []byte, accumulate bool) (bool, error) {
			orderBook, err := types.UnmarshalOrderBook(k.cdc, value)
			if err != nil {
				return false, err
			}

			if req.Status != "" && !strings.EqualFold(orderBook.Status.String(), req.Status) {
				return false, nil
			}

			if accumulate {
				orderBooks = append(orderBooks, orderBook)
			}

			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOrderBooksResponse{Orderbooks: orderBooks, Pagination: pageRes}, nil
}
