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

	payload := types.PubkeysChangeProposalPayload{}
	err := k.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload)
	if err != nil {
		return nil, err
	}

	var pp interface{}
	err = k.VerifyTicketUnmarshal(goCtx, msg.Ticket, &pp)
	if err != nil {
		return nil, err
	}

	// remove duplicates in public keys
	payload.PublicKeys = utils.RemoveDuplicateStrs(payload.PublicKeys)

	err = payload.Validate(payload.PublicKeys, payload.LeaderIndex)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ticket payload is not valid %s", err)
	}

	stats := k.GetProposalStats(ctx)
	stats.PubkeysChangeCount++

	// set proposal
	k.Keeper.SetPubkeysChangeProposal(ctx,
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
