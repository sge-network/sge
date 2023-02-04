package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// KeeperTest is a wrapper object for the keeper, It is being used
// to export unexported methods of the keeper
type KeeperTest = Keeper

func (k KeeperTest) UpdateSrPool(ctx sdk.Context, newLockedAmount sdk.Int,
	newUnlockedAmount sdk.Int,
) {
	k.updateSrPool(ctx, newLockedAmount, newUnlockedAmount)
}

func (k KeeperTest) TransferFundsFromUserToModule(ctx sdk.Context,
	address sdk.AccAddress, moduleAccName string, amount sdk.Int,
) error {
	return k.transferFundsFromUserToModule(ctx, address, moduleAccName,
		amount)
}

func (k KeeperTest) TransferFundsFromModuleToUser(ctx sdk.Context,
	moduleAccName string, address sdk.AccAddress, amount sdk.Int,
) error {
	return k.transferFundsFromModuleToUser(ctx, moduleAccName, address,
		amount)
}

func (k KeeperTest) TransferFundsFromModuleToModule(ctx sdk.Context,
	senderModule string, recipientModule string, amount sdk.Int,
) error {
	return k.transferFundsFromModuleToModule(ctx, senderModule,
		recipientModule, amount)
}

func (k KeeperTest) SetPayoutLock(ctx sdk.Context, uniqueLock string) {
	k.setPayoutLock(ctx, uniqueLock)
}
