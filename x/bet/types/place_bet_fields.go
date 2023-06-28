package types

import (
	"strings"

	"github.com/sge-network/sge/utils"
)

// BetFieldsValidation validates fields of the given bet
func BetFieldsValidation(bet *PlaceBetFields) error {
	if !utils.IsValidUID(bet.UID) {
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
