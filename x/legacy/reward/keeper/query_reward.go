package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/x/legacy/reward/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Rewards(goCtx context.Context, req *types.QueryRewardsRequest) (*types.QueryRewardsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var rewards []types.Reward
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	rewardStore := prefix.NewStore(store, types.RewardKeyPrefix)

	pageRes, err := query.Paginate(rewardStore, req.Pagination, func(_ []byte, value []byte) error {
		var reward types.Reward
		if err := k.cdc.Unmarshal(value, &reward); err != nil {
			return err
		}

		rewards = append(rewards, reward)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRewardsResponse{Rewards: rewards, Pagination: pageRes}, nil
}

func (k Keeper) Reward(goCtx context.Context, req *types.QueryRewardRequest) (*types.QueryRewardResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetReward(
		ctx,
		req.Uid,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryRewardResponse{Reward: val}, nil
}

func (k Keeper) RewardsByAddress(goCtx context.Context, req *types.QueryRewardsByAddressRequest) (*types.QueryRewardsByAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var rewards []types.RewardByCategory
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := k.getRewardByReceiverAndCategoryStore(ctx)
	rewardStore := prefix.NewStore(store, types.GetRewardsOfReceiverByPromoterPrefix(req.PromoterUid, req.Address))

	pageRes, err := query.Paginate(rewardStore, req.Pagination, func(_ []byte, value []byte) error {
		var reward types.RewardByCategory
		if err := k.cdc.Unmarshal(value, &reward); err != nil {
			return err
		}

		rewards = append(rewards, reward)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRewardsByAddressResponse{Rewards: rewards, Pagination: pageRes}, nil
}

func (k Keeper) RewardsByAddressAndCategory(goCtx context.Context, req *types.QueryRewardsByAddressAndCategoryRequest) (*types.QueryRewardsByAddressAndCategoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var rewards []types.RewardByCategory
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := k.getRewardByReceiverAndCategoryStore(ctx)
	rewardStore := prefix.NewStore(store, types.GetRewardsOfReceiverByPromoterAndCategoryPrefix(req.PromoterUid, req.Address, req.Category))

	pageRes, err := query.Paginate(rewardStore, req.Pagination, func(_ []byte, value []byte) error {
		var reward types.RewardByCategory
		if err := k.cdc.Unmarshal(value, &reward); err != nil {
			return err
		}

		rewards = append(rewards, reward)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRewardsByAddressAndCategoryResponse{Rewards: rewards, Pagination: pageRes}, nil
}

func (k Keeper) RewardsByCampaign(goCtx context.Context, req *types.QueryRewardsByCampaignRequest) (*types.QueryRewardsByCampaignResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var rewards []types.RewardByCampaign
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := k.getRewardsByCampaignStore(ctx)
	rewardStore := prefix.NewStore(store, types.GetRewardsByCampaignPrefix(req.CampaignUid))

	pageRes, err := query.Paginate(rewardStore, req.Pagination, func(_ []byte, value []byte) error {
		var reward types.RewardByCampaign
		if err := k.cdc.Unmarshal(value, &reward); err != nil {
			return err
		}

		rewards = append(rewards, reward)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRewardsByCampaignResponse{Rewards: rewards, Pagination: pageRes}, nil
}
