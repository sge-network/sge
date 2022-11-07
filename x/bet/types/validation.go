package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// BetFieldsValidation validates fields of the given bet
func BetFieldsValidation(bet *BetPlaceFields) error {
	if !IsValidUID(bet.UID) {
		return ErrInvalidBetUID
	}

	if bet.Amount.IsNil() || !bet.Amount.IsPositive() {
		return ErrInvalidAmount
	}

	if bet.OddsType < 1 || bet.OddsType > 3 {
		return ErrInvalidOddsType
	}

	if bet.Ticket == "" || strings.Contains(bet.Ticket, " ") {
		return ErrInvalidTicket
	}

	// regexp check of ticket will reside here if requested

	return nil
}

// TicketFieldsValidation validates fields of the given ticketData
func TicketFieldsValidation(ticketData *BetPlacementTicketPayload) error {

	if ticketData.SelectedOdds == nil {
		return ErrOddsDataNotFound
	}

	if ticketData.KycData == nil {
		return ErrNoKycField
	}

	if !IsValidUID(ticketData.SelectedOdds.SportEventUID) {
		return ErrInvalidSportEventUID
	}

	if !IsValidUID(ticketData.SelectedOdds.UID) {
		return ErrInvalidOddsUID
	}

	if ticketData.SelectedOdds.Value == "" {
		return ErrInvalidOddsValue
	}

	ticketOddsValue, err := sdk.NewDecFromStr(ticketData.SelectedOdds.Value)
	if err != nil {
		return sdkerrors.Wrapf(ErrInConvertingOddsToDec, "%s", err)
	}

	if ticketOddsValue.IsNil() || ticketOddsValue.LTE(sdk.OneDec()) {
		return ErrInvalidOddsValue
	}

	if ticketData.KycData.KycRequired && ticketData.KycData.KycId == "" {
		return ErrNoKycIdField
	}

	return nil
}
