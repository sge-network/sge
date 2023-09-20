package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// typeMsgCreate is type of message MsgCreate
	typeMsgCreate = "subaccount_create"
)

var _ sdk.Msg = &MsgCreate{}

// Route returns the module's message router key.
func (*MsgCreate) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgCreate) Type() string { return typeMsgCreate }

func (msg *MsgCreate) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgCreate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs a basic validation of the MsgCreate fields.
func (msg *MsgCreate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrortypes.ErrInvalidAddress, "%s", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.SubAccountOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrortypes.ErrInvalidAddress, "%s", err)
	}

	for _, balanceUnlock := range msg.LockedBalances {
		if err = balanceUnlock.Validate(); err != nil {
			return sdkerrors.Wrapf(ErrInvalidLockedBalance, "%s", err)
		}
	}

	return nil
}
