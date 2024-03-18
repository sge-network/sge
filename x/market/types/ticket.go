package types

import (
	"strings"

	"github.com/mrz1836/go-sanitize"
	"github.com/spf13/cast"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/utils"
)

// Validate validates market add ticket payload.
func (payload *MarketAddTicketPayload) Validate(ctx sdk.Context) error {
	// remove xss attach prone characters
	payload.Meta = sanitize.XSS(payload.Meta)
	if err := validateMarketTS(ctx, payload.StartTS, payload.EndTS); err != nil {
		return err
	}

	if !(payload.Status == MarketStatus_MARKET_STATUS_ACTIVE ||
		payload.Status == MarketStatus_MARKET_STATUS_INACTIVE) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "acceptable status is active or inactive")
	}

	if !utils.IsValidUID(payload.UID) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid uid for the market")
	}

	if len(payload.Odds) < 2 {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "not provided enough odds for the market")
	}

	if strings.TrimSpace(payload.Meta) == "" {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "meta is mandatory for the market")
	}

	if len(payload.Meta) > MaxAllowedCharactersForMeta {
		return sdkerrors.Wrapf(
			sdkerrtypes.ErrInvalidRequest,
			"meta length should be less than %d characters",
			MaxAllowedCharactersForMeta,
		)
	}

	oddsSet := make(map[string]Odds, len(payload.Odds))
	for _, o := range payload.Odds {
		if o.Meta == "" {
			return sdkerrors.Wrapf(
				sdkerrtypes.ErrInvalidRequest,
				"meta is mandatory for odds with uuid %s",
				o.UID,
			)
		}
		if len(o.Meta) > MaxAllowedCharactersForMeta {
			return sdkerrors.Wrapf(
				sdkerrtypes.ErrInvalidRequest,
				"meta length should be less than %d characters",
				MaxAllowedCharactersForMeta,
			)
		}
		o.Meta = sanitize.XSS(o.Meta)
		if !utils.IsValidUID(o.UID) {
			return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "odds-uid passed is invalid")
		}
		if _, exist := oddsSet[o.UID]; exist {
			return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "duplicate odds-uid in request")
		}
		oddsSet[o.UID] = Odds{}
	}

	return nil
}

// Validate validates market update ticket payload.
func (payload *MarketUpdateTicketPayload) Validate(ctx sdk.Context) error {
	// updating the status to something other than active and inactive
	if !(payload.Status == MarketStatus_MARKET_STATUS_ACTIVE ||
		payload.Status == MarketStatus_MARKET_STATUS_INACTIVE) {
		return sdkerrors.Wrapf(
			sdkerrtypes.ErrInvalidRequest,
			"supported update status is active or inactive",
		)
	}

	return validateMarketTS(ctx, payload.StartTS, payload.EndTS)
}

// Validate validates market resolution ticket payload.
func (payload *MarketResolutionTicketPayload) Validate() error {
	// resolution status should be canceled, aborted or result declared
	if !(payload.Status == MarketStatus_MARKET_STATUS_CANCELED ||
		payload.Status == MarketStatus_MARKET_STATUS_ABORTED ||
		payload.Status == MarketStatus_MARKET_STATUS_RESULT_DECLARED) {
		return sdkerrors.Wrapf(
			sdkerrtypes.ErrInvalidRequest,
			"resolution status passed for the market is invalid",
		)
	}

	switch payload.Status {
	case MarketStatus_MARKET_STATUS_RESULT_DECLARED:
		if len(payload.WinnerOddsUIDs) > maxWinnerUIDs {
			return sdkerrors.Wrapf(
				sdkerrtypes.ErrInvalidRequest,
				"currently only %d winner uid is allowed",
				maxWinnerUIDs,
			)
		}
	default:
		if len(payload.WinnerOddsUIDs) > 0 {
			return sdkerrors.Wrapf(
				sdkerrtypes.ErrInvalidRequest,
				"winner odds should be set if the status is 'result declared'",
			)
		}
	}

	if payload.ResolutionTS == 0 {
		return sdkerrors.Wrapf(
			sdkerrtypes.ErrInvalidRequest,
			"invalid resolution timestamp for the market",
		)
	}

	if !utils.IsValidUID(payload.UID) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid uid for the market")
	}

	if payload.Status == MarketStatus_MARKET_STATUS_RESULT_DECLARED && len(payload.WinnerOddsUIDs) < 1 {
		return sdkerrors.Wrapf(
			sdkerrtypes.ErrInvalidRequest,
			"not provided enough winner odds for the market",
		)
	}

	for _, wid := range payload.WinnerOddsUIDs {
		if !utils.IsValidUID(wid) {
			return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "odds-uid passed is invalid")
		}
	}

	if payload.SgePrice.IsNil() || payload.SgePrice.LTE(sdk.ZeroDec()) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "sge price should be positive")
	}

	return nil
}

// ValidateWinnerOdds validates market resolution ticket payload winner odds.
func (payload *MarketResolutionTicketPayload) ValidateWinnerOdds(market *Market) error {
	if payload.Status == MarketStatus_MARKET_STATUS_RESULT_DECLARED {
		if payload.ResolutionTS < market.StartTS {
			return ErrResolutionTimeLessThenStartTime
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
			return ErrInvalidWinnerOdds
		}
	}

	return nil
}

// validateMarketTS validates start and end timestamp of a market.
func validateMarketTS(ctx sdk.Context, startTS, endTS uint64) error {
	if endTS <= cast.ToUint64(ctx.BlockTime().Unix()) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid end timestamp for the market")
	}

	if startTS >= endTS || startTS == 0 {
		return sdkerrors.Wrapf(
			sdkerrtypes.ErrInvalidRequest,
			"invalid start timestamp for the market, cannot be (greater than eql to EndTs) or 0",
		)
	}

	return nil
}
