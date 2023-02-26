package keeper

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/x/orderbook/types"
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
	bookStore := prefix.NewStore(store, types.BookKeyPrefix)

	pageRes, err := query.FilteredPaginate(bookStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		orderBook, err := types.UnmarshalBook(k.cdc, value)
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

// OrderBook queries orderbook info for given order book id
func (k Keeper) OrderBook(c context.Context, req *types.QueryOrderBookRequest) (*types.QueryOrderBookResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	orderBook, found := k.GetBook(ctx, req.BookUid)
	if !found {
		return nil, status.Errorf(codes.NotFound, "order book %s not found", req.BookUid)
	}

	return &types.QueryOrderBookResponse{Orderbook: orderBook}, nil
}

// BookParticipations queries participations info for given orderbook
func (k Keeper) BookParticipations(c context.Context, req *types.QueryBookParticipationsRequest) (*types.QueryBookParticipationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var bookParticipations []types.BookParticipation
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookParticipationKeyPrefix)
	participationStore := prefix.NewStore(store, types.GetBookParticipationsKey(req.BookUid))
	pageRes, err := query.FilteredPaginate(participationStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var bookParticipation types.BookParticipation
		if err := k.cdc.Unmarshal(value, &bookParticipation); err != nil {
			return false, err
		}

		if accumulate {
			bookParticipations = append(bookParticipations, bookParticipation)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBookParticipationsResponse{
		BookParticipations: bookParticipations,
		Pagination:         pageRes,
	}, nil
}

// BookParticipation queries book participation info for given order book id and participation number
func (k Keeper) BookParticipation(c context.Context, req *types.QueryBookParticipationRequest) (*types.QueryBookParticipationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.ParticipationIndex < 1 {
		return nil, status.Error(codes.InvalidArgument, "participation number can not be less than 1")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bookParticipation, found := k.GetBookParticipation(ctx, req.BookUid, req.ParticipationIndex)
	if !found {
		return nil, status.Errorf(codes.NotFound, "book participation %s, %d not found", req.BookUid, req.ParticipationIndex)
	}

	return &types.QueryBookParticipationResponse{BookParticipation: bookParticipation}, nil
}

// BookExposures queries exposures info for given orderbook
func (k Keeper) BookExposures(c context.Context, req *types.QueryBookExposuresRequest) (*types.QueryBookExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var bookExposures []types.BookOddsExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookOddsExposureKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetBookOddsExposuresKey(req.BookUid))
	pageRes, err := query.FilteredPaginate(exposureStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var bookExposure types.BookOddsExposure
		if err := k.cdc.Unmarshal(value, &bookExposure); err != nil {
			return false, err
		}

		if accumulate {
			bookExposures = append(bookExposures, bookExposure)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBookExposuresResponse{
		BookExposures: bookExposures,
		Pagination:    pageRes,
	}, nil
}

// BookExposure queries book exposure info for given order book id and odds id
func (k Keeper) BookExposure(c context.Context, req *types.QueryBookExposureRequest) (*types.QueryBookExposureResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.OddsUid == "" {
		return nil, status.Error(codes.InvalidArgument, "odds id can not be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bookExposure, found := k.GetBookOddsExposure(ctx, req.BookUid, req.OddsUid)
	if !found {
		return nil, status.Errorf(codes.NotFound, "book exposure %s, %s not found", req.BookUid, req.OddsUid)
	}

	return &types.QueryBookExposureResponse{BookExposure: bookExposure}, nil
}
