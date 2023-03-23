package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/market/types"
)

// AddMarket accepts ticket containing creation market and return response after processing
func (k msgServer) AddMarket(goCtx context.Context, msg *types.MsgAddMarket) (*types.MsgAddMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var addPayload types.MarketAddTicketPayload
	if err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &addPayload); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	params := k.GetParams(ctx)

	if err := addPayload.Validate(ctx, &params); err != nil {
		return nil, sdkerrors.Wrap(err, "validate add market")
	}

	_, found := k.Keeper.GetMarket(ctx, addPayload.UID)
	if found {
		return nil, types.ErrMarketAlreadyExist
	}

	var oddsUIDs []string
	for _, odds := range addPayload.Odds {
		oddsUIDs = append(oddsUIDs, odds.UID)
	}
	err := k.srKeeper.InitiateBook(ctx, addPayload.UID, addPayload.SrContributionForHouse, oddsUIDs)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInOrderBookInitiation, "%s", err)
	}

	market := types.NewMarket(
		addPayload.UID,
		msg.Creator,
		addPayload.StartTS,
		addPayload.EndTS,
		addPayload.Odds,
		params.NewMarketBetConstraints(addPayload.MinBetAmount, addPayload.BetFee),
		addPayload.Meta,
		addPayload.UID,
		addPayload.SrContributionForHouse,
		addPayload.Status,
	)

	k.Keeper.SetMarket(ctx, market)

	response := &types.MsgAddMarketResponse{
		Error: "",
		Data:  &market,
	}
	emitTransactionEvent(ctx, types.TypeMsgCreateMarkets, response.Data.UID, addPayload.UID, msg.Creator)

	return response, nil
}

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
	if !storedMarket.IsUpdateAllowed() {
		return nil, types.ErrMarketCanNotBeAltered
	}

	params := k.GetParams(ctx)

	// update market is not valid, return error
	if err := updatePayload.Validate(ctx, &params); err != nil {
		return nil, sdkerrors.Wrap(err, "update validation failed")
	} // TODO: Should these errors not be registered?

	// replace current data with payload values
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
