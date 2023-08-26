package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	housetypes "github.com/sge-network/sge/x/house/types"
	"github.com/sge-network/sge/x/subaccount/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m msgServer) HouseDeposit(goCtx context.Context, msg *types.MsgHouseDeposit) (*types.MsgHouseDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if subaccount exists
	subAccountAddr, exists := m.keeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(msg.Msg.Creator))
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	if err := m.houseDeposit(ctx, msg.Msg); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deposit")
	}

	// get subaccount balance, and check if it can spend
	balance, exists := m.keeper.GetBalance(ctx, subAccountAddr)
	if !exists {
		panic("data corruption: subaccount balance not found")
	}

	err := balance.Spend(msg.Msg.Amount)
	if err != nil {
		return nil, err
	}

	// send house deposit from subaccount on behalf of the owner
	participationIndex, err := m.keeper.houseKeeper.Deposit(
		ctx,
		subAccountAddr.String(),
		subAccountAddr.String(),
		msg.Msg.MarketUID,
		msg.Msg.Amount,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to deposit")
	}

	// update subaccount balance
	m.keeper.SetBalance(ctx, subAccountAddr, balance)

	// emit event
	msg.Msg.EmitEvent(&ctx, subAccountAddr.String(), participationIndex)

	return &types.MsgHouseDepositResponse{
		Response: &housetypes.MsgDepositResponse{
			MarketUID:          msg.Msg.MarketUID,
			ParticipationIndex: participationIndex,
		},
	}, nil
}

// TODO: This is a copy of the Deposit function from x/house/keeper/msg_server_deposit.go
func (m msgServer) houseDeposit(ctx sdk.Context, msg *housetypes.MsgDeposit) error {
	params := m.keeper.houseKeeper.GetParams(ctx)
	if err := msg.ValidateSanity(ctx, &params); err != nil {
		return sdkerrors.Wrap(err, "invalid deposit")
	}

	var payload housetypes.DepositTicketPayload
	if err := m.keeper.ovmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Ticket, &payload); err != nil {
		return sdkerrors.Wrapf(housetypes.ErrInTicketVerification, "%s", err)
	}

	if payload.DepositorAddress != "" {
		return sdkerrors.Wrapf(housetypes.ErrInTicketPayloadValidation, "in subaccount the depositor address must be empty")
	}

	depositorAddr := msg.Creator

	if err := payload.Validate(depositorAddr); err != nil {
		return sdkerrors.Wrapf(housetypes.ErrInTicketPayloadValidation, "%s", err)
	}

	return nil
}

func (m msgServer) HouseWithdraw(goCtx context.Context, withdraw *types.MsgHouseWithdraw) (*types.MsgHouseWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if subaccount exists
	subAccountAddr, exists := m.keeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(withdraw.Msg.Creator))
	if !exists {
		return nil, types.ErrSubaccountDoesNotExist
	}

	subAccountBalance, exists := m.keeper.GetBalance(ctx, subAccountAddr)
	if !exists {
		panic("data corruption: subaccount balance not found")
	}

	withdrawable, resp, err := m.houseWithdraw(ctx, withdraw.Msg, subAccountAddr)
	if err != nil {
		return nil, err
	}

	err = subAccountBalance.Unspend(withdrawable)
	if err != nil {
		panic("data corruption: it must be possible to unspend an house withdrawal")
	}

	m.keeper.SetBalance(ctx, subAccountAddr, subAccountBalance)
	return &types.MsgHouseWithdrawResponse{
		Response: resp,
	}, nil
}

func (m msgServer) houseWithdraw(ctx sdk.Context, msg *housetypes.MsgWithdraw, subAccAddr sdk.AccAddress) (math.Int, *housetypes.MsgWithdrawResponse, error) {
	var payload housetypes.WithdrawTicketPayload
	if err := m.keeper.ovmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Ticket, &payload); err != nil {
		return math.Int{}, nil, sdkerrors.Wrapf(housetypes.ErrInTicketVerification, "%s", err)
	}

	if payload.DepositorAddress != "" {
		return math.Int{}, nil, status.Errorf(codes.InvalidArgument, "in subaccount the depositor address must be empty")
	}

	if err := payload.Validate(msg.Creator); err != nil {
		return math.Int{}, nil, sdkerrors.Wrapf(housetypes.ErrInTicketPayloadValidation, "%s", err)
	}

	// Get the deposit object
	deposit, found := m.keeper.houseKeeper.GetDeposit(ctx, subAccAddr.String(), msg.MarketUID, msg.ParticipationIndex)
	if !found {
		return math.Int{}, nil, sdkerrors.Wrapf(housetypes.ErrDepositNotFound, ": %s, %d", msg.MarketUID, msg.ParticipationIndex)
	}

	withdrawable, err := m.keeper.obKeeper.CalcWithdrawalAmount(ctx,
		subAccAddr.String(),
		msg.MarketUID,
		msg.ParticipationIndex,
		msg.Mode,
		deposit.TotalWithdrawalAmount,
		msg.Amount,
	)
	if err != nil {
		return math.Int{}, nil, sdkerrors.Wrapf(housetypes.ErrInTicketVerification, "%s", err)
	}

	id, err := m.keeper.houseKeeper.Withdraw(ctx, deposit, msg.Creator, subAccAddr.String(), msg.MarketUID,
		msg.ParticipationIndex, msg.Mode, withdrawable)
	if err != nil {
		return math.Int{}, nil, sdkerrors.Wrap(err, "process withdrawal")
	}

	msg.EmitEvent(&ctx, subAccAddr.String(), id)

	return withdrawable, &housetypes.MsgWithdrawResponse{
		ID:                 id,
		MarketUID:          msg.MarketUID,
		ParticipationIndex: msg.ParticipationIndex,
	}, nil
}
