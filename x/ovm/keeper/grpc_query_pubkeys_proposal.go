package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sge-network/sge/consts"
	"github.com/sge-network/sge/x/ovm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PublicKeysChangeProposal returns a specific proposal by its id and status
func (k Keeper) PublicKeysChangeProposal(
	c context.Context,
	req *types.QueryPublicKeysChangeProposalRequest,
) (*types.QueryPublicKeysChangeProposalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetPubkeysChangeProposal(
		ctx,
		req.Status,
		req.Id,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryPublicKeysChangeProposalResponse{Proposal: val}, nil
}

// PublicKeysChangeProposals returns list of the pubkeys change proposal
func (k Keeper) PublicKeysChangeProposals(
	goCtx context.Context,
	req *types.QueryPublicKeysChangeProposalsRequest,
) (*types.QueryPublicKeysChangeProposalsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, consts.ErrTextInvalidRequest)
	}

	var proposals []types.PublicKeysChangeProposal
	ctx := sdk.UnwrapSDKContext(goCtx)

	marketStore := k.getPubKeysChangeProposalStore(ctx)
	proposalStore := prefix.NewStore(marketStore, types.PubkeysChangeProposalPrefix(req.Status))

	pageRes, err := query.Paginate(proposalStore, req.Pagination, func(key []byte, value []byte) error {
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

	return &types.QueryPublicKeysChangeProposalsResponse{Proposals: proposals, Pagination: pageRes}, nil
}
