package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/market/types"
)

// Resolve accepts ticket containing resolution markets and return response after processing
func (k msgServer) Resolve(
	goCtx context.Context,
	msg *types.MsgResolve,
) (*types.MsgResolveResponse, error) {
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

	if err := k.Keeper.TopUpPriceLockPool(ctx, market.Creator, market.PriceStats.PriceReimbursement); err != nil {
		return nil, err
	}

	resolvedMarket := k.Keeper.Resolve(ctx, market, &resolutionPayload)

	msg.EmitEvent(&ctx, market.UID)

	return &types.MsgResolveResponse{
		Data: resolvedMarket,
	}, nil
}
