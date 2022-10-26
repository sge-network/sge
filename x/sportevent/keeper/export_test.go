package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/sportevent/types"
)

// KeeperTest is a wrapper object for the keeper, It is being used
// to export unexported methods of the keeper
type KeeperTest = Keeper

func (k KeeperTest) ValidateResolutionEvent(event types.ResolutionEvent) error {
	return validateResolutionEvent(event)
}

func (k KeeperTest) MsgServerValidateCreationEvent(ctx sdk.Context, event types.SportEvent) error {
	msgSrv := &msgServer{Keeper: k}
	return msgSrv.validateEventAdd(ctx, &event)
}

func (k KeeperTest) MsgServerValidateUpdateEvent(ctx sdk.Context, event, previousEvent types.SportEvent) error {
	msgSrv := &msgServer{Keeper: k}
	return msgSrv.validateEventUpdate(ctx, event, previousEvent)
}
