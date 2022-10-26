package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/bet/types"
)

func (k msgServer) PlaceBet(goCtx context.Context, msg *types.MsgPlaceBet) (*types.MsgPlaceBetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetBet(ctx, msg.Bet.UID)
	if isFound {
		return nil, types.ErrDuplicateUID
	}

	ticketData := &types.BetOdds{}
	err := k.dvmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), msg.Bet.Ticket, &ticketData)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	if err = types.TicketFieldsValidation(ticketData); err != nil {
		return nil, err
	}

	bet, err := types.NewBet(msg.Creator, msg.Bet, ticketData)
	if err != nil {
		return nil, err
	}

	if err := k.Keeper.PlaceBet(ctx, bet); err != nil {
		return nil, err
	}
	return &types.MsgPlaceBetResponse{}, nil
}

func (k msgServer) PlaceBetSlip(goCtx context.Context, msg *types.MsgPlaceBetSlip) (*types.MsgPlaceBetSlipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var (
		successBetUIDsList     = []string{}
		failureBetUIDsErrorMap = make(map[string]string)
	)

	for _, betFields := range msg.Bets {
		// Check if the value already exists
		_, isFound := k.GetBet(ctx, betFields.UID)
		if isFound {
			failureBetUIDsErrorMap[betFields.UID] = types.ErrDuplicateUID.Error()
			continue
		}

		// Fields Validation of single bet
		if err := types.BetFieldsValidation(betFields); err != nil {
			failureBetUIDsErrorMap[betFields.UID] = err.Error()
			continue
		}

		ticketData := &types.BetOdds{}
		err := k.dvmKeeper.VerifyTicketUnmarshal(sdk.WrapSDKContext(ctx), betFields.Ticket, &ticketData)
		if err != nil {
			failureBetUIDsErrorMap[betFields.UID] = err.Error()
			continue
		}

		if err = types.TicketFieldsValidation(ticketData); err != nil {
			failureBetUIDsErrorMap[betFields.UID] = err.Error()
			continue
		}

		bet, err := types.NewBet(msg.Creator, betFields, ticketData)
		if err != nil {
			return nil, err
		}

		if err := k.Keeper.PlaceBet(ctx, bet); err != nil {
			failureBetUIDsErrorMap[bet.UID] = err.Error()
			continue
		}

		successBetUIDsList = append(successBetUIDsList, bet.UID)
	}
	response := &types.MsgPlaceBetSlipResponse{
		SuccessfulBetUIDsList: successBetUIDsList,
		FailedBetUIDsErrorMap: failureBetUIDsErrorMap,
	}

	return response, nil
}

func (k msgServer) SettleBet(goCtx context.Context, msg *types.MsgSettleBet) (*types.MsgSettleBetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.Keeper.SettleBet(ctx, msg.BetUID); err != nil {
		return nil, err
	}
	return &types.MsgSettleBetResponse{}, nil
}

func (k msgServer) SettleBetBulk(goCtx context.Context, msg *types.MsgSettleBetBulk) (*types.MsgSettleBetBulkResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	successBetUIDs := []string{}
	failedBetUIDs := make(map[string]string)

	for _, betUID := range msg.BetUIDs {
		if err := k.Keeper.SettleBet(ctx, betUID); err == nil {
			successBetUIDs = append(successBetUIDs, betUID)
		} else {
			failedBetUIDs[betUID] = err.Error()
		}
	}

	settleBetBulkResp := &types.MsgSettleBetBulkResponse{
		SuccessfulBetUIDsList: successBetUIDs,
		FailedBetUIDsErrorMap: failedBetUIDs,
	}

	return settleBetBulkResp, nil
}
