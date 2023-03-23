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
	err := k.orderBookKeeper.InitiateBook(ctx, addPayload.UID, addPayload.SrContributionForHouse, oddsUIDs)
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
