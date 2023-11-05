package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/reward/types"
	subaccounttypes "github.com/sge-network/sge/x/subaccount/types"
)

// DistributeRewards distributes the rewards according to the input distribution list.
func (k Keeper) DistributeRewards(ctx sdk.Context, funderAddr string, allocation types.Allocation) error {
	if allocation.MainAcc.Amount.GT(sdkmath.ZeroInt()) {
		if err := k.modFunder.Refund(
			types.RewardPoolFunder{}, ctx,
			sdk.MustAccAddressFromBech32(allocation.MainAcc.Addr),
			allocation.MainAcc.Amount,
		); err != nil {
			return err
		}
	}

	if allocation.SubAcc.Amount.GT(sdkmath.ZeroInt()) {
		if _, err := k.subaccountKeeper.TopUp(ctx, funderAddr, allocation.SubAcc.Addr,
			[]subaccounttypes.LockedBalance{
				{
					UnlockTS: allocation.SubAcc.UnlockTS,
					Amount:   allocation.SubAcc.Amount,
				},
			}); err != nil {
			return sdkerrors.Wrapf(types.ErrSubAccRewardTopUp, "subaccount address %s, %s", allocation.SubAcc.Addr, err)
		}
	}

	return nil
}
