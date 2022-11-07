package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/sportevent/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SportEventAll returns all the sport events
func (k Keeper) SportEventAll(c context.Context, req *types.QuerySportEventListAllRequest) (*types.QuerySportEventListAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var sportEvents []types.SportEvent
	ctx := sdk.UnwrapSDKContext(c)

	sportEventStore := k.getSportEventsStore(ctx)

	pageRes, err := query.Paginate(sportEventStore, req.Pagination, func(key []byte, value []byte) error {
		var sportEvent types.SportEvent
		if err := k.cdc.Unmarshal(value, &sportEvent); err != nil {
			return err
		}

		sportEvents = append(sportEvents, sportEvent)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QuerySportEventListAllResponse{SportEvent: sportEvents, Pagination: pageRes}, nil
}

// SportEvent returns a specific sport events by its UID
func (k Keeper) SportEvent(c context.Context, req *types.QuerySportEventRequest) (*types.QuerySportEventResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetSportEvent(
		ctx,
		req.Uid,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QuerySportEventResponse{SportEvent: val}, nil
}

// SportEventListByUIDs return success events and failed events id only back to the caller
func (k Keeper) SportEventListByUIDs(goCtx context.Context, req *types.QuerySportEventListByUIDsRequest) (*types.QuerySportEventListByUIDsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	success := make([]types.SportEvent, 0, len(req.Uids))
	failed := make([]string, 0)
	for _, id := range req.GetUids() {
		val, found := k.GetSportEvent(ctx, id)
		if !found {
			failed = append(failed, id)
			continue
		}
		success = append(success, val)
	}

	return &types.QuerySportEventListByUIDsResponse{
		SportEvents:  success,
		FailedEvents: failed,
	}, nil
}
