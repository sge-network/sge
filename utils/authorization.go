package utils

import (
	"context"
	"time"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/x/authz"
)

// AuthzKeeper defines the expected authz keeper.
type AuthzKeeper interface {
	GetAuthorization(
		ctx context.Context,
		grantee sdk.AccAddress,
		granter sdk.AccAddress,
		msgType string,
	) (authz.Authorization, *time.Time)
	SaveGrant(
		ctx context.Context,
		grantee, granter sdk.AccAddress,
		authorization authz.Authorization,
		expiration *time.Time,
	) error
	DeleteGrant(
		ctx context.Context,
		grantee, granter sdk.AccAddress,
		msgType string,
	) error
}

func ValidateMsgAuthorization(
	ctx context.Context,
	authzKeeper AuthzKeeper,
	creator, depositor string,
	msg sdk.Msg,
	errAuthorizationNotFound, errAuthorizationNotAccepted error,
) error {
	granteeAddr := sdk.MustAccAddressFromBech32(creator)
	granterAddr, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid granter address (%s)", err)
	}
	authorization, expiration := authzKeeper.GetAuthorization(
		ctx,
		granteeAddr,
		granterAddr,
		sdk.MsgTypeURL(msg),
	)
	if authorization == nil {
		return sdkerrors.Wrapf(
			errAuthorizationNotFound,
			"grantee: %s, granter: %s",
			creator,
			depositor,
		)
	}
	authRes, err := authorization.Accept(ctx, msg)
	if err != nil {
		return sdkerrors.Wrapf(errAuthorizationNotAccepted, "%s", err)
	}

	if authRes.Delete {
		err = authzKeeper.DeleteGrant(ctx, granteeAddr, granterAddr, sdk.MsgTypeURL(msg))
	} else if authRes.Updated != nil {
		err = authzKeeper.SaveGrant(ctx, granteeAddr, granterAddr, authRes.Updated, expiration)
	}
	if err != nil {
		return err
	}

	if !authRes.Accept {
		return errAuthorizationNotAccepted
	}

	return nil
}
