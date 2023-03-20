package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/dvm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ActivePublicKeysChangeProposal returns a specific active proposal by its UID
func (k Keeper) ActivePublicKeysChangeProposal(c context.Context, req *types.QueryActivePublicKeysChangeProposalRequest) (*types.QueryActivePublicKeysChangeProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetActivePubkeysChangeProposal(
		ctx,
		req.Id,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryActivePublicKeysChangeProposalResponse{Proposal: val}, nil
}

// ActivePublicKeysChangeProposal returns list of the active pubkeys change proposal
func (k Keeper) ActivePublicKeysChangeProposals(goCtx context.Context, req *types.QueryActivePublicKeysChangeProposalsRequest) (*types.QueryActivePublicKeysChangeProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var proposals []types.PublicKeysChangeProposal
	ctx := sdk.UnwrapSDKContext(goCtx)

	MarketStore := k.getActivePubKeysChangeProposalStore(ctx)

	pageRes, err := query.Paginate(MarketStore, req.Pagination, func(key []byte, value []byte) error {
		var proposal types.PublicKeysChangeProposal
		if err := k.cdc.Unmarshal(value, &proposal); err != nil {
			return err
		}

		proposals = append(proposals, proposal)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryActivePublicKeysChangeProposalsResponse{Proposals: proposals, Pagination: pageRes}, nil
}

// PublicKeysChangeProposal returns a specific finished proposal by its UID
func (k Keeper) FinishedPublicKeysChangeProposal(c context.Context, req *types.QueryFinishedPublicKeysChangeProposalRequest) (*types.QueryFinishedPublicKeysChangeProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetFinishedPubkeysChangeProposal(
		ctx,
		req.Id,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryFinishedPublicKeysChangeProposalResponse{Proposal: val}, nil
}

// FinishedPublicKeysChangeProposal returns list of the finished pubkeys change proposal
func (k Keeper) FinishedPublicKeysChangeProposals(goCtx context.Context, req *types.QueryFinishedPublicKeysChangeProposalsRequest) (*types.QueryFinishedPublicKeysChangeProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var proposals []types.PublicKeysChangeFinishedProposal
	ctx := sdk.UnwrapSDKContext(goCtx)

	MarketStore := k.getActivePubKeysChangeProposalStore(ctx)

	pageRes, err := query.Paginate(MarketStore, req.Pagination, func(key []byte, value []byte) error {
		var proposal types.PublicKeysChangeFinishedProposal
		if err := k.cdc.Unmarshal(value, &proposal); err != nil {
			return err
		}

		proposals = append(proposals, proposal)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryFinishedPublicKeysChangeProposalsResponse{Proposals: proposals, Pagination: pageRes}, nil
}
