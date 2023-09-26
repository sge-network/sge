package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
)

// Validate validates fields of the given ticketData
func (payload *WagerTicketPayload) Validate(creator string) error {
	if payload.SelectedOdds == nil {
		return ErrOddsDataNotFound
	}

	if err := payload.ValidateOdds(*payload.SelectedOdds); err != nil {
		return sdkerrors.Wrapf(err, "%s", payload.SelectedOdds.UID)
	}

	for _, odd := range payload.Odds {
		if err := payload.ValidateOdds(*odd); err != nil {
			return sdkerrors.Wrapf(err, "%s", odd.UID)
		}
	}

	if !payload.KycData.Validate(creator) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", creator)
	}

	if payload.OddsType < OddsType_ODDS_TYPE_DECIMAL ||
		payload.OddsType > OddsType_ODDS_TYPE_MONEYLINE {
		return ErrInvalidOddsType
	}

	return nil
}

func (payload *WagerTicketPayload) ValidateOdds(odds BetOdds) error {
	if !utils.IsValidUID(odds.MarketUID) {
		return ErrInvalidMarketUID
	}

	if !utils.IsValidUID(odds.UID) {
		return ErrInvalidOddsUID
	}

	if len(strings.TrimSpace(odds.Value)) == 0 {
		return ErrEmptyOddsValue
	}

	if odds.MaxLossMultiplier.IsNil() || odds.MaxLossMultiplier.LTE(sdk.ZeroDec()) {
		return ErrMaxLossMultiplierCanNotBeZero
	}

	if odds.MaxLossMultiplier.GT(sdk.OneDec()) {
		return ErrMaxLossMultiplierCanNotBeMoreThanOne
	}

	return nil
}

func (payload *WagerTicketPayload) OddsMap() map[string]*BetOdds {
	oddMap := make(map[string]*BetOdds)
	for _, odd := range payload.Odds {
		oddMap[odd.UID] = odd
	}
	return oddMap
}
