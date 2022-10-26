package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/sportevent/types"
)

// AddEvent accepts ticket containing multiple creation events and return batch response after processing
func (k msgServer) AddEvent(goCtx context.Context, msg *types.MsgAddEvent) (*types.MsgSportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	type input struct {
		Events []types.SportEvent `json:"events"`
	}
	var sportData input
	if err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &sportData); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	uniqueEvents := make(map[string]struct{}, len(sportData.Events))
	failed := make([]*types.FailedEvent, 0)
	success := make([]string, 0, len(sportData.Events))

	for _, event := range sportData.Events {
		if err := k.validateEventAdd(ctx, &event); err != nil {
			failed = append(failed, &types.FailedEvent{
				ID:  event.UID,
				Err: err.Error(),
			})
			continue
		}

		_, found := k.Keeper.GetSportEvent(ctx, event.UID)
		if found {
			failed = append(failed, &types.FailedEvent{
				ID:  event.UID,
				Err: types.ErrEventAlreadyExist.Error(),
			})
			continue
		}

		if _, duplicateFound := uniqueEvents[event.UID]; duplicateFound {
			failed = append(failed, &types.FailedEvent{
				ID:  event.UID,
				Err: types.ErrDuplicateEventsExist.Error(),
			})
			continue
		}

		event.Creator = msg.Creator
		uniqueEvents[event.UID] = struct{}{}

		k.Keeper.SetSportEvent(ctx, event)
		success = append(success, event.UID)
	}

	response := &types.MsgSportResponse{
		SuccessEvents: success,
		FailedEvents:  failed,
	}
	emitTransactionEvent(ctx, types.TypeMsgCreateSportEvents, response, msg.Creator)

	return response, nil
}

// validateEventAdd validates individual event acceptability
func (k msgServer) validateEventAdd(ctx sdk.Context, event *types.SportEvent) error {

	if err := validateEventTS(event); err != nil {
		return err
	}

	if !utils.IsValidUID(event.UID) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid uid for the sport event")
	}

	if len(event.OddsUIDs) < 2 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "not provided enough odds for the event")
	}

	oddsSet := make(map[string]struct{}, len(event.OddsUIDs))
	for _, oddsUID := range event.OddsUIDs {
		if !utils.IsValidUID(oddsUID) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "odds-uid passed is invalid")
		}
		if _, exist := oddsSet[oddsUID]; exist {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate odds-uid in request")
		}
		oddsSet[oddsUID] = struct{}{}
	}

	params := k.GetParams(ctx)

	//init all Bet constraints if nil
	if event.BetConstraints == nil {
		event.BetConstraints = &types.EventBetConstraints{
			MaxBetCap:          params.EventMaxBetCap,
			MinAmount:          params.EventMinBetAmount,
			BetFee:             params.EventMinBetFee,
			CurrentTotalAmount: sdk.NewInt(0),
		}
		return nil
	}

	//init individual params if any one of them is nil
	if event.BetConstraints.BetFee.IsNil() {
		event.BetConstraints.BetFee = params.EventMinBetFee
	}
	if event.BetConstraints.MaxBetCap.IsNil() {
		event.BetConstraints.MaxBetCap = params.EventMaxBetCap
	}
	if event.BetConstraints.MinAmount.IsNil() {
		event.BetConstraints.MinAmount = params.EventMinBetAmount
	}
	event.BetConstraints.CurrentTotalAmount = sdk.NewInt(0)

	if err := validateBetConstraints(event, &params); err != nil {
		return err
	}

	return nil
}

func validateEventTS(event *types.SportEvent) error {
	if event.EndTS <= uint64(time.Now().Unix()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid end timestamp for the sport event")
	}

	if event.StartTS >= event.EndTS || event.StartTS == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid start timestamp for the sport event, cannot be (greater than eql to EndTs) or 0")
	}

	return nil
}

func validateBetConstraints(event *types.SportEvent, params *types.Params) error {
	if !(event.BetConstraints.BetFee.IsLT(params.EventMinBetFee) || event.BetConstraints.BetFee.IsEqual(params.EventMinBetFee)) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event bet fee is out of threshold limit")
	}

	if event.BetConstraints.MinAmount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min amount can not be negative")
	}

	if event.BetConstraints.MinAmount.LT(params.EventMinBetAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min bet amount is less than threshold")
	}

	if event.BetConstraints.MaxBetCap.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event max bet can not be negative")
	}

	if event.BetConstraints.MaxBetCap.GT(params.EventMaxBetCap) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event max bet cap is greater than threshold")
	}

	if event.BetConstraints.MinAmount.GTE(event.BetConstraints.MaxBetCap) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "min bet amount cannot be greater than or equals to to max bet capacity")
	}

	return nil
}
