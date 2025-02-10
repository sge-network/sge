package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"github.com/sge-network/sge/utils"
)

// Validate validates fields of the given ticketData
func (payload *WagerTicketPayload) Validate(creator string) error {
	if payload.SelectedOdds == nil {
		return ErrOddsDataNotFound
	}

	if OddsType_ODDS_TYPE_UNSPECIFIED > payload.Meta.SelectedOddsType ||
		OddsType_ODDS_TYPE_MONEYLINE < payload.Meta.SelectedOddsType {
		return sdkerrors.Wrapf(ErrMetaOddsType, "%s", payload.Meta.SelectedOddsType)
	}

	if err := payload.ValidateOdds(); err != nil {
		return sdkerrors.Wrapf(err, "%s", payload.SelectedOdds.UID)
	}

	for _, odd := range payload.AllOdds {
		if err := payload.ValidateCompactOdds(*odd); err != nil {
			return sdkerrors.Wrapf(err, "%s", odd.UID)
		}
	}

	if !payload.KycData.Validate(creator) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", creator)
	}

	return nil
}

func (payload *WagerTicketPayload) ValidateOdds() error {
	if !utils.IsValidUID(payload.SelectedOdds.MarketUID) {
		return ErrInvalidMarketUID
	}

	if !utils.IsValidUID(payload.SelectedOdds.UID) {
		return ErrInvalidOddsUID
	}

	if len(strings.TrimSpace(payload.SelectedOdds.Value)) == 0 {
		return ErrEmptyOddsValue
	}

	if payload.SelectedOdds.MaxLossMultiplier.IsNil() || payload.SelectedOdds.MaxLossMultiplier.LTE(sdkmath.LegacyZeroDec()) {
		return ErrMaxLossMultiplierCanNotBeZero
	}

	if payload.SelectedOdds.MaxLossMultiplier.GT(sdkmath.LegacyOneDec()) {
		return ErrMaxLossMultiplierCanNotBeMoreThanOne
	}

	return nil
}

func (payload *WagerTicketPayload) ValidateCompactOdds(odds BetOddsCompact) error {
	if !utils.IsValidUID(odds.UID) {
		return ErrInvalidOddsUID
	}

	if odds.MaxLossMultiplier.IsNil() || odds.MaxLossMultiplier.LTE(sdkmath.LegacyZeroDec()) {
		return ErrMaxLossMultiplierCanNotBeZero
	}

	if odds.MaxLossMultiplier.GT(sdkmath.LegacyOneDec()) {
		return ErrMaxLossMultiplierCanNotBeMoreThanOne
	}

	return nil
}

func (payload *WagerTicketPayload) OddsMap() map[string]*BetOddsCompact {
	oddMap := make(map[string]*BetOddsCompact)
	for _, odd := range payload.AllOdds {
		oddMap[odd.UID] = odd
	}
	return oddMap
}
