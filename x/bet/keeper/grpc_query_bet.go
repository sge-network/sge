package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
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

// Bets returns all bets of certain creator sort-able by pagination attributes
func (k Keeper) BetsByCreator(c context.Context, req *types.QueryBetsByCreatorRequest) (*types.QueryBetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var bets []types.Bet
	ctx := sdk.UnwrapSDKContext(c)

	betStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BetListByCreatorKey(req.Creator))

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

	count := len(req.GetItems())
	if count > types.MaxAllowedQueryBetsCount {
		return nil, types.ErrCanNotQueryLargeNumberOfBets
	}

	req.Items = RemoveDuplicateUIDs(req.Items)

	foundBets := make([]types.Bet, 0, count)
	notFoundBets := make([]string, 0)
	for _, item := range req.GetItems() {
		uid2ID, found := k.GetBetID(ctx, item.Uid)
		if !found {
			notFoundBets = append(notFoundBets, item.Uid)
			continue
		}

		val, found := k.GetBet(ctx, item.Creator, uid2ID.ID)
		if !found {
			notFoundBets = append(notFoundBets, item.Uid)
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

	uid2ID, found := k.GetBetID(ctx, req.Uid)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	val, found := k.GetBet(
		ctx,
		req.Creator,
		uid2ID.ID,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryBetResponse{Bet: val}, nil
}

// RemoveDuplicateUIDs retuns input array without duplicates
func RemoveDuplicateUIDs(strSlice []*types.QueryBetRequest) (list []*types.QueryBetRequest) {
	keys := make(map[string]bool)

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range strSlice {
		if _, value := keys[entry.Uid]; !value {
			keys[entry.Uid] = true
			list = append(list, entry)
		}
	}
	return
}
