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

	resolvedMarket, err := k.processMarketResolution(ctx, &resolutionPayload)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "process resolution market")
		// TODO: Should we not register these error messages?
	}

	response := &types.MsgResolveMarketResponse{
		Data: resolvedMarket,
	}
	emitTransactionEvent(ctx, types.TypeMsgResolveMarkets, response.Data.UID, response.Data.BookUID, msg.Creator)
	return response, nil
}

func (k msgServer) processMarketResolution(ctx sdk.Context, resolutionPayload *types.MarketResolutionTicketPayload) (*types.Market, error) {
	if err := resolutionPayload.Validate(); err != nil {
		return nil, sdkerrors.Wrap(err, "validate resolution data")
	}

	market, err := k.getMarketToResolve(ctx, *resolutionPayload)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "getting market")
	}

	if err := extractWinnerOddsUIDs(&market, resolutionPayload); err != nil {
		return nil, sdkerrors.Wrap(err, "extract winner odds id")
	}

	resolvedMarket, err := k.Keeper.ResolveMarket(ctx, resolutionPayload)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "resolve market")
	}

	return resolvedMarket, nil
}

func (k msgServer) getMarketToResolve(ctx sdk.Context, resolutionPayload types.MarketResolutionTicketPayload) (types.Market, error) {
	market, found := k.Keeper.GetMarket(ctx, resolutionPayload.UID)
	if !found {
		return types.Market{}, types.ErrMarketNotFound
	}

	if market.Status != types.MarketStatus_MARKET_STATUS_ACTIVE &&
		market.Status != types.MarketStatus_MARKET_STATUS_INACTIVE {
		return types.Market{}, types.ErrMarketIsNotActiveOrInactive
	}
	//TODO: this is repeated check, I found similar in keeper.ResolveMarket

	return market, nil
}

func extractWinnerOddsUIDs(market *types.Market, payload *types.MarketResolutionTicketPayload) error {
	if payload.Status == types.MarketStatus_MARKET_STATUS_RESULT_DECLARED {
		if payload.ResolutionTS < market.StartTS {
			return types.ErrResolutionTimeLessThenStartTime
		}

		validWinnerOdds := true
		for _, wid := range payload.WinnerOddsUIDs {
			validWinnerOdds = false
			for _, o := range market.Odds {
				if o.UID == wid {
					validWinnerOdds = true
				}
			}
			if !validWinnerOdds {
				break
			}
		}

		if !validWinnerOdds {
			return types.ErrInvalidWinnerOdds
		}
	}

	return nil
}
