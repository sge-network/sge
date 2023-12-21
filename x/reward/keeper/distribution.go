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
	if receiver.SubAccountAmount.GT(sdk.ZeroInt()) {
		moduleAccAddr := types.RewardPoolFunder{}.GetModuleAcc()
		unlockTS = cast.ToUint64(ctx.BlockTime().Unix()) + receiver.UnlockPeriod
		if _, err := k.subaccountKeeper.TopUp(ctx, k.accountKeeper.GetModuleAddress(moduleAccAddr).String(), receiver.MainAccountAddr,
			[]subaccounttypes.LockedBalance{
				{
					UnlockTS: unlockTS,
					Amount:   receiver.SubAccountAmount,
				},
			}); err != nil {
			return unlockTS, sdkerrors.Wrapf(types.ErrSubAccRewardTopUp, "subaccount address %s, %s", receiver.SubAccountAddr, err)
		}
	}
	if receiver.MainAccountAmount.GT(sdk.ZeroInt()) {
		if err := k.modFunder.Refund(
			types.RewardPoolFunder{}, ctx,
			sdk.MustAccAddressFromBech32(receiver.MainAccountAddr),
			receiver.MainAccountAmount,
		); err != nil {
			return unlockTS, err
		}
	}

	return unlockTS, nil
}
