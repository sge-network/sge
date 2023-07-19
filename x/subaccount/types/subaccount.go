package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

// Available reports the coins that are available in the subaccount.
func (m *Balance) Available() sdk.Int {
	return m.DepositedAmount.
		Sub(m.WithdrawmAmount).
		Sub(m.SpentAmount).
		Sub(m.LostAmount)
}
