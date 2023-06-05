package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/ovm/types"
	"github.com/spf13/cast"
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
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ticket payload is not valid %s", err)
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

	emitSumitProposalEvent(
		ctx,
		types.TypeMsgCreatePubKeyChaneProposal,
		proposal.Id,
		proposal.Creator,
		proposal.Modifications.String(),
	)

	return &types.MsgSubmitPubkeysChangeProposalResponse{Success: true}, nil
}

func emitSumitProposalEvent(
	ctx sdk.Context,
	emitType string,
	id uint64,
	creator string,
	content string,
) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			emitType,
			sdk.NewAttribute(types.AttributeKeyPubkeysChangeProposalID, cast.ToString(id)),
			sdk.NewAttribute(types.AttributeKeyCreator, creator),
			sdk.NewAttribute(types.AttributeKeyContent, content),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, emitType),
			sdk.NewAttribute(sdk.AttributeKeySender, creator),
		),
	})
}
