package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/sge-network/sge/x/sportevent/types"
)

type resoveEventInput struct {
	Events []types.ResolutionEvent `json:"events"`
}

// ResolveEvent accepts ticket containing multiple resolution events and return batch response after processing
func (k msgServer) ResolveEvent(goCtx context.Context, msg *types.MsgResolveEvent) (*types.MsgSportResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var resolvedEvent resoveEventInput
	err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &resolvedEvent)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	success, failed := k.processEvents(ctx, &resolvedEvent)

	response := &types.MsgSportResponse{
		SuccessEvents: success,
		FailedEvents:  failed,
	}
	emitTransactionEvent(ctx, types.TypeMsgResolveSportEvents, response, msg.Creator)
	return response, nil
}

func (k msgServer) processEvents(ctx sdk.Context, resolvedEvent *resoveEventInput) ([]string, []*types.FailedEvent) {
	success := make([]string, 0, len(resolvedEvent.Events))
	failed := make([]*types.FailedEvent, 0)

	for _, event := range resolvedEvent.Events {
		sportEvent, err := k.getSportEvent(ctx, event)
		if err != nil {
			failed = append(failed, &types.FailedEvent{
				ID:  event.UID,
				Err: err.Error(),
			})
			continue
		}

		if err := extractWinnerOddsIDs(&sportEvent, &event); err != nil {
			failed = append(failed, &types.FailedEvent{
				ID:  event.UID,
				Err: err.Error(),
			})
			continue
		}

		if err := k.Keeper.ResolveSportEvents(ctx, []types.ResolutionEvent{event}); err != nil {
			failed = append(failed, &types.FailedEvent{
				ID:  event.UID,
				Err: err.Error(),
			})
			continue
		}
		success = append(success, sportEvent.UID)
	}

	return success, failed
}

func (k msgServer) getSportEvent(ctx sdk.Context, event types.ResolutionEvent) (types.SportEvent, error) {
	if err := validateResolutionEvent(event); err != nil {
		return types.SportEvent{}, err
	}

	sportEvent, found := k.Keeper.GetSportEvent(ctx, event.UID)
	if !found {
		return types.SportEvent{}, types.ErrEventNotFound
	}

	if sportEvent.Status != types.SportEventStatus_STATUS_PENDING {

		return types.SportEvent{}, types.ErrEventIsNotPending
	}

	return sportEvent, nil
}

func extractWinnerOddsIDs(sportEvent *types.SportEvent, event *types.ResolutionEvent) error {
	winnerOddUids := make(map[string][]byte, len(event.WinnerOddsUIDs))

	if event.Status == types.SportEventStatus_STATUS_RESULT_DECLARED {
		if event.ResolutionTS < sportEvent.StartTS {
			return types.ErrResolutionTimeLessTnStart
		}

		validWinnerOdds := true
		for key := range event.WinnerOddsUIDs {
			if !utils.StrSliceContains(sportEvent.OddsUIDs, key) {
				validWinnerOdds = false
				break
			}
			winnerOddUids[key] = nil
		}

		if !validWinnerOdds {
			return types.ErrInvalidWinnerOdd
		}
	}

	event.WinnerOddsUIDs = winnerOddUids
	return nil
}

// validateResolutionEvent validates individual event acceptability
func validateResolutionEvent(event types.ResolutionEvent) error {
	// NOTE: Will have discussion for this in future, according to real scenerios
	if event.Status == types.SportEventStatus_STATUS_PENDING || event.Status == types.SportEventStatus_STATUS_INVALID {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "resolution status passed for the sports event is invalid")
	}

	if event.ResolutionTS == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid resolution timestamp for the sport event")
	}

	if !utils.IsValidUID(event.UID) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid uid for the sport event")
	}

	if event.Status == types.SportEventStatus_STATUS_RESULT_DECLARED && len(event.WinnerOddsUIDs) < 1 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "not provided enough winner odds for the sports event")
	}

	for key := range event.WinnerOddsUIDs {
		if !utils.IsValidUID(key) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "odds-uid passed is invalid")
		}
	}

	if event.Status > 4 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid event resolution status ")
	}

	return nil
}
