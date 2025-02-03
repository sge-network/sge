package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/orderbook/types"
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
		func(_, value []byte, accumulate bool) (bool, error) {
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
