package keeper

import (
	"context"

	"github.com/spf13/cast"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/legacy/ovm/types"
)

// VotePubkeysChange is the main transaction of OVM to vote for a pubkeys list change proposal.
func (k msgServer) VotePubkeysChange(
	goCtx context.Context,
	msg *types.MsgVotePubkeysChangeRequest,
) (*types.MsgVotePubkeysChangeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	keyVault, found := k.GetKeyVault(ctx)
	if !found {
		return nil, types.ErrKeyVaultNotFound
	}

	// voter index is out of range of current public keys in the vault
	if cast.ToUint32(len(keyVault.PublicKeys)) <= msg.VoterKeyIndex {
		return nil, sdkerrors.Wrapf(
			sdkerrtypes.ErrInvalidPubKey,
			"index is not in the valid range of public keys indices",
		)
	}

	// pick the voter public key from the public key list of the vault
	voterPubKey := keyVault.PublicKeys[msg.VoterKeyIndex]

	// the ticket should be signed by the private key associated with the
	// stored public key in the public keys list of the vault at index
	// that is proposed in the message body
	payload := types.ProposalVotePayload{}
	err := k.verifyTicketWithKeyUnmarshal(goCtx, msg.Ticket, &payload, voterPubKey)
	if err != nil {
		return nil, sdkerrors.Wrapf(
			sdkerrtypes.ErrInvalidRequest,
			"ticket should be signed by the provided pub key: %s",
			err,
		)
	}

	err = payload.Validate()
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "ticket payload is not valid %s", err)
	}

	// get public key change proposal
	proposal, found := k.GetPubkeysChangeProposal(
		ctx,
		types.ProposalStatus_PROPOSAL_STATUS_ACTIVE,
		payload.ProposalId,
	)
	if !found {
		return nil, sdkerrors.Wrapf(
			sdkerrtypes.ErrInvalidRequest,
			"proposal not fount with id %d",
			&payload.ProposalId,
		)
	}

	// check if the voter public key already voted or not
	// a public key can participate in the voting just for one time
	for _, voter := range proposal.Votes {
		if voter.PublicKey == voterPubKey {
			return nil, sdkerrors.Wrapf(
				sdkerrtypes.ErrInvalidRequest,
				"vote already set for this pubkey %s",
				msg.Creator,
			)
		}
	}

	// append new vote to the votes of the proposal
	vote := types.NewVote(
		voterPubKey,
		payload.Vote,
	)
	proposal.Votes = append(proposal.Votes, vote)

	// set proposal with updated votes in the module state
	k.SetPubkeysChangeProposal(ctx, proposal)

	msg.EmitEvent(&ctx, proposal.Id, vote.PublicKey, vote.Vote)

	return &types.MsgVotePubkeysChangeResponse{Success: true}, nil
}
