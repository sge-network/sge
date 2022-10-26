package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/app/params"
	bet "github.com/sge-network/sge/x/bet/types"
	"github.com/sge-network/sge/x/strategicreserve/types"
)

// ProcessBetPlacement transfers the bet fee from the bettor's account
// to the bet module account and the bet amount from the given bettor's
// account to the `bet_reserve` account of SR and locks the extra
// payout in the `sr_pool` account of SR.
// payout = bet amount * odds value
// extra payout = payout - bet amount
func (k Keeper) ProcessBetPlacement(
	ctx sdk.Context,
	bettorAddress sdk.AccAddress,
	betFee sdk.Coin,
	betAmount sdk.Int,
	extraPayout sdk.Int,
	uniqueLock string,
	endTs uint64,
) error {

	// If lock exists, return error
	// Lock already exists means the bet is already placed for the given bet-uid.
	if k.payoutLockExists(ctx, uniqueLock) {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrLockAlreadyExists, uniqueLock))
		return types.ErrLockAlreadyExists
	}

	// For a valid `endTs`, `dayID` can never be negative
	dayID := getDayID(endTs)
	if dayID.IsNegative() {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrInvalidEndTimestamp, endTs))
		return types.ErrInvalidEndTimestamp
	}

	reserver := k.GetReserver(ctx)

	// If SR unlocked amount has insufficient balance, return error
	// NOTE: This check may result in completely emptying the SR pool.
	// replacement of SR pool capacity calculation to minimise risking
	// the SR pool due to payout. Risk management should also be done here if requested.
	if reserver.SrPool.UnlockedAmount.LT(extraPayout) {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrInsufficientUnlockedAmountInSrPool,
			reserver.SrPool.UnlockedAmount, extraPayout))
		return types.ErrInsufficientUnlockedAmountInSrPool
	}

	// Transfer bet fee from bettor to the `bet` module account
	err := k.transferFundsFromUserToModule(ctx, bettorAddress, bet.ModuleName,
		betFee.Amount)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrTransferOfFundsFailed, betFee,
			bettorAddress, bet.ModuleName, err.Error()))
		return err
	}

	// Transfer bet amount from bettor to `bet_reserve` Account
	err = k.transferFundsFromUserToModule(ctx, bettorAddress, types.BetReserveName,
		betAmount)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrTransferOfFundsFailed, betAmount,
			bettorAddress, types.BetReserveName, err.Error()))
		return err
	}

	k.updateSrPool(ctx, reserver.SrPool.LockedAmount.Add(extraPayout),
		reserver.SrPool.UnlockedAmount.Sub(extraPayout))

	// Create a unique lock in the Payout Store for the bet
	k.setPayoutLock(ctx, uniqueLock)

	k.Logger(ctx).Info(fmt.Sprintf(types.LogInfoBetAccepted, betAmount.String()))
	return nil
}

// BettorWins pays the payout to the bettor from the `bet_reserve`
// of the SR. Also, it transfers the extra payout locked in the
// `sr_pool` to the `bet_reserve`. It should be called when the
// bettor wins the bet.
// payout = bet amount * odds value
// extra payout = payout - bet amount
func (k Keeper) BettorWins(
	ctx sdk.Context,
	bettorAddress sdk.AccAddress,
	betAmount sdk.Int,
	extraPayout sdk.Int,
	uniqueLock string,
) error {

	// Idempotency check: If lock does not exist, return error
	if !k.payoutLockExists(ctx, uniqueLock) {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrPayoutLockDoesnotExist, uniqueLock))
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	reserver := k.GetReserver(ctx)

	// If SR locked amount has insufficient balance, return error
	if reserver.SrPool.LockedAmount.LT(extraPayout) {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrInsufficientLockedAmountInSrPool,
			reserver.SrPool.LockedAmount, extraPayout))
		return types.ErrInsufficientLockedAmountInSrPool
	}

	// Transfer extra payout from `sr_pool` to `bet_reserve` account
	err := k.transferFundsFromModuleToModule(ctx, types.SRPoolName,
		types.BetReserveName, extraPayout)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrTransferOfFundsFailed, extraPayout,
			types.SRPoolName, types.BetReserveName, err))
		return err
	}

	k.updateSrPool(ctx, reserver.SrPool.LockedAmount.Sub(extraPayout),
		reserver.SrPool.UnlockedAmount)

	payout := betAmount.Add(extraPayout)

	// Transfer payout from the `bet_reserve` account to bettor
	err = k.transferFundsFromModuleToUser(ctx, types.BetReserveName, bettorAddress,
		payout)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrTransferOfFundsFailed, payout,
			types.BetReserveName, bettorAddress, err))
		return err
	}

	// Delete lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)

	k.Logger(ctx).Info(fmt.Sprintf(types.LogInfoBettorReceivedPayout, bettorAddress,
		payout))
	return nil
}

