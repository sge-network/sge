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

	if bet.Ticket == "" || strings.Contains(bet.Ticket, " ") {
		return ErrInvalidTicket
	}

	// regexp check of ticket will reside here if requested

	return nil
}

// TicketFieldsValidation validates fields of the given ticketData
func TicketFieldsValidation(ticketData *BetPlacementTicketPayload) error {

	for _, o := range ticketData.Odds {

		if !IsValidUID(o.SportEventUID) {
			return ErrInvalidSportEventUID
		}

		if !IsValidUID(o.UID) {
			return ErrInvalidOddsUID
		}

		if o.Value == "" {
			return ErrInvalidOddsValue
		}

		ticketOddsValue, err := sdk.NewDecFromStr(o.Value)
		if err != nil {
			return sdkerrors.Wrapf(ErrInConvertingOddsToDec, "%s", err)
		}

		if ticketOddsValue.IsNil() || ticketOddsValue.LTE(sdk.OneDec()) {
			return ErrInvalidOddsValue
		}

	}

	return nil
}
