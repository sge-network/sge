package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/sportevent/types"
)

// KeeperTest is a wrapper object for the keeper, It is being used
// to export unexported methods of the keeper
type KeeperTest = Keeper

func (k KeeperTest) ValidateEventResolution(resolutionPayload types.SportEventResolutionTicketPayload) error {
	return resolutionPayload.Validate()
}

func (k KeeperTest) ValidateEventAdd(ctx sdk.Context, addPayload types.SportEventAddTicketPayload) error {
	params := k.GetParams(ctx)
	return addPayload.Validate(ctx, &params)
}

func (k KeeperTest) ValidateEventUpdate(ctx sdk.Context, updatePayload types.SportEventUpdateTicketPayload, previousEvent types.SportEvent) error {
	params := k.GetParams(ctx)
	return updatePayload.Validate(ctx, &params)
}
