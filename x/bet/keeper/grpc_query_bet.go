package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/bet/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Bets returns all bets
func (k Keeper) Bets(c context.Context, req *types.QueryBetsRequest) (*types.QueryBetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var bets []types.Bet
	ctx := sdk.UnwrapSDKContext(c)

	betStore := k.getBetStore(ctx)

	pageRes, err := query.Paginate(betStore, req.Pagination, func(key []byte, value []byte) error {
		var bet types.Bet
		if err := k.cdc.Unmarshal(value, &bet); err != nil {
			return err
		}

		bets = append(bets, bet)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryBetsResponse{Bet: bets, Pagination: pageRes}, nil
}

// Bets returns bets with selected uids
func (k Keeper) BetsByUIDs(c context.Context, req *types.QueryBetsByUIDsRequest) (*types.QueryBetsByUIDsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	ctx := sdk.UnwrapSDKContext(c)

	count := len(req.GetUids())
	if count > types.MaxQueriableBetsCount {
		return nil, types.ErrCanNotQueryLargeNumberOfBets
	}

	req.Uids = utils.RemoveDuplicateStrs(req.Uids)

	foundBets := make([]types.Bet, 0, count)
	notFoundBets := make([]string, 0)
	for _, id := range req.GetUids() {
		val, found := k.GetBet(ctx, id)
		if !found {
			notFoundBets = append(notFoundBets, id)
			continue
		}
		foundBets = append(foundBets, val)
	}

	return &types.QueryBetsByUIDsResponse{
		Bets:           foundBets,
		NotFoundEvents: notFoundBets,
	}, nil

}

// Bet returns a specific bet by its UID
func (k Keeper) Bet(c context.Context, req *types.QueryBetRequest) (*types.QueryBetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetBet(
		ctx,
		req.Uid,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryBetResponse{Bet: val}, nil
}
