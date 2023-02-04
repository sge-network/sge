package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/x/house/types"
)

// validateDeposit validates deposit acceptability
func (k msgServer) validateDeposit(ctx sdk.Context, msg *types.MsgDeposit) error {
	if msg.Amount.Amount.LT(sdk.NewIntFromUint64(k.Keeper.MinDeposit(ctx))) {
		return sdkerrors.Wrapf(
			types.ErrDepositTooSmall, ": got %s, expected %d", msg.Amount.Amount.String(), k.Keeper.MinDeposit(ctx),
		)
	}

	return nil
}
