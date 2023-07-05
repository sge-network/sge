package types

import (
	"fmt"
)

// Validate performs a basic validation of the LockedBalance fields.
func (lb *LockedBalance) Validate() error {
	if lb.UnlockTime.IsZero() {
		return fmt.Errorf("unlock time is zero")
	}

	if lb.Amount.IsNil() {
		return fmt.Errorf("amount is nil")
	}

	if lb.Amount.IsNegative() {
		return fmt.Errorf("amount is negative")
	}

	return nil
}
