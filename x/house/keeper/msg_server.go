package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/house/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the house MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// Deposit defines a method for performing a deposit of coins to become part of the house correspondifg to a sport event.
func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	depAddr, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidDepositor, err)
	}

	if err := k.validateDeposit(ctx, msg); err != nil {
		return nil, sdkerrors.Wrap(err, "validate deposit")
	}

	err, participantId := k.Keeper.Deposit(ctx, depAddr, msg.SportEventUid, msg.Amount.Amount)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "process deposit")
	}

	return &types.MsgDepositResponse{
		SportEventUid: msg.SportEventUid,
		ParticipantId: participantId,
	}, nil
}

// Withdraw defines a method for performing a withdrawal of coins of unused amount corresponding to a deposit.
func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	depAddr, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, types.ErrTextInvalidDepositor, err)
	}

	err, withdrawalId := k.Keeper.Withdraw(ctx, depAddr, msg.SportEventUid, msg.ParticipantId, msg.Mode, msg.Amount.Amount)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "process withdrawal")
	}

	return &types.MsgWithdrawResponse{
		SportEventUid:    msg.SportEventUid,
		ParticipantId:    msg.ParticipantId,
		WithdrawalNumber: withdrawalId,
	}, nil
}
