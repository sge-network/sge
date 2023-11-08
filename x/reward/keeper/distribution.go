package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/reward/types"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

// DistributeRewards distributes the rewards according to the input distribution list.
func (k Keeper) DistributeRewards(ctx sdk.Context, receiver types.Receiver) error {
	if receiver.SubAccountAmount.GT(sdk.ZeroInt()) {
		if _, err := k.subaccountKeeper.TopUp(ctx, types.RewardPoolFunder{}.GetModuleAcc(), receiver.SubAccountAddr,
			[]subaccounttypes.LockedBalance{
				{
					UnlockTS: receiver.UnlockTS,
					Amount:   receiver.SubAccountAmount,
				},
			}); err != nil {
			return sdkerrors.Wrapf(types.ErrSubAccRewardTopUp, "subaccount address %s, %s", receiver.SubAccountAddr, err)
		}
	}
	if receiver.MainAccountAmount.GT(sdk.ZeroInt()) {
		if err := k.modFunder.Refund(
			types.RewardPoolFunder{}, ctx,
			sdk.MustAccAddressFromBech32(receiver.MainAccountAddr),
			receiver.MainAccountAmount,
		); err != nil {
			return err
		}
	}

	return nil
}
