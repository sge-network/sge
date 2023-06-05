package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/x/house/types"
)

func (k Keeper) ValidateMsgAuthorization(
	ctx sdk.Context,
	creator, depositor string,
	msg sdk.Msg,
) error {
	granteeAddr := sdk.MustAccAddressFromBech32(creator)
	granterAddr, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid depositor address (%s)", err)
	}
	authorization, _ := k.authzKeeper.GetCleanAuthorization(
		ctx,
		granteeAddr,
		granterAddr,
		sdk.MsgTypeURL(msg),
	)
	if authorization == nil {
		return sdkerrors.Wrapf(
			types.ErrAuthorizationNotFound,
			"grantee: %s, granter: %s",
			creator,
			depositor,
		)
	}
	_, err = authorization.Accept(ctx, msg)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrAuthorizationNotAccepted, "%s", err)
	}
	return nil
}
