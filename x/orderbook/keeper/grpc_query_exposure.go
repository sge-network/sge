package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/orderbook/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderBookExposure queries book exposure info for given order book id and odds id
func (k Keeper) OrderBookExposure(
	c context.Context,
	req *types.QueryOrderBookExposureRequest,
) (*types.QueryOrderBookExposureResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "order book id can not be empty")
	}

	if req.OddsUid == "" {
		return nil, status.Error(codes.InvalidArgument, "odds id can not be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	orderBookExposure, found := k.GetOrderBookOddsExposure(ctx, req.OrderBookUid, req.OddsUid)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"order book exposure %s, %s not found",
			req.OrderBookUid,
			req.OddsUid,
		)
	}

	return &types.QueryOrderBookExposureResponse{OrderBookExposure: orderBookExposure}, nil
}

// OrderBookExposures queries exposures info for given orderbook
func (k Keeper) OrderBookExposures(
	c context.Context,
	req *types.QueryOrderBookExposuresRequest,
) (*types.QueryOrderBookExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "order book id cannot be empty")
	}
	var orderBookExposures []types.OrderBookOddsExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OrderBookOddsExposureKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetOrderBookOddsExposuresKey(req.OrderBookUid))
	pageRes, err := query.FilteredPaginate(
		exposureStore,
		req.Pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			var orderBookExposure types.OrderBookOddsExposure
			if err := k.cdc.Unmarshal(value, &orderBookExposure); err != nil {
				return false, err
			}

			if accumulate {
				orderBookExposures = append(orderBookExposures, orderBookExposure)
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOrderBookExposuresResponse{
		OrderBookExposures: orderBookExposures,
		Pagination:         pageRes,
	}, nil
}

// OrderBookParticipationExposures queries participation exposures info for a given orderbook
func (k Keeper) OrderBookParticipationExposures(
	c context.Context,
	req *types.QueryOrderBookParticipationExposuresRequest,
) (*types.QueryOrderBookParticipationExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var participationExposures []types.ParticipationExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipationExposureKeyPrefix)
	exposureStore := prefix.NewStore(
		store,
		types.GetParticipationExposuresByOrderBookKey(req.OrderBookUid),
	)
	pageRes, err := query.FilteredPaginate(
		exposureStore,
		req.Pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			var participationExposure types.ParticipationExposure
			if err := k.cdc.Unmarshal(value, &participationExposure); err != nil {
				return false, err
			}

			if accumulate {
				participationExposures = append(participationExposures, participationExposure)
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOrderBookParticipationExposuresResponse{
		ParticipationExposures: participationExposures,
		Pagination:             pageRes,
	}, nil
}

// ParticipationExposures queries participation exposure info for given order book id and participation index
func (k Keeper) ParticipationExposures(
	c context.Context,
	req *types.QueryParticipationExposuresRequest,
) (*types.QueryParticipationExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.ParticipationIndex < 1 {
		return nil, status.Error(codes.InvalidArgument, "participation index can not be less than 1")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var participationExposures []types.ParticipationExposure

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipationExposureByIndexKeyPrefix)
	exposureStore := prefix.NewStore(
		store,
		types.GetParticipationByIndexKey(req.OrderBookUid, req.ParticipationIndex),
	)
	pageRes, err := query.FilteredPaginate(
		exposureStore,
		req.Pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			var participationExposure types.ParticipationExposure
			if err := k.cdc.Unmarshal(value, &participationExposure); err != nil {
				return false, err
			}

			if accumulate {
				participationExposures = append(participationExposures, participationExposure)
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParticipationExposuresResponse{
		ParticipationExposures: participationExposures,
		Pagination:             pageRes,
	}, nil
}

// HistoricalParticipationExposures queries historical participation exposures info for given orderbook
func (k Keeper) HistoricalParticipationExposures(
	c context.Context, req *types.QueryHistoricalParticipationExposuresRequest,
) (*types.QueryHistoricalParticipationExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var historicalParticipationExposures []types.ParticipationExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HistoricalParticipationExposureKeyPrefix)
	exposureStore := prefix.NewStore(
		store,
		types.GetParticipationExposuresByOrderBookKey(req.OrderBookUid),
	)
	pageRes, err := query.FilteredPaginate(
		exposureStore,
		req.Pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			var participationExposure types.ParticipationExposure
			if err := k.cdc.Unmarshal(value, &participationExposure); err != nil {
				return false, err
			}

			if accumulate {
				historicalParticipationExposures = append(
					historicalParticipationExposures,
					participationExposure,
				)
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryHistoricalParticipationExposuresResponse{
		ParticipationExposures: historicalParticipationExposures,
		Pagination:             pageRes,
	}, nil
}
