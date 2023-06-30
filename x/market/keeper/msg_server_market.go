package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/market/types"
)

// AddMarket accepts ticket containing creation market and return response after processing
func (k msgServer) AddMarket(
	goCtx context.Context,
	msg *types.MsgAddMarket,
) (*types.MsgAddMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var addPayload types.MarketAddTicketPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &addPayload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if err := addPayload.Validate(ctx); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	_, found := k.Keeper.GetMarket(ctx, addPayload.UID)
	if found {
		return nil, types.ErrMarketAlreadyExist
	}

	var oddsUIDs []string
	for _, odds := range addPayload.Odds {
		oddsUIDs = append(oddsUIDs, odds.UID)
	}
	err := k.orderbookKeeper.InitiateOrderBook(ctx, addPayload.UID, oddsUIDs)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInOrderBookInitiation, "%s", err)
	}

	market := types.NewMarket(
		addPayload.UID,
		msg.Creator,
		addPayload.StartTS,
		addPayload.EndTS,
		addPayload.Odds,
		addPayload.Meta,
		addPayload.UID,
		addPayload.Status,
	)

	k.Keeper.SetMarket(ctx, market)

	msg.EmitEvent(&ctx, market.UID, market.BookUID)

	return &types.MsgAddMarketResponse{
		Error: "",
		Data:  &market,
	}, nil
}

// UpdateMarket accepts ticket containing update market and return response after processing
func (k msgServer) UpdateMarket(
	goCtx context.Context,
	msg *types.MsgUpdateMarket,
) (*types.MsgUpdateMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var updatePayload types.MarketUpdateTicketPayload
	if err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &updatePayload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	market, found := k.Keeper.GetMarket(ctx, updatePayload.GetUID())
	if !found {
		return nil, types.ErrMarketNotFound
	}

	// if stored market is inactive it is not updatable
	// active status can be changed to inactive or vice versa in the updating
	if !market.IsUpdateAllowed() {
		return nil, sdkerrors.Wrapf(types.ErrMarketCanNotBeAltered, "%s", market.Status)
	}

	// update market is not valid, return error
	if err := updatePayload.Validate(ctx); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	// replace current data with payload values
	market.StartTS = updatePayload.StartTS
	market.EndTS = updatePayload.EndTS
	market.Status = updatePayload.Status

	// update market is successful, update the module state
	k.Keeper.SetMarket(ctx, market)

	msg.EmitEvent(&ctx, market.UID)

	return &types.MsgUpdateMarketResponse{Data: &market}, nil
}
