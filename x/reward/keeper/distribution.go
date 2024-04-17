package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/x/reward/types"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

// DistributeRewards distributes the rewards according to the input distribution list.
func (k Keeper) DistributeRewards(ctx sdk.Context, receiver types.Receiver) (uint64, error) {
	unlockTS := uint64(0)
	if receiver.RewardAmount.SubaccountAmount.GT(sdk.ZeroInt()) {
		moduleAccAddr := types.RewardPoolFunder{}.GetModuleAcc()
		unlockTS = cast.ToUint64(ctx.BlockTime().Unix()) + receiver.RewardAmount.UnlockPeriod
		if _, err := k.subaccountKeeper.TopUp(ctx, k.accountKeeper.GetModuleAddress(moduleAccAddr).String(), receiver.MainAccountAddr,
			[]subaccounttypes.LockedBalance{
				{
					UnlockTS: unlockTS,
					Amount:   receiver.RewardAmount.SubaccountAmount,
				},
			}); err != nil {
			return unlockTS, sdkerrors.Wrapf(types.ErrSubaccountRewardTopUp, "subaccount address %s, %s", receiver.SubaccountAddr, err)
		}
	}
	if receiver.RewardAmount.MainAccountAmount.GT(sdk.ZeroInt()) {
		if err := k.modFunder.Refund(
			types.RewardPoolFunder{}, ctx,
			sdk.MustAccAddressFromBech32(receiver.MainAccountAddr),
			receiver.RewardAmount.MainAccountAmount,
		); err != nil {
			return unlockTS, err
		}
	}

	return unlockTS, nil
}
