package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sge-network/sge/x/market/types"
)

func (k Keeper) TopUpPriceLockPool(
	ctx sdk.Context,
	funder string, amount sdkmath.Int,
) error {
	if err := k.modFunder.Fund(types.PriceLockFunder{}, ctx, sdk.MustAccAddressFromBech32(funder), amount); err != nil {
		return sdkerrors.Wrapf(types.ErrInsufficientBalanceInPriceLockFunder, "%s", err)
	}

	return nil
}
