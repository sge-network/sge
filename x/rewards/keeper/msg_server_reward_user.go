package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/rewards/types"
)

func (k msgServer) RewardUser(goCtx context.Context, msg *types.MsgRewardUser) (*types.MsgRewardUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if err := msg.ValidateSanity(ctx, &params); err != nil {
		return nil, sdkerrors.Wrap(err, "Invalid Reward Requests")
	}
	var rewardsTicketPayload types.RewardsTicketPayload

	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &rewardsTicketPayload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	emitTransactionEvent(ctx, "Reward User Initiate", len(msg.Reward.Awardees), msg.Creator)

	err := k.Keeper.RewardUsers(ctx, msg, rewardsTicketPayload.IncentiveUID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Something failed")
	}

	emitTransactionEvent(ctx, "Reward User Done", len(msg.Reward.Awardees), msg.Creator)
	return &types.MsgRewardUserResponse{}, nil
}

func emitTransactionEvent(ctx sdk.Context, emitType string, lenAwardees int, creator string) {
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			emitType,
			sdk.NewAttribute("Awardees Len ", strconv.Itoa(lenAwardees)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeyAction, emitType),
			sdk.NewAttribute(sdk.AttributeKeySender, creator),
		),
	})
}