// BettorLoses unlocks the extra payout in the `sr_pool`. It transfers the
// bet amount (house winnings) from the `bet_reserve` to the `winnings_collector`
// module account of SR. It should be called when the bettor loses the bet.
// payout = bet amount * odds value
// extra payout = payout - bet amount
func (k Keeper) BettorLoses(ctx sdk.Context, address sdk.AccAddress,
	betAmount sdk.Int, extraPayout sdk.Int, uniqueLock string) error {

	// Idempotency check: If lock does not exist, return error
	if !k.payoutLockExists(ctx, uniqueLock) {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrPayoutLockDoesnotExist, uniqueLock))
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	// Transfer bet amount from `bet_reserve` to `winnings_collector` module
	// account of the SR
	err := k.transferFundsFromModuleToModule(ctx, types.BetReserveName,
		types.WinningsCollectorName, betAmount)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrTransferOfFundsFailed, betAmount,
			types.BetReserveName, types.WinningsCollectorName, err))
		return err
	}

	reserver := k.GetReserver(ctx)
	k.updateSrPool(ctx, reserver.SrPool.LockedAmount.Sub(extraPayout),
		reserver.SrPool.UnlockedAmount.Add(extraPayout))

	// Delete lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)

	k.Logger(ctx).Info(fmt.Sprintf(types.LogInfoHouseReceivedWinnings, betAmount,
		types.WinningsCollectorName))
	return nil
}

// RefundBettor refunds back the bet amount from the `bet_reserve` to
// the bettor in case a sports event gets cancelled or aborted.
// It should be called when a sports event is cancelled or aborted
// and the bet amount needs to be refunded back to the bettor.
// payout = bet amount * odds value
// extra payout = payout - bet amount
func (k Keeper) RefundBettor(ctx sdk.Context, bettorAddress sdk.AccAddress,
	betAmount sdk.Int, extraPayout sdk.Int, uniqueLock string) error {

	// Idempotency check: If lock does not exist, return error
	if !k.payoutLockExists(ctx, uniqueLock) {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrPayoutLockDoesnotExist, uniqueLock))
		return sdkerrors.Wrapf(types.ErrPayoutLockDoesnotExist, uniqueLock)
	}

	reserver := k.GetReserver(ctx)

	// If SR locked amount has insufficient balance, return error
	if reserver.SrPool.LockedAmount.LT(extraPayout) {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrInsufficientLockedAmountInSrPool,
			reserver.SrPool.LockedAmount, extraPayout))
		return types.ErrInsufficientLockedAmountInSrPool
	}

	// Transfer bet amount from `bet_reserve` to bettor's account
	err := k.transferFundsFromModuleToUser(ctx, types.BetReserveName, bettorAddress,
		betAmount)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrTransferOfFundsFailed, betAmount,
			types.BetReserveName, bettorAddress, err))
		return err
	}

	k.updateSrPool(ctx, reserver.SrPool.LockedAmount.Sub(extraPayout),
		reserver.SrPool.UnlockedAmount.Add(extraPayout))

	// Delete the lock from the payout store as the bet is settled
	k.removePayoutLock(ctx, uniqueLock)

	k.Logger(ctx).Info(fmt.Sprintf(types.LogInfoBettorRefunded, bettorAddress,
		betAmount))
	return nil
}

// updateSrPool updates the Reserver.SrPool with the new amounts
func (k Keeper) updateSrPool(ctx sdk.Context, newLockedAmount sdk.Int,
	newUnlockedAmount sdk.Int) {

	reserver := k.GetReserver(ctx)

	// Update the Reserver.SrPool
	reserver.SrPool.LockedAmount = newLockedAmount
	reserver.SrPool.UnlockedAmount = newUnlockedAmount

	// Set the updated reserver
	k.SetReserver(ctx, reserver)
}

