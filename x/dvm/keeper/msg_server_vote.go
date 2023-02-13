package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/dvm/types"
)

// VotePubkeysChange is the main transaction of DVM to vote for a pubkeys list change proposal.
func (k msgServer) VotePubkeysChange(goCtx context.Context, msg *types.MsgVotePubkeysChangeRequest) (*types.MsgVotePubkeysChangeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	payload := types.ProposalVotePayload{}
	err := k.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload)
	if err != nil {
		return nil, err
	}

	proposal, found := k.GetActivePubkeysChangeProposal(ctx, payload.ProposalId)

	if !found {
		return nil, types.ErrNoPublicKeysFound
	}

	for _, approver := range proposal.ApprovedBy {
		if approver == msg.Creator {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "approval already set for this address %s", msg.Creator)
		}
	}

	proposal.ApprovedBy = append(proposal.ApprovedBy, msg.Creator)

	return &types.MsgVotePubkeysChangeResponse{Success: true}, nil
}
