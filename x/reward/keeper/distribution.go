package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sge-network/sge/x/reward/types"
)

func (k Keeper) DistributeRewards(ctx sdk.Context, distributions []types.Distribution) error {
	for _, d := range distributions {
		if err := d.Allocation.CheckExpiration(uint64(ctx.BlockTime().Unix())); err != nil {
			return err
		}

		switch d.Allocation.ReceiverAccType {
		case types.ReceiverAccType_RECEIVER_ACC_TYPE_MAIN:
			return k.modFunder.Refund(
				types.RewardPoolFunder{}, ctx,
				sdk.MustAccAddressFromBech32(d.AccAddr),
				d.Allocation.Amount,
			)
		case types.ReceiverAccType_RECEIVER_ACC_TYPE_SUB:
			_, found := k.subaccountKeeper.GetSubAccountByOwner(ctx, sdk.MustAccAddressFromBech32(d.AccAddr))
			if !found {
				return sdkerrors.Wrapf(types.ErrSubAccountNotfoundForTheOwner, "owner address %s", d.AccAddr)
			}

		default:
			return types.ErrUnknownAccType
		}
	}
	return nil
}
