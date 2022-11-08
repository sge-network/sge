package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/sportevent/types"
)

// AddEvent accepts ticket containing multiple creation events and return batch response after processing
func (k msgServer) AddEvent(goCtx context.Context, msg *types.MsgAddEvent) (*types.SportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var sportData types.SportEvent
	if err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &sportData); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	if err := k.validateEventAdd(ctx, &sportData); err != nil {
		return nil, sdkerrors.Wrap(err, "validate add event")
	}

	_, found := k.Keeper.GetSportEvent(ctx, sportData.UID)
	if found {
		return nil, types.ErrEventAlreadyExist
	}

	sportData.Creator = msg.Creator
	k.Keeper.SetSportEvent(ctx, sportData)

	response := &types.SportResponse{
		Error: "",
		Data:  &sportData,
	}
	emitTransactionEvent(ctx, types.TypeMsgCreateSportEvents, response, msg.Creator)

	return response, nil
}

// validateEventAdd validates individual event acceptability
func (k msgServer) validateEventAdd(ctx sdk.Context, event *types.SportEvent) error {

	if err := validateEventTS(ctx, event); err != nil {
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
	initBetConstraints(event, params)

	if err := validateBetConstraints(event, &params); err != nil {
		return err
	}

	return nil
}

func initBetConstraints(event *types.SportEvent, params types.Params) {
	// init all Bet constraints if nil
	if event.BetConstraints == nil {
		event.BetConstraints = &types.EventBetConstraints{
			MinAmount: params.EventMinBetAmount,
			BetFee:    params.EventMinBetFee,
		}
		return
	}

	// init individual params if any one of them is nil
	if event.BetConstraints.BetFee.IsNil() {
		event.BetConstraints.BetFee = params.EventMinBetFee
	}
	if event.BetConstraints.MinAmount.IsNil() {
		event.BetConstraints.MinAmount = params.EventMinBetAmount
	}
}

func validateEventTS(ctx sdk.Context, event *types.SportEvent) error {
	if event.EndTS <= uint64(ctx.BlockTime().Unix()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid end timestamp for the sport event")
	}

	if event.StartTS >= event.EndTS || event.StartTS == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid start timestamp for the sport event, cannot be (greater than eql to EndTs) or 0")
	}

	return nil
}

func validateBetConstraints(event *types.SportEvent, params *types.Params) error {
	//check the validity constraints as there is no GT method on coin type
	if !(event.BetConstraints.BetFee.IsLT(params.EventMinBetFee) || event.BetConstraints.BetFee.IsEqual(params.EventMinBetFee)) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event bet fee is out of threshold limit")
	}

	if event.BetConstraints.MinAmount.IsNegative() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min amount can not be negative")
	}

	if event.BetConstraints.MinAmount.LT(params.EventMinBetAmount) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "event min bet amount is less than threshold")
	}

	return nil
}
