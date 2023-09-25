package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/reward/types"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

// DistributeRewards distributes the rewards according to the input distribution list.
func (k Keeper) DistributeRewards(ctx sdk.Context, funderAddr string, distributions []types.Distribution) error {
	for _, d := range distributions {
		switch d.Allocation.ReceiverAccType {
		case types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN:
			return k.modFunder.Refund(
				types.RewardPoolFunder{}, ctx,
				sdk.MustAccAddressFromBech32(d.AccAddr),
				d.Allocation.Amount,
			)
		case types.ReceiverAccType_RECEIVER_ACC_TYPE_SUB:
			if _, err := k.subaccountKeeper.TopUp(ctx, funderAddr, d.AccAddr,
				[]subaccounttypes.LockedBalance{
					{
						UnlockTS: d.Allocation.UnlockTS,
						Amount:   d.Allocation.Amount,
					},
				}); err != nil {
				return sdkerrors.Wrapf(types.ErrSubAccRewardTopUp, "owner address %s", d.AccAddr)
			}
		default:
			return types.ErrUnknownAccType
		}
	}
	return nil
}
