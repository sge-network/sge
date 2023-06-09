package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/house/types"
)

// Withdraw performs withdrawal of unused tokens corresponding to a deposit.
func (k msgServer) Withdraw(goCtx context.Context,
	msg *types.MsgWithdraw,
) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var payload types.WithdrawTicketPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	depositorAddr := msg.Creator
	if payload.DepositorAddress != "" {
		if err := k.ValidateMsgAuthorization(ctx, msg.Creator, payload.DepositorAddress, msg); err != nil {
			return nil, err
		}
		depositorAddr = payload.DepositorAddress
	}

	if err := payload.Validate(depositorAddr); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	id, err := k.Keeper.Withdraw(ctx, msg.Creator, depositorAddr, msg.MarketUID,
		msg.ParticipationIndex, msg.Mode, msg.Amount)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "process withdrawal")
	}

	msg.EmitEvent(&ctx, depositorAddr)

	return &types.MsgWithdrawResponse{
		ID:                 id,
		MarketUID:          msg.MarketUID,
		ParticipationIndex: msg.ParticipationIndex,
	}, nil
}