// transferFundsFromUserToModule transfers the given amount from
// the given account address to the module account passed.
// Returns an error if the account holder has insufficient balance.
func (k Keeper) transferFundsFromUserToModule(ctx sdk.Context,
	address sdk.AccAddress, moduleAccName string, amount sdk.Int) error {

	// Get the spendable balance of the account holder
	usgeCoins := k.bankKeeper.SpendableCoins(ctx,
		address).AmountOf(params.DefaultBondDenom)

	// If account holder has insufficient balance, return error
	if usgeCoins.LT(amount) {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrInsufficientUserBalance,
			address, usgeCoins, amount))
		return sdkerrors.Wrapf(types.ErrInsufficientUserBalance, address.String())
	}

	amt := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, amount))

	// Transfer funds
	err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, address, moduleAccName, amt)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrFromBankModule, err))
		return sdkerrors.Wrapf(types.ErrFromBankModule, err.Error())
	}

	k.Logger(ctx).Info(fmt.Sprintf(types.LogInfoFundsTransferred, amount,
		address, moduleAccName))
	return nil
}

// transferFundsFromModuleToUser transfers the given amount from a module
// account to the given account address.
// Returns an error if the account holder has insufficient balance.
func (k Keeper) transferFundsFromModuleToUser(ctx sdk.Context,
	moduleAccName string, address sdk.AccAddress, amount sdk.Int) error {

	// Get the balance of the sender module account
	balance := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(
		moduleAccName), params.DefaultBondDenom)

	amt := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, amount))

	// If module account has insufficient balance, return error
	if balance.Amount.LT(amt.AmountOf(params.DefaultBondDenom)) {
		return sdkerrors.Wrapf(types.ErrInsufficientBalanceInModuleAccount,
			moduleAccName)
	}

	// Transfer funds
	err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAccName,
		address, amt)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrFromBankModule, err))
		return sdkerrors.Wrapf(types.ErrFromBankModule, err.Error())
	}

	k.Logger(ctx).Info(fmt.Sprintf(types.LogInfoFundsTransferred, amount,
		moduleAccName, address))
	return nil
}

// transferFundsFromModuleToModule transfers the given amount from a module
// account to another module account.
// Returns an error if the sender module has insufficient balance.
func (k Keeper) transferFundsFromModuleToModule(ctx sdk.Context,
	senderModule string, recipientModule string, amount sdk.Int) error {

	if senderModule == recipientModule {
		return types.ErrDuplicateSenderAndRecipientModule
	}

	amt := sdk.NewCoins(sdk.NewCoin(params.DefaultBondDenom, amount))

	// Get the balance of the sender module account
	balance := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(
		senderModule), params.DefaultBondDenom)

	// If module account has insufficient balance, return error
	if balance.Amount.LT(amt.AmountOf(params.DefaultBondDenom)) {
		return sdkerrors.Wrapf(types.ErrInsufficientBalanceInModuleAccount,
			senderModule)
	}

	// Transfer funds
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, senderModule,
		recipientModule, amt)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrFromBankModule, err))
		return sdkerrors.Wrapf(types.ErrFromBankModule, err.Error())
	}

	k.Logger(ctx).Info(fmt.Sprintf(types.LogInfoFundsTransferred, amount,
		senderModule, recipientModule))
	return nil
}

// AcceptRewards gets rewards from module account and transfer to SRPool module account
func (k Keeper) AcceptRewards(ctx sdk.Context, senderModuleAccName string,
	rewardAmount sdk.Int) error {

	err := k.transferFundsFromModuleToModule(ctx, senderModuleAccName,
		types.SRPoolName, rewardAmount)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf(types.LogErrTransferOfFundsFailed, rewardAmount,
			senderModuleAccName, types.SRPoolName, err))
		return err
	}

	reserver := k.GetReserver(ctx)
	k.updateSrPool(ctx, reserver.SrPool.LockedAmount,
		reserver.SrPool.UnlockedAmount.Add(rewardAmount))

	return nil
}

// getDayID returns the dayID for the end timestamp passed
func getDayID(timestamp uint64) sdk.Int {
	currentEpochTime := sdk.NewInt(time.Now().Unix())

	// Difference in current time and the passed timestamp
	diff := sdk.NewIntFromUint64(timestamp).Sub(currentEpochTime)

	// If difference is zero or negative, return -1
	if diff.IsNegative() || diff.IsZero() {
		return sdk.OneInt().Neg()
	}

	// block_generation_time = 6 seconds (approx.)
	// one_day = 24 * 60 * 60 seconds = 86,400 seconds
	// blocks_per_day = one_day / block_generation_time
	//     = 86400s / 6s = 14,400 blocks per day (approx.)
	blocksPerDay := sdk.NewInt(14400)

	// dayID = diff / blocksPerDay
	return diff.Quo(blocksPerDay)
}
