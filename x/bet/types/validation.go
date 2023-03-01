package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// BetFieldsValidation validates fields of the given bet
func BetFieldsValidation(bet *PlaceBetFields) error {
	if !IsValidUID(bet.UID) {
		return ErrInvalidBetUID
	}

	if bet.Amount.IsNil() || !bet.Amount.IsPositive() {
		return ErrInvalidAmount
	}

	if bet.Ticket == "" || strings.Contains(bet.Ticket, " ") {
		return ErrInvalidTicket
	}

	return nil
}

// Validate validates fields of the given ticketData
func (ticketData *BetPlacementTicketPayload) Validate(creator string) error {
	if ticketData.SelectedOdds == nil {
		return ErrOddsDataNotFound
	}

	if !ticketData.KycData.validate(creator) {
		return sdkerrors.Wrapf(ErrUserKycFailed, "%s", creator)
	}

	if !IsValidUID(ticketData.SelectedOdds.SportEventUID) {
		return ErrInvalidSportEventUID
	}

	if !IsValidUID(ticketData.SelectedOdds.UID) {
		return ErrInvalidOddsUID
	}

	if len(strings.TrimSpace(ticketData.SelectedOdds.Value)) == 0 {
		return ErrEmptyOddsValue
	}

	if ticketData.SelectedOdds.MaxLossMultiplier.IsNil() || ticketData.SelectedOdds.MaxLossMultiplier.LTE(sdk.ZeroDec()) {
		return ErrMaxLossMultiplierCanNotBeZero
	}

	if ticketData.SelectedOdds.MaxLossMultiplier.GT(sdk.OneDec()) {
		return ErrMaxLossMultiplierCanNotBeMoreThanOne
	}

	if ticketData.OddsType < OddsType_ODDS_TYPE_DECIMAL ||
		ticketData.OddsType > OddsType_ODDS_TYPE_MONEYLINE {
		return ErrInvalidOddsType
	}

	return nil
}

// validate checks whether the kyc data is valid for a particular bettor
// if the kyc is required
func (payload KycDataPayload) validate(address string) bool {
	// ignore is true means that kyc check should be ignored
	if payload.Ignore {
		return true
	}

	if payload.Approved && payload.ID == address {
		return true
	}

	return false
}
