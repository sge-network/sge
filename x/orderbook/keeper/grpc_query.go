package keeper

import (
	"context"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sge-network/sge/x/orderbook/types"
)

var _ types.QueryServer = Keeper{}

// OrderBooks queries all order books that match the given status
func (k Keeper) OrderBooks(c context.Context, req *types.QueryOrderBooksRequest) (*types.QueryOrderBooksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// validate the provided status, return all the orderbooks if the status is empty
	if req.Status != "" && !(req.Status == types.OrderBookStatus_STATUS_ACTIVE.String() || req.Status == types.OrderBookStatus_STATUS_SETTLED.String()) {
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

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	orderBook, found := k.GetBook(ctx, req.BookId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "order book %s not found", req.BookId)
	}

	return &types.QueryOrderBookResponse{Orderbook: orderBook}, nil
}

// BookParticipants queries participants info for given orderbook
func (k Keeper) BookParticipants(c context.Context, req *types.QueryBookParticipantsRequest) (*types.QueryBookParticipantsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var bookParticipants []types.BookParticipant
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookParticipantKeyPrefix)
	participantStore := prefix.NewStore(store, types.GetBookParticipantsKey(req.BookId))
	pageRes, err := query.FilteredPaginate(participantStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		bookParticipant, err := types.UnmarshalBookParticipant(k.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			bookParticipants = append(bookParticipants, bookParticipant)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBookParticipantsResponse{
		BookParticipants: bookParticipants,
		Pagination:       pageRes,
	}, nil
}

// BookParticipant queries book participant info for given order book id and participant number
func (k Keeper) BookParticipant(c context.Context, req *types.QueryBookParticipantRequest) (*types.QueryBookParticipantResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.ParticipantNumber < 1 {
		return nil, status.Error(codes.InvalidArgument, "participant number can not be less than 1")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bookParticipant, found := k.GetBookParticipant(ctx, req.BookId, req.ParticipantNumber)
	if !found {
		return nil, status.Errorf(codes.NotFound, "book participant %s, %d not found", req.BookId, req.ParticipantNumber)
	}

	return &types.QueryBookParticipantResponse{BookParticipant: bookParticipant}, nil
}

// BookExposures queries exposures info for given orderbook
func (k Keeper) BookExposures(c context.Context, req *types.QueryBookExposuresRequest) (*types.QueryBookExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var bookExposures []types.BookOddExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BookOddExposureKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetBookOddExposuresKey(req.BookId))
	pageRes, err := query.FilteredPaginate(exposureStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		bookExposure, err := types.UnmarshalBookOddExposure(k.cdc, value)
		if err != nil {
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

// BookExposure queries book exposure info for given order book id and odd id
func (k Keeper) BookExposure(c context.Context, req *types.QueryBookExposureRequest) (*types.QueryBookExposureResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.OddId == "" {
		return nil, status.Error(codes.InvalidArgument, "odd id can not be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	bookExposure, found := k.GetBookOddExposure(ctx, req.BookId, req.OddId)
	if !found {
		return nil, status.Errorf(codes.NotFound, "book exposure %s, %s not found", req.BookId, req.OddId)
	}

	return &types.QueryBookExposureResponse{BookExposure: bookExposure}, nil
}

// ParticipantExposures queries participant exposures info for given orderbook
func (k Keeper) ParticipantExposures(c context.Context, req *types.QueryParticipantExposuresRequest) (*types.QueryParticipantExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var participantExposures []types.ParticipantExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipantExposureKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetParticipantExposuresByBookKey(req.BookId))
	pageRes, err := query.FilteredPaginate(exposureStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		participantExposure, err := types.UnmarshalParticipantExposure(k.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			participantExposures = append(participantExposures, participantExposure)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParticipantExposuresResponse{
		ParticipantExposures: participantExposures,
		Pagination:           pageRes,
	}, nil
}

// ParticipantExposure queries participant exposure info for given order book id and participant number
func (k Keeper) ParticipantExposure(c context.Context, req *types.QueryParticipantExposureRequest) (*types.QueryParticipantExposureResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.ParticipantNumber < 1 {
		return nil, status.Error(codes.InvalidArgument, "participant number can not be less than 1")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var participantExposures []types.ParticipantExposure

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipantExposureByPNKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetParticipantExposuresByPNKey(req.BookId, req.ParticipantNumber))
	pageRes, err := query.FilteredPaginate(exposureStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		participantExposure, err := types.UnmarshalParticipantExposure(k.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			participantExposures = append(participantExposures, participantExposure)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParticipantExposureResponse{
		ParticipantExposure: participantExposures,
		Pagination:          pageRes,
	}, nil
}

// HistoricalParticipantExposures queries historical participant exposures info for given orderbook
func (k Keeper) HistoricalParticipantExposures(c context.Context, req *types.QueryHistoricalParticipantExposuresRequest) (*types.QueryHistoricalParticipantExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var historicalParticipantExposures []types.ParticipantExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HistoricalParticipantExposureKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetParticipantExposuresByBookKey(req.BookId))
	pageRes, err := query.FilteredPaginate(exposureStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		participantExposure, err := types.UnmarshalParticipantExposure(k.cdc, value)
		if err != nil {
			return false, err
		}

		if accumulate {
			historicalParticipantExposures = append(historicalParticipantExposures, participantExposure)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryHistoricalParticipantExposuresResponse{
		ParticipantExposures: historicalParticipantExposures,
		Pagination:           pageRes,
	}, nil
}
