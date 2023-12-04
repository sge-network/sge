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

// OrderBookParticipation queries book participation info for given order book id and participation index
func (k Keeper) OrderBookParticipation(
	c context.Context,
	req *types.QueryOrderBookParticipationRequest,
) (*types.QueryOrderBookParticipationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "order book id can not be empty")
	}

	if req.ParticipationIndex < 1 {
		return nil, status.Error(codes.InvalidArgument, "participation index can not be less than 1")
	}

	ctx := sdk.UnwrapSDKContext(c)
	orderBookParticipation, found := k.GetOrderBookParticipation(
		ctx,
		req.OrderBookUid,
		req.ParticipationIndex,
	)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"order book participation %s, %d not found",
			req.OrderBookUid,
			req.ParticipationIndex,
		)
	}

	return &types.QueryOrderBookParticipationResponse{
		OrderBookParticipation: orderBookParticipation,
	}, nil
}

// OrderBookParticipations queries participation info for a given orderbook
func (k Keeper) OrderBookParticipations(
	c context.Context,
	req *types.QueryOrderBookParticipationsRequest,
) (*types.QueryOrderBookParticipationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	if req.OrderBookUid == "" {
		return nil, status.Error(codes.InvalidArgument, "order book id cannot be empty")
	}
	var orderBookParticipations []types.OrderBookParticipation
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OrderBookParticipationKeyPrefix)
	participationStore := prefix.NewStore(store, types.GetOrderBookParticipationsKey(req.OrderBookUid))
	pageRes, err := query.FilteredPaginate(
		participationStore,
		req.Pagination,
		func(key []byte, value []byte, accumulate bool) (bool, error) {
			var bookParticipation types.OrderBookParticipation
			if err := k.cdc.Unmarshal(value, &bookParticipation); err != nil {
				return false, err
			}

			if accumulate {
				orderBookParticipations = append(orderBookParticipations, bookParticipation)
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOrderBookParticipationsResponse{
		OrderBookParticipations: orderBookParticipations,
		Pagination:              pageRes,
	}, nil
}

// SettledOrderbookParticipationsOfHeight returns settled orderbook participations of a certain height
func (k Keeper) SettledOrderbookParticipationsOfHeight(
	c context.Context,
	req *types.QuerySettledOrderBookParticipationsOfHeightRequest,
) (*types.QuerySettledOrderBookParticipationsOfHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var orderbookParticipations []types.OrderBookParticipation
	ctx := sdk.UnwrapSDKContext(c)

	settledOrderbookParticipationStore := prefix.NewStore(
		ctx.KVStore(k.storeKey),
		types.SettledOrderbookParticipationListOfBlockHeightPrefix(req.BlockHeight),
	)

	pageRes, err := query.Paginate(
		settledOrderbookParticipationStore,
		req.Pagination,
		func(key []byte, value []byte) error {
			var settledOrderbookParticipation types.SettledOrderbookParticipation
			if err := k.cdc.Unmarshal(value, &settledOrderbookParticipation); err != nil {
				return err
			}

			orderBookParticipation, found := k.GetOrderBookParticipation(
				ctx,
				settledOrderbookParticipation.OrderBookUID,
				settledOrderbookParticipation.Index,
			)
			if found {
				orderbookParticipations = append(orderbookParticipations, orderBookParticipation)
			}

			return nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySettledOrderBookParticipationsOfHeightResponse{Participations: orderbookParticipations, Pagination: pageRes}, nil
}
