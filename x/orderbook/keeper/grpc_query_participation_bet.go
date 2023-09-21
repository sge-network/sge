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

// ParticipationFulfilledBets queries participation fulfilled bets info for given order book id and participation index
func (k Keeper) ParticipationFulfilledBets(
	c context.Context,
	req *types.QueryParticipationFulfilledBetsRequest,
) (*types.QueryParticipationFulfilledBetsResponse, error) {
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
	var participationBets []types.ParticipationBetPair

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ParticipationBetPairKeyPrefix)
	betsStore := prefix.NewStore(
		store,
		types.GetParticipationByIndexKey(req.OrderBookUid, req.ParticipationIndex),
	)
	pageRes, err := query.FilteredPaginate(
		betsStore,
		req.Pagination,
		func(key, value []byte, accumulate bool) (bool, error) {
			var participationBet types.ParticipationBetPair
			if err := k.cdc.Unmarshal(value, &participationBet); err != nil {
				return false, err
			}

			if accumulate {
				participationBets = append(participationBets, participationBet)
			}
			return true, nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParticipationFulfilledBetsResponse{
		ParticipationBets: participationBets,
		Pagination:        pageRes,
	}, nil
}
