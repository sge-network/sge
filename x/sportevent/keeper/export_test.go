package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/sportevent/types"
)

// KeeperTest is a wrapper object for the keeper, It is being used
// to export unexported methods of the keeper
type KeeperTest = Keeper

func (k KeeperTest) ValidateEventResolution(event types.SportEventResolutionTicketPayload) error {
	return validateEventResolution(event)
}

func (k KeeperTest) MsgServerValidateCreationEvent(ctx sdk.Context, event types.SportEvent) error {
	msgSrv := &msgServer{Keeper: k}
	return msgSrv.validateAddEvent(ctx, &event)
}

func (k KeeperTest) MsgServerValidateEventUpdate(ctx sdk.Context, updatePayload types.SportEventUpdateTicketPayload, previousEvent types.SportEvent) error {
	msgSrv := &msgServer{Keeper: k}
	return msgSrv.validateEventUpdate(ctx, updatePayload, previousEvent)
}
