package keeper

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/bet/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	mapSeparator = ":"
)

// Bets returns all bets
func (k Keeper) Bets(c context.Context, req *types.QueryBetsRequest) (*types.QueryBetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
		// TODO: Registered errors?
	}

	var bets []types.Bet
	ctx := sdk.UnwrapSDKContext(c)

	betStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BetListPrefix)

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

// BetsByCreator returns all bets of certain creator sort-able by pagination attributes
func (k Keeper) BetsByCreator(c context.Context, req *types.QueryBetsByCreatorRequest) (*types.QueryBetsByCreatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var bets []types.Bet
	ctx := sdk.UnwrapSDKContext(c)

	betStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.BetListByCreatorPrefix(req.Creator))

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

	return &types.QueryBetsByCreatorResponse{Bet: bets, Pagination: pageRes}, nil
}

// BetsByUIDs returns bets with selected uids
func (k Keeper) BetsByUIDs(c context.Context, req *types.QueryBetsByUIDsRequest) (*types.QueryBetsByUIDsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	ctx := sdk.UnwrapSDKContext(c)

	count := len(req.GetItems())
	if count > types.MaxAllowedQueryBetsCount {
		return nil, types.ErrCanNotQueryLargeNumberOfBets
	}

	var items []*types.QueryBetRequest
	for _, val := range req.Items {
		pair := strings.Split(val, mapSeparator)
		if len(pair) != 2 {
			return nil, fmt.Errorf("each pair should be separated by colon ex. creator:uid")
		}

		items = append(items, &types.QueryBetRequest{
			Creator: pair[0],
			Uid:     pair[1],
		})
	}

	items = removeDuplicateUIDs(items)

	foundBets := make([]types.Bet, 0, count)
	notFoundBets := make([]string, 0)
	for _, item := range items {
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
		Bets:         foundBets,
		UidsNotFound: notFoundBets,
		// TODO: NotFoundBets
	}, nil
}

// ActiveBets returns active bets of a market
func (k Keeper) ActiveBets(c context.Context, req *types.QueryActiveBetsRequest) (*types.QueryActiveBetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var bets []types.Bet
	ctx := sdk.UnwrapSDKContext(c)

	activeBetStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActiveBetListOfMarketPrefix(req.MarketUid))

	pageRes, err := query.Paginate(activeBetStore, req.Pagination, func(key []byte, value []byte) error {
		var activeBet types.ActiveBet
		if err := k.cdc.Unmarshal(value, &activeBet); err != nil {
			return err
		}

		uid2ID, found := k.GetBetID(ctx, activeBet.UID)
		if found {
			bet, found := k.GetBet(ctx, activeBet.Creator, uid2ID.ID)
			if found {
				bets = append(bets, bet)
			}
		}

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryActiveBetsResponse{Bet: bets, Pagination: pageRes}, nil
}

// SettledBetsOfHeight returns settled bets of a certain height
func (k Keeper) SettledBetsOfHeight(c context.Context, req *types.QuerySettledBetsOfHeightRequest) (*types.QuerySettledBetsOfHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var bets []types.Bet
	ctx := sdk.UnwrapSDKContext(c)

	settledBetStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.SettledBetListOfBlockHeightPrefix(req.BlockHeight))

	pageRes, err := query.Paginate(settledBetStore, req.Pagination, func(key []byte, value []byte) error {
		var settledBet types.SettledBet
		if err := k.cdc.Unmarshal(value, &settledBet); err != nil {
			return err
		}

		uid2ID, found := k.GetBetID(ctx, settledBet.UID)
		if found {
			bet, found := k.GetBet(ctx, settledBet.BettorAddress, uid2ID.ID)
			if found {
				bets = append(bets, bet)
			}
		}

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySettledBetsOfHeightResponse{Bet: bets, Pagination: pageRes}, nil
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

	market, found := k.marketKeeper.GetMarket(ctx, val.MarketUID)
	if !found {
		return nil, status.Errorf(codes.NotFound, "corresponding market with id %s not found", val.MarketUID)
	}

	return &types.QueryBetResponse{Bet: val, Market: market}, nil
}

// removeDuplicateUIDs returns input array without duplicates
func removeDuplicateUIDs(strSlice []*types.QueryBetRequest) (list []*types.QueryBetRequest) {
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
