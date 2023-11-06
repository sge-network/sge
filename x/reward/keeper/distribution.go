package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/reward/types"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

// DistributeRewards distributes the rewards according to the input distribution list.
func (k Keeper) DistributeRewards(ctx sdk.Context, funderAddr string, isSubAccount bool, receiver types.Receiver) error {
	if isSubAccount {
		if _, err := k.subaccountKeeper.TopUp(ctx, funderAddr, receiver.Addr,
			[]subaccounttypes.LockedBalance{
				{
					UnlockTS: receiver.UnlockTS,
					Amount:   receiver.Amount,
				},
			}); err != nil {
			return sdkerrors.Wrapf(types.ErrSubAccRewardTopUp, "subaccount address %s, %s", receiver.Addr, err)
		}
	} else {
		if receiver.Amount.GT(sdkmath.ZeroInt()) {
			if err := k.modFunder.Refund(
				types.RewardPoolFunder{}, ctx,
				sdk.MustAccAddressFromBech32(receiver.Addr),
				receiver.Amount,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
