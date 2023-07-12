package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateSubAccount{}
var _ sdk.Msg = &MsgTopUp{}

func (msg *MsgCreateSubAccount) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// ValidateBasic performs a basic validation of the MsgCreateSubAccountRequest fields.
func (msg *MsgCreateSubAccount) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.ErrInvalidAddress
	}

	_, err = sdk.AccAddressFromBech32(msg.SubAccountOwner)

	if err != nil {
		return errors.ErrInvalidAddress
	}

	for _, balanceUnlock := range msg.LockedBalances {
		if err = balanceUnlock.Validate(); err != nil {
			return errors.Wrapf(err, "invalid locked balance")
		}
	}

	return nil
}

func (msg *MsgTopUp) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.ErrInvalidAddress
	}

	_, err = sdk.AccAddressFromBech32(msg.SubAccount)

	if err != nil {
		return errors.ErrInvalidAddress
	}

	for _, balanceUnlock := range msg.LockedBalances {
		if err = balanceUnlock.Validate(); err != nil {
			return errors.Wrapf(err, "invalid locked balance")
		}
	}

	return nil
}

func (msg *MsgTopUp) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}
