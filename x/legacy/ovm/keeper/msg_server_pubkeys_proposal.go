package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/legacy/ovm/types"
)

// SubmitPubkeysChangeProposal is the main transaction of OVM to add or delete the keys to the chain.
func (k msgServer) SubmitPubkeysChangeProposal(
	goCtx context.Context,
	msg *types.MsgSubmitPubkeysChangeProposalRequest,
) (*types.MsgSubmitPubkeysChangeProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	keyVault, found := k.GetKeyVault(ctx)
	if !found {
		return nil, types.ErrKeyVaultNotFound
	}

	payload := types.PubkeysChangeProposalPayload{}
	err := k.verifyTicketWithKeyUnmarshal(goCtx, msg.Ticket, &payload, keyVault.PublicKeys...)
	if err != nil {
		return nil, err
	}

	// remove duplicates in public keys
	payload.PublicKeys = utils.RemoveDuplicateStrs(payload.PublicKeys)

	err = payload.Validate(payload.LeaderIndex)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "ticket payload is not valid %s", err)
	}

	stats := k.GetProposalStats(ctx)
	stats.PubkeysChangeCount++

	// set proposal
	proposal := types.NewPublicKeysChangeProposal(
		stats.PubkeysChangeCount,
		msg.Creator,
		payload,
		ctx.BlockTime().Unix(),
	)
	k.Keeper.SetPubkeysChangeProposal(ctx, proposal)

	// set proposal statistics
	k.Keeper.SetProposalStats(ctx, stats)

	msg.EmitEvent(&ctx, proposal.Id, proposal.Modifications.String())

	return &types.MsgSubmitPubkeysChangeProposalResponse{Success: true}, nil
}
