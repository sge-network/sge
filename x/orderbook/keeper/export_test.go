package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// KeeperTest is a wrapper object for the keeper, It is being used
// to export unexported methods of the keeper
type KeeperTest = Keeper

func (k KeeperTest) Fund(
	mf iModuleFunder,
	ctx sdk.Context,
	senderAcc sdk.AccAddress,
	amount sdk.Int,
) error {
	return k.fund(mf, ctx, senderAcc, amount)
}

func (k KeeperTest) ReFund(
	mf iModuleFunder,
	ctx sdk.Context,
	receiverAcc sdk.AccAddress,
	amount sdk.Int,
) error {
	return k.refund(mf, ctx, receiverAcc, amount)
}
