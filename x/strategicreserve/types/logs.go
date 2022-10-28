package types

const (
	// ─── ERROR LOGS ─────────────────────────────────────────────────────────

	// LogErrInsufficientUserBalance is logged when the user has insufficient balance
	LogErrInsufficientUserBalance = "User has insufficient balance. Address: %s, Balance: %s usge, Required: %s usge"

	// LogErrInsufficientUnlockedAmountInSrPool is logged when SR_Pool unlocked amount is
	// less than the required amount
	LogErrInsufficientUnlockedAmountInSrPool = "Unlocked amount in the SR_Pool is insufficient. Unlocked amount: %s usge, Required amount: %s usge"

	// LogErrInsufficientLockedAmountInSrPool is logged when SR_Pool locked amount is
	// less than the required amount
	LogErrInsufficientLockedAmountInSrPool = "Locked amount in the SR_Pool is insufficient. Locked amount: %s usge, Required amount: %s usge"

	// LogErrFromBankModule is logged when an error is returned from the bank module
	// while trying to transfer funds
	LogErrFromBankModule = "Funds transfer failed from the bank module with error: %s"

	// LogErrPayoutLockDoesnotExist is logged when the payout lock doesn't exist
	LogErrPayoutLockDoesnotExist = "lock %s for payout does not exist"

	// LogErrLockAlreadyExists is logged when locks already exists
	LogErrLockAlreadyExists = "lock already exists for the lock: %s"

	// LogErrTransferOfFundsFailed is logged when funds transfer gets failed
	LogErrTransferOfFundsFailed = "Transfer of %s usge from %s to %s failed with error: %s"

	// ─── INFO LOGS ─────────────────────────────────────────────────────────

	// LogInfoFundsTransferred is logged when the funds are transferred
	// to the required account successfully
	LogInfoFundsTransferred = "transferred %s usge from %s to %s"

	// LogInfoBettorReceivedPayout is logged when the user receives the payout successfully
	LogInfoBettorReceivedPayout = "bettor %s received the payout of %s usge"

	// LogInfoHouseReceivedWinnings is logged when the house receives the bet
	// winnings successfully
	LogInfoHouseReceivedWinnings = "house received %s usge as winnings in %s account"

	// LogInfoBetAccepted is logged when a bet is accepted
	LogInfoBetAccepted = "bet with bet amount %s usge accepted"

	// LogInfoBettorRefunded is logged when a bettor is refunded with the bet amount
	LogInfoBettorRefunded = "bettor %s refunded with bet amount %s usge"
)
