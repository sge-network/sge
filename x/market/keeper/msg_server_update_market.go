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

	currentData, found := k.Keeper.GetMarket(ctx, updatePayload.GetUID())
	if !found {
		return nil, types.ErrEventNotFound
	}

	// if current data is not active or inactive it is not updatable
	// active status can be changed to inactive or vice versa in the updating
	if currentData.Status != types.MarketStatus_MARKET_STATUS_ACTIVE &&
		currentData.Status != types.MarketStatus_MARKET_STATUS_INACTIVE {
		return nil, types.ErrCanNotBeAltered
	}

	// update market is not valid, return error
	params := k.GetParams(ctx)

	if err := updatePayload.Validate(ctx, &params); err != nil {
		return nil, sdkerrors.Wrap(err, "validate update data")
	}

	// replace current data with payload values
	currentData.StartTS = updatePayload.StartTS
	currentData.EndTS = updatePayload.EndTS
	currentData.BetConstraints = params.NewMarketBetConstraints(updatePayload.MinBetAmount, updatePayload.BetFee)
	currentData.Status = updatePayload.Status

	// update market is successful, update the module state
	k.Keeper.SetMarket(ctx, currentData)

	response := &types.MsgUpdateMarketResponse{
		Data: &currentData,
	}
	emitTransactionEvent(ctx, types.TypeMsgUpdateMarkets, response.Data.UID, response.Data.BookUID, msg.Creator)
	return response, nil
}
