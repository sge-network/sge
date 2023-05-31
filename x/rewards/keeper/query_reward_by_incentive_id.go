package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/x/rewards/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RewardByIncentiveId(goCtx context.Context, req *types.QueryRewardByIncentiveIdRequest) (*types.QueryRewardByIncentiveIdResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var rewards []*types.RewardK

	store := prefix.NewStore(k.getRewardStore(ctx), types.GetRewardKey(req.IncentiveId))
	_, err := query.Paginate(store, nil, func(key []byte, value []byte) error {
		var reward types.RewardK
		if err := k.cdc.Unmarshal(value, &reward); err != nil {
			return err
		}

		rewards = append(rewards, &reward)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRewardByIncentiveIdResponse{Rewards: rewards}, nil
}
