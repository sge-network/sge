package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/house/types"
)

// Deposit performs deposit operation to participate as a house in a specific market/order book
func (k msgServer) Deposit(goCtx context.Context,
	msg *types.MsgDeposit,
) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	depositorAddr, err := k.ParseTicketAndValidate(goCtx, ctx, msg, true)
	if err != nil {
		return nil, err
	}

	participationIndex, err := k.Keeper.Deposit(
		ctx,
		msg.Creator,
		depositorAddr,
		msg.MarketUID,
		msg.Amount,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deposit")
	}

	msg.EmitEvent(&ctx, depositorAddr, participationIndex)

	return &types.MsgDepositResponse{
		MarketUID:          msg.MarketUID,
		ParticipationIndex: participationIndex,
	}, nil
}

// ParseTicketAndDeposit parses the deposit payload ticket and deposit.
func (k Keeper) ParseTicketAndValidate(
	goCtx context.Context,
	ctx sdk.Context,
	msg *types.MsgDeposit,
	authzAllowed bool,
) (string, error) {
	params := k.GetParams(ctx)
	if err := msg.ValidateSanity(ctx, &params); err != nil {
		return "", sdkerrors.Wrap(err, "invalid deposit")
	}

	var payload types.DepositTicketPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &payload); err != nil {
		return "", sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	depositorAddr := msg.Creator
	if payload.DepositorAddress != "" &&
		payload.DepositorAddress != msg.Creator {
		if !authzAllowed {
			return "", sdkerrors.Wrapf(types.ErrAuthorizationNotAllowed, "%s")
		}
		if err := k.ValidateMsgAuthorization(ctx, msg.Creator, payload.DepositorAddress, msg); err != nil {
			return "", err
		}
		depositorAddr = payload.DepositorAddress
	}

	if err := payload.Validate(depositorAddr); err != nil {
		return "", sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	return depositorAddr, nil
}
