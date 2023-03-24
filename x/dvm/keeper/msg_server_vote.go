package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/dvm/types"
	"github.com/spf13/cast"
)

// VotePubkeysChange is the main transaction of DVM to vote for a pubkeys list change proposal.
func (k msgServer) VotePubkeysChange(goCtx context.Context, msg *types.MsgVotePubkeysChangeRequest) (*types.MsgVotePubkeysChangeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	payload := types.ProposalVotePayload{}

	keyVault, found := k.GetKeyVault(ctx)
	if !found {
		return nil, types.ErrKeyVaultNotFound
	}

	if cast.ToUint32(len(keyVault.PublicKeys)) <= msg.VoterKeyIndex {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "index is not in the valid range of public keys indices")
	}

	voterPubKey := keyVault.PublicKeys[msg.VoterKeyIndex]

	err := k.verifyTicketWithKeyUnmarshal(goCtx, msg.Ticket, &payload, voterPubKey)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ticket should be signed by the provided pub key: %s", err)
	}

	proposal, found := k.GetPubkeysChangeProposal(ctx, types.ProposalStatus_PROPOSAL_STATUS_ACTIVE, payload.ProposalId)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "proposal not fount with id %d", &payload.ProposalId)
	}

	for _, voter := range proposal.Votes {
		if voter.PublicKey == voterPubKey {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "vote already set for this pubkey %s", msg.Creator)
		}
	}

	proposal.Votes = append(proposal.Votes,
		types.NewVote(
			voterPubKey,
			payload.Vote,
		),
	)

	// set in the state
	k.SetPubkeysChangeProposal(ctx, proposal)

	return &types.MsgVotePubkeysChangeResponse{Success: true}, nil
}
