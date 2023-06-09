package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/market/types"
)

// ResolveMarket accepts ticket containing resolution markets and return response after processing
func (k msgServer) ResolveMarket(
	goCtx context.Context,
	msg *types.MsgResolveMarket,
) (*types.MsgResolveMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var resolutionPayload types.MarketResolutionTicketPayload
	err := k.ovmKeeper.VerifyTicketUnmarshal(goCtx, msg.Ticket, &resolutionPayload)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketVerification, "%s", err)
	}

	if err := resolutionPayload.Validate(); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInTicketPayloadValidation, "%s", err)
	}

	market, found := k.Keeper.GetMarket(ctx, resolutionPayload.UID)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrMarketNotFound, "%s", market.UID)
	}

	if !market.IsResolveAllowed() {
		return nil, sdkerrors.Wrapf(types.ErrMarketResolutionNotAllowed, "%s", market.Status)
	}

	if err := resolutionPayload.ValidateWinnerOdds(&market); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidWinnerOdds, "%s", err)
	}

	resolvedMarket := k.Keeper.ResolveMarket(ctx, market, &resolutionPayload)

	msg.EmitEvent(&ctx, market.UID)

	return &types.MsgResolveMarketResponse{
		Data: resolvedMarket,
	}, nil
}
