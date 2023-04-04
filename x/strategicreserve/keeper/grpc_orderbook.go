package keeper

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/x/strategicreserve/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderBooks queries all order books that match the given status
func (k Keeper) OrderBooks(c context.Context, req *types.QueryOrderBooksRequest) (*types.QueryOrderBooksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// validate the provided status, return all the orderbooks if the status is empty
	if req.Status != "" && !(req.Status == types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_ACTIVE.String() ||
		req.Status == types.OrderBookStatus_ORDER_BOOK_STATUS_STATUS_SETTLED.String()) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order book status %s", req.Status)
	}

	var orderBooks []types.OrderBook
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	bookStore := prefix.NewStore(store, types.OrderBookKeyPrefix)

	pageRes, err := query.FilteredPaginate(bookStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
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
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOrderBooksResponse{Orderbooks: orderBooks, Pagination: pageRes}, nil
}

// OrderBook queries strategicreserve info for given order book id
func (k Keeper) OrderBook(c context.Context, req *types.QueryOrderBookRequest) (*types.QueryOrderBookResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	orderBook, found := k.GetOrderBook(ctx, req.OrderBookUid)
	if !found {
		return nil, status.Errorf(codes.NotFound, "order book %s not found", req.OrderBookUid)
	}

	return &types.QueryOrderBookResponse{OrderBook: orderBook}, nil
}

// OrderBookParticipations queries participations info for given strategicreserve
func (k Keeper) OrderBookParticipations(c context.Context, req *types.QueryOrderBookParticipationsRequest) (*types.QueryOrderBookParticipationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var orderBookParticipations []types.OrderBookParticipation
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OrderBookParticipationKeyPrefix)
	participationStore := prefix.NewStore(store, types.GetOrderBookParticipationsKey(req.OrderBookUid))
	pageRes, err := query.FilteredPaginate(participationStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var bookParticipation types.OrderBookParticipation
		if err := k.cdc.Unmarshal(value, &bookParticipation); err != nil {
			return false, err
		}

		if accumulate {
			orderBookParticipations = append(orderBookParticipations, bookParticipation)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOrderBookParticipationsResponse{
		OrderBookParticipations: orderBookParticipations,
		Pagination:              pageRes,
	}, nil
}

// OrderBookParticipation queries book participation info for given order book id and participation index
func (k Keeper) OrderBookParticipation(c context.Context, req *types.QueryOrderBookParticipationRequest) (*types.QueryOrderBookParticipationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.ParticipationIndex < 1 {
		return nil, status.Error(codes.InvalidArgument, "participation index can not be less than 1")
	}

	ctx := sdk.UnwrapSDKContext(c)
	orderBookParticipation, found := k.GetOrderBookParticipation(ctx, req.OrderBookUid, req.ParticipationIndex)
	if !found {
		return nil, status.Errorf(codes.NotFound, "book participation %s, %d not found", req.OrderBookUid, req.ParticipationIndex)
	}

	return &types.QueryOrderBookParticipationResponse{OrderBookParticipation: orderBookParticipation}, nil
}

// OrderBookExposures queries exposures info for given strategicreserve
func (k Keeper) OrderBookExposures(c context.Context, req *types.QueryOrderBookExposuresRequest) (*types.QueryOrderBookExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var orderBookExposures []types.OrderBookOddsExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OrderBookOddsExposureKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetOrderBookOddsExposuresKey(req.OrderBookUid))
	pageRes, err := query.FilteredPaginate(exposureStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var orderBookExposure types.OrderBookOddsExposure
		if err := k.cdc.Unmarshal(value, &orderBookExposure); err != nil {
			return false, err
		}

		if accumulate {
			orderBookExposures = append(orderBookExposures, orderBookExposure)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOrderBookExposuresResponse{
		OrderBookExposures: orderBookExposures,
		Pagination:         pageRes,
	}, nil
}

// OrderBookExposure queries book exposure info for given order book id and odds id
func (k Keeper) OrderBookExposure(c context.Context, req *types.QueryOrderBookExposureRequest) (*types.QueryOrderBookExposureResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.OddsUid == "" {
		return nil, status.Error(codes.InvalidArgument, "odds id can not be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	orderBookExposure, found := k.GetOrderBookOddsExposure(ctx, req.OrderBookUid, req.OddsUid)
	if !found {
		return nil, status.Errorf(codes.NotFound, "book exposure %s, %s not found", req.OrderBookUid, req.OddsUid)
	}

	return &types.QueryOrderBookExposureResponse{OrderBookExposure: orderBookExposure}, nil
}
