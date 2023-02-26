package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/dvm/types"
)

// SubmitPubkeysChangeProposal is the main transaction of DVM to add or delete the keys to the chain.
func (k msgServer) SubmitPubkeysChangeProposal(goCtx context.Context, msg *types.MsgSubmitPubkeysChangeProposalRequest) (*types.MsgSubmitPubkeysChangeProposalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	keys, found := k.GetKeyVault(ctx)

	if !found {
		return nil, types.ErrNoPublicKeysFound
	}

	payload := types.PubkeysChangeProposalPayload{}
	err := k.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload)
	if err != nil {
		return nil, err
	}

	// remove duplicate additions and deletions
	payload.Additions = utils.RemoveDuplicateStrs(payload.Additions)
	payload.Deletions = utils.RemoveDuplicateStrs(payload.Deletions)

	err = payload.Validate(keys.PublicKeys)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ticket payload is not valid %s", err)
	}

	stats := k.GetProposalStats(ctx)
	stats.PubkeysChangeCount++

	// set active proposal
	k.Keeper.SetActivePubkeysChangeProposal(ctx,
		types.NewPublicKeysChangeProposal(
			stats.PubkeysChangeCount,
			msg.Creator,
			payload,
			ctx.BlockTime().Unix(),
		),
	)

	// set proposal statistics
	k.Keeper.SetProposalStats(ctx, stats)

	return &types.MsgSubmitPubkeysChangeProposalResponse{Success: true}, nil
}
