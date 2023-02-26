package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/x/orderbook/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ParticipationExposures queries participation exposures info for given orderbook
func (k Keeper) ParticipationExposures(c context.Context, req *types.QueryParticipationExposuresRequest) (*types.QueryParticipationExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var participationExposures []types.ParticipationExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipationExposureKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetParticipationExposuresByBookKey(req.BookId))
	pageRes, err := query.FilteredPaginate(exposureStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var participationExposure types.ParticipationExposure
		if err := k.cdc.Unmarshal(value, &participationExposure); err != nil {
			return false, err
		}

		if accumulate {
			participationExposures = append(participationExposures, participationExposure)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParticipationExposuresResponse{
		ParticipationExposures: participationExposures,
		Pagination:             pageRes,
	}, nil
}

// ParticipationExposure queries participation exposure info for given order book id and participation number
func (k Keeper) ParticipationExposure(c context.Context, req *types.QueryParticipationExposureRequest) (*types.QueryParticipationExposureResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.ParticipationIndex < 1 {
		return nil, status.Error(codes.InvalidArgument, "participation number can not be less than 1")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var participationExposures []types.ParticipationExposure

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipationExposureByIndexKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetParticipationByIndexKey(req.BookId, req.ParticipationIndex))
	pageRes, err := query.FilteredPaginate(exposureStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var participationExposure types.ParticipationExposure
		if err := k.cdc.Unmarshal(value, &participationExposure); err != nil {
			return false, err
		}

		if accumulate {
			participationExposures = append(participationExposures, participationExposure)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParticipationExposureResponse{
		ParticipationExposure: participationExposures,
		Pagination:            pageRes,
	}, nil
}

// HistoricalParticipationExposures queries historical participation exposures info for given orderbook
func (k Keeper) HistoricalParticipationExposures(c context.Context, req *types.QueryHistoricalParticipationExposuresRequest) (*types.QueryHistoricalParticipationExposuresResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id cannot be empty")
	}
	var historicalParticipationExposures []types.ParticipationExposure
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HistoricalParticipationExposureKeyPrefix)
	exposureStore := prefix.NewStore(store, types.GetParticipationExposuresByBookKey(req.BookId))
	pageRes, err := query.FilteredPaginate(exposureStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var participationExposure types.ParticipationExposure
		if err := k.cdc.Unmarshal(value, &participationExposure); err != nil {
			return false, err
		}

		if accumulate {
			historicalParticipationExposures = append(historicalParticipationExposures, participationExposure)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryHistoricalParticipationExposuresResponse{
		ParticipationExposures: historicalParticipationExposures,
		Pagination:             pageRes,
	}, nil
}

// ParticipationFulfilledBets queries participation fulfilled bets info for given order book id and participation number
func (k Keeper) ParticipationFulfilledBets(c context.Context, req *types.QueryParticipationFulfilledBetsRequest) (*types.QueryParticipationFulfilledBetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.BookId == "" {
		return nil, status.Error(codes.InvalidArgument, "book id can not be empty")
	}

	if req.ParticipationIndex < 1 {
		return nil, status.Error(codes.InvalidArgument, "participation number can not be less than 1")
	}

	ctx := sdk.UnwrapSDKContext(c)
	var participationBets []types.ParticipationBetPair

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipationBetPairKeyPrefix)
	betsStore := prefix.NewStore(store, types.GetParticipationByIndexKey(req.BookId, req.ParticipationIndex))
	pageRes, err := query.FilteredPaginate(betsStore, req.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var participationBet types.ParticipationBetPair
		if err := k.cdc.Unmarshal(value, &participationBet); err != nil {
			return false, err
		}

		if accumulate {
			participationBets = append(participationBets, participationBet)
		}
		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParticipationFulfilledBetsResponse{
		ParticipationBets: participationBets,
		Pagination:        pageRes,
	}, nil
}
