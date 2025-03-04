package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// typeMsgUpdateParams is type of message MsgUpdateParams
	typeMsgUpdateParams = "ovm_update_params"
)

var _ sdk.Msg = &MsgUpdateParams{}

// Route returns the module's message router key.
func (*MsgUpdateParams) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgUpdateParams) Type() string { return typeMsgUpdateParams }

// GetSigners returns the expected signers for a MsgUpdateParams message.
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// GetSignBytes implements the LegacyMsg interface.
func (msg MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(Amino.MustMarshalJSON(&msg))
}

// ValidateBasic does a sanity check on the provided data.
func (msg *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.Wrap(err, "invalid authority address")
	}

	if err := msg.Params.Validate(); err != nil {
		return err
	}

	return nil
}
