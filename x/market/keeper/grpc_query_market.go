package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/market/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Markets returns all the markets
func (k Keeper) Markets(c context.Context, req *types.QueryMarketsRequest) (*types.QueryMarketsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var markets []types.Market
	ctx := sdk.UnwrapSDKContext(c)

	marketStore := k.getMarketsStore(ctx)

	pageRes, err := query.Paginate(marketStore, req.Pagination, func(key []byte, value []byte) error {
		var market types.Market
		if err := k.cdc.Unmarshal(value, &market); err != nil {
			return err
		}

		markets = append(markets, market)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMarketsResponse{Market: markets, Pagination: pageRes}, nil
}

// Market returns a specific markets by its UID
func (k Keeper) Market(c context.Context, req *types.QueryMarketRequest) (*types.QueryMarketResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetMarket(
		ctx,
		req.Uid,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryMarketResponse{Market: val}, nil
}

// MarketsByUIDs return success markets and failed markets by uids back to the caller
func (k Keeper) MarketsByUIDs(goCtx context.Context, req *types.QueryMarketsByUIDsRequest) (*types.QueryMarketsByUIDsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	success := make([]types.Market, 0, len(req.Uids))
	failed := make([]string, 0)
	for _, id := range req.GetUids() {
		val, found := k.GetMarket(ctx, id)
		if !found {
			failed = append(failed, id)
			continue
		}
		success = append(success, val)
	}

	return &types.QueryMarketsByUIDsResponse{
		Markets:       success,
		FailedMarkets: failed,
	}, nil
}
