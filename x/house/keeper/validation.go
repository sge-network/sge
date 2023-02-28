package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/house/types"
)

// validateDeposit validates deposit acceptability
func (k msgServer) validateDeposit(ctx sdk.Context, msg *types.MsgDeposit) error {
	if msg.Amount.LT(k.Keeper.GetMinDepositAllowedAmount(ctx)) {
		return sdkerrors.Wrapf(
			types.ErrDepositTooSmall, ": got %s, expected %d", msg.Amount.String(), k.Keeper.GetMinDepositAllowedAmount(ctx),
		)
	}

	return nil
}
