package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/market/types"
)

// ResolveMarket accepts ticket containing resolution markets and return response after processing
func (k msgServer) ResolveMarket(goCtx context.Context, msg *types.MsgResolveMarket) (*types.MsgResolveMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var resolutionPayload types.MarketResolutionTicketPayload
	err := k.dvmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &resolutionPayload)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInVerification, "%s", err)
	}

	if err := resolutionPayload.Validate(); err != nil {
		return nil, sdkerrors.Wrap(err, "validate resolution data")
	}

	market, found := k.Keeper.GetMarket(ctx, resolutionPayload.UID)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrMarketNotFound, "getting market")
	}

	if !market.IsResolveAllowed() {
		return nil, sdkerrors.Wrap(types.ErrMarketIsNotActiveOrInactive, "getting market")
	}

	if err := resolutionPayload.ValidateWinnerOdds(&market); err != nil {
		return nil, sdkerrors.Wrap(err, "extract winner odds id")
	}

	resolvedMarket := k.Keeper.ResolveMarket(ctx, market, &resolutionPayload)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "resolve market")
	}

	emitTransactionEvent(ctx, types.TypeMsgResolveMarkets, resolvedMarket.UID, resolvedMarket.BookUID, msg.Creator)

	return &types.MsgResolveMarketResponse{
		Data: resolvedMarket,
	}, nil
}
