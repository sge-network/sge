package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BetFieldsValidation validates fields of the given bet
func BetFieldsValidation(bet *PlaceBetFields) error {
	if !IsValidUID(bet.UID) {
		return ErrInvalidBetUID
	}

	if bet.Amount.IsNil() || !bet.Amount.IsPositive() {
		return ErrInvalidAmount
	}

	if bet.OddsType < OddsType_ODDS_TYPE_DECIMAL ||
		bet.OddsType > OddsType_ODDS_TYPE_MONEYLINE {
		return ErrInvalidOddsType
	}

	if bet.Ticket == "" || strings.Contains(bet.Ticket, " ") {
		return ErrInvalidTicket
	}

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

	if len(strings.TrimSpace(ticketData.SelectedOdds.Value)) == 0 {
		return ErrEmptyOddsValue
	}

	if ticketData.SelectedOdds.MaxLossMultiplier.IsNil() || ticketData.SelectedOdds.MaxLossMultiplier.LTE(sdk.ZeroDec()) {
		return ErrMaxLossMultiplierCanNotBeZero
	}

	if ticketData.KycData.KycRequired && ticketData.KycData.KycId == "" {
		return ErrNoKycIDField
	}

	return nil
}
