package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

// Validate performs a basic validation of the LockedBalance fields.
func (lb *LockedBalance) Validate() error {
	if lb.UnlockTS == 0 {
		return fmt.Errorf("unlock time is zero %d", lb.UnlockTS)
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
func (m *AccountSummary) Available() sdkmath.Int {
	return m.DepositedAmount.
		Sub(m.WithdrawnAmount).
		Sub(m.SpentAmount).
		Sub(m.LostAmount)
}

// Spend modifies the spent amount of subaccount balance according to the spent value.
func (m *AccountSummary) Spend(amt sdkmath.Int) error {
	if amt.IsNegative() {
		return fmt.Errorf("amount is not positive")
	}
	if amt.GT(m.Available()) {
		return fmt.Errorf("amount is greater than available")
	}
	m.SpentAmount = m.SpentAmount.Add(amt)
	return nil
}

// Unspend modifies the spent amount of subaccount balance according to the undpent value.
func (m *AccountSummary) Unspend(amt sdkmath.Int) error {
	if amt.IsNegative() {
		return fmt.Errorf("amount is not positive")
	}
	if amt.GT(m.SpentAmount) {
		return fmt.Errorf("amount is greater than spent")
	}
	m.SpentAmount = m.SpentAmount.Sub(amt)
	return nil
}

// AddLoss adds to the lost amout of subaccount balance after losing the bet.
func (m *AccountSummary) AddLoss(amt sdkmath.Int) error {
	if amt.IsNegative() {
		return fmt.Errorf("amount is not positive")
	}
	m.LostAmount = m.LostAmount.Add(amt)
	return nil
}
