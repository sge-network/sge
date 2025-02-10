package types

import (
	"strings"

	"github.com/sge-network/sge/utils"
)

// WagerValidation validates fields of the given bet
func WagerValidation(props *WagerProps) error {
	if !utils.IsValidUID(props.UID) {
		return ErrInvalidBetUID
	}

	if props.Amount.IsNil() || !props.Amount.IsPositive() {
		return ErrInvalidAmount
	}

	if props.Ticket == "" || strings.Contains(props.Ticket, " ") {
		return ErrInvalidTicket
	}

	return nil
}
