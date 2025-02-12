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

func (k Keeper) Promoters(goCtx context.Context, req *types.QueryPromotersRequest) (*types.QueryPromotersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var promoters []types.Promoter
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	promoterStore := prefix.NewStore(store, types.PromoterKeyPrefix)

	pageRes, err := query.Paginate(promoterStore, req.Pagination, func(_ []byte, value []byte) error {
		var promoter types.Promoter
		if err := k.cdc.Unmarshal(value, &promoter); err != nil {
			return err
		}

		promoters = append(promoters, promoter)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPromotersResponse{Promoter: promoters, Pagination: pageRes}, nil
}

func (k Keeper) PromoterByAddress(goCtx context.Context, req *types.QueryPromoterByAddressRequest) (*types.QueryPromoterByAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	promoterByAddress, found := k.GetPromoterByAddress(
		ctx,
		req.Addr,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "promoter with the provided address not found")
	}

	promoter, found := k.GetPromoter(
		ctx,
		promoterByAddress.PromoterUID,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "promoter with the provided uid found")
	}

	return &types.QueryPromoterByAddressResponse{Promoter: promoter}, nil
}
