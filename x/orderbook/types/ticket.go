package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate validates feegrant invoke ticket payload.
func (payload *InvokeFeeGrantPayload) Validate() error {
	_, err := sdk.AccAddressFromBech32(payload.Grantee)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid grantee address (%s)", err)
	}

	return nil
}

// Validate validates feegrant revoke ticket payload.
func (payload *RevokeFeeGrantPayload) Validate() error {
	_, err := sdk.AccAddressFromBech32(payload.Grantee)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid grantee address (%s)", err)
	}

	return nil
}
