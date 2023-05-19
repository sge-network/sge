package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/rewards/types"
)

func (k msgServer) RewardUser(goCtx context.Context, msg *types.MsgRewardUser) (*types.MsgRewardUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if err := msg.ValidateSanity(ctx, &params); err != nil {
		return nil, sdkerrors.Wrap(err, "Invalid Reward Requests")
	}

	// TODO: See if want to use OVM
	//if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
	//	return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	//}
	for _, awardee := range msg.Reward.Awardees {
		k.Keeper.RewardUser(ctx, msg.Creator, msg.Reward.RewardType.String(), awardee.Amount, awardee.Address)
	}
	return &types.MsgRewardUserResponse{}, nil
}
