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
func (as *AccountSummary) Available() sdkmath.Int {
	return as.DepositedAmount.
		Sub(as.WithdrawnAmount).
		Sub(as.SpentAmount).
		Sub(as.LostAmount)
}

// Spend modifies the spent amount of subaccount balance according to the spent value.
func (as *AccountSummary) Spend(amt sdkmath.Int) error {
	if amt.IsNegative() {
		return fmt.Errorf("amount is not positive")
	}
	if amt.GT(as.Available()) {
		return fmt.Errorf("amount is greater than available")
	}
	as.SpentAmount = as.SpentAmount.Add(amt)
	return nil
}

// Unspend modifies the spent amount of subaccount balance according to the undpent value.
func (as *AccountSummary) Unspend(amt sdkmath.Int) error {
	if amt.IsNegative() {
		return fmt.Errorf("amount is not positive")
	}
	if amt.GT(as.SpentAmount) {
		return fmt.Errorf("amount is greater than spent")
	}
	as.SpentAmount = as.SpentAmount.Sub(amt)
	return nil
}

// AddLoss adds to the lost amout of subaccount balance after losing the bet.
func (as *AccountSummary) AddLoss(amt sdkmath.Int) error {
	if amt.IsNegative() {
		return fmt.Errorf("amount is not positive")
	}
	as.LostAmount = as.LostAmount.Add(amt)
	return nil
}

// Withdraw sends deposited amount to withdrawn
func (as *AccountSummary) Withdraw(amt sdkmath.Int) error {
	if amt.IsNegative() {
		return fmt.Errorf("amount is not positive")
	}
	if amt.GT(as.Available()) {
		return fmt.Errorf("amount is greater than available")
	}
	as.WithdrawnAmount = as.WithdrawnAmount.Add(amt)

	return nil
}

// WithdrawableBalance returns withdrawable balance of a subaccount
func (as *AccountSummary) WithdrawableBalance(unlockedBalance, bankBalance sdkmath.Int) sdkmath.Int {
	// calculate withdrawable balance, which is the minimum between the available balance, and
	// what has been unlocked so far. Also, it cannot be greater than the bank balance.
	// Available reports the deposited amount - spent amount - lost amount - withdrawn amount.
	return sdkmath.MinInt(sdkmath.MinInt(as.Available(), unlockedBalance), bankBalance)
}
