package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/market/types"
)

// UpdateMarket accepts ticket containing update market and return response after processing
func (k msgServer) UpdateMarket(goCtx context.Context, msg *types.MsgUpdateMarket) (*types.MsgUpdateMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var updatePayload types.MarketUpdateTicketPayload
	if err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &updatePayload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	storedMarket, found := k.Keeper.GetMarket(ctx, updatePayload.GetUID())
	if !found {
		return nil, types.ErrMarketNotFound
	}

	// if stored market is inactive it is not updatable
	// active status can be changed to inactive or vice versa in the updating
	if storedMarket.Status != types.MarketStatus_MARKET_STATUS_ACTIVE &&
		storedMarket.Status != types.MarketStatus_MARKET_STATUS_INACTIVE {
		return nil, types.ErrMarketCanNotBeAltered
	} //TODO: Something wrong with this check?

	params := k.GetParams(ctx)

	// update market is not valid, return error
	if err := updatePayload.Validate(ctx, &params); err != nil {
		return nil, sdkerrors.Wrap(err, "update validation failed")
	} // TODO: Should these errors not be registered?

	// replace current data with payload values
	// TODO: If any of these values are null because it is not sent in ticket then it will assign null value
	storedMarket.StartTS = updatePayload.StartTS
	storedMarket.EndTS = updatePayload.EndTS
	storedMarket.BetConstraints = params.NewMarketBetConstraints(updatePayload.MinBetAmount, updatePayload.BetFee)
	storedMarket.Status = updatePayload.Status

	// update market is successful, update the module state
	k.Keeper.SetMarket(ctx, storedMarket)

	response := &types.MsgUpdateMarketResponse{
		Data: &storedMarket,
	}
	emitTransactionEvent(ctx, types.TypeMsgUpdateMarkets, response.Data.UID, response.Data.BookUID, msg.Creator)
	return response, nil
}
