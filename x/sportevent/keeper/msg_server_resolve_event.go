package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/sportevent/types"
)

// ResolveSportEvent accepts ticket containing multiple resolution events and return batch response after processing
func (k msgServer) ResolveSportEvent(goCtx context.Context, msg *types.MsgResolveSportEventRequest) (*types.MsgResolveSportEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var resolutionPayload types.SportEventResolutionTicketPayload
	err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &resolutionPayload)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	if err := k.processSportEventResolution(ctx, &resolutionPayload); err != nil {
		return nil, sdkerrors.Wrap(err, "process resolution event")
	}

	sportEvent, _ := k.getSportEventToResolve(ctx, resolutionPayload)
	response := &types.MsgResolveSportEventResponse{
		Data: &sportEvent,
	}
	emitTransactionEvent(ctx, types.TypeMsgResolveSportEvents, response.Data.UID, msg.Creator)
	return response, nil
}

func (k msgServer) processSportEventResolution(ctx sdk.Context, resolutionPayload *types.SportEventResolutionTicketPayload) error {
	sportEvent, err := k.getSportEventToResolve(ctx, *resolutionPayload)
	if err != nil {
		return sdkerrors.Wrap(err, "getting sport-event")
	}

	if err := extractWinnerOddsIDs(&sportEvent, resolutionPayload); err != nil {
		return sdkerrors.Wrap(err, "extract winner odds id")
	}

	if err := k.Keeper.ResolveSportEvent(ctx, resolutionPayload); err != nil {
		return sdkerrors.Wrap(err, "resolve sport-event")
	}

	return nil
}

func (k msgServer) getSportEventToResolve(ctx sdk.Context, event types.SportEventResolutionTicketPayload) (types.SportEvent, error) {
	if err := validateEventResolution(event); err != nil {
		return types.SportEvent{}, err
	}

	sportEvent, found := k.Keeper.GetSportEvent(ctx, event.UID)
	if !found {
		return types.SportEvent{}, types.ErrEventNotFound
	}

	if sportEvent.Status != types.SportEventStatus_SPORT_EVENT_STATUS_UNSPECIFIED {
		return types.SportEvent{}, types.ErrEventIsNotPending
	}

	return sportEvent, nil
}

func extractWinnerOddsIDs(sportEvent *types.SportEvent, event *types.SportEventResolutionTicketPayload) error {
	if event.Status == types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED {
		if event.ResolutionTS < sportEvent.StartTS {
			return types.ErrResolutionTimeLessTnStart
		}

		validWinnerOdds := true
		for _, wid := range event.WinnerOddsUIDs {
			validWinnerOdds = false
			for _, o := range sportEvent.Odds {
				if o.UID == wid {
					validWinnerOdds = true
				}
			}
			if !validWinnerOdds {
				break
			}
		}

		if !validWinnerOdds {
			return types.ErrInvalidWinnerOdd
		}
	}

	return nil
}

// validateEventResolution validates event resolution payload
func validateEventResolution(event types.SportEventResolutionTicketPayload) error {
	if event.Status == types.SportEventStatus_SPORT_EVENT_STATUS_UNSPECIFIED || event.Status == types.SportEventStatus_SPORT_EVENT_STATUS_INVALID {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "resolution status passed for the sports event is invalid")
	}

	if event.ResolutionTS == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid resolution timestamp for the sport-event")
	}

	if !utils.IsValidUID(event.UID) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid uid for the sport-event")
	}

	if event.Status == types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED && len(event.WinnerOddsUIDs) < 1 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "not provided enough winner odds for the sports event")
	}

	for _, wid := range event.WinnerOddsUIDs {
		if !utils.IsValidUID(wid) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "odds-uid passed is invalid")
		}
	}

	if event.Status > types.SportEventStatus_SPORT_EVENT_STATUS_RESULT_DECLARED {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid event resolution status ")
	}

	return nil
}
