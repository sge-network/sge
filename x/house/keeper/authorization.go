package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

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
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid depositor address (%s)", err)
	}
	authorization, expiration := k.authzKeeper.GetAuthorization(
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
	authRes, err := authorization.Accept(ctx, msg)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrAuthorizationNotAccepted, "%s", err)
	}

	if authRes.Delete {
		err = k.authzKeeper.DeleteGrant(ctx, granteeAddr, granterAddr, sdk.MsgTypeURL(msg))
	} else if authRes.Updated != nil {
		err = k.authzKeeper.SaveGrant(ctx, granteeAddr, granterAddr, authRes.Updated, expiration)
	}
	if err != nil {
		return err
	}

	if !authRes.Accept {
		return types.ErrAuthorizationNotAccepted
	}

	return nil
}
