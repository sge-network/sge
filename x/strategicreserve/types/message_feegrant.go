package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const typeMsgInvokeFeeGrant = "invoke_feegrant"

var _ sdk.Msg = &MsgInvokeFeeGrant{}

// NewMsgInvokeFeeGrant creates the new input for invoking a feegrant to blockchain.
func NewMsgInvokeFeeGrant(creator string, ticket string) *MsgInvokeFeeGrant {
	return &MsgInvokeFeeGrant{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (msg *MsgInvokeFeeGrant) Route() string {
	return RouterKey
}

// Type returns the msg add market type
func (msg *MsgInvokeFeeGrant) Type() string {
	return typeMsgInvokeFeeGrant
}

// GetSigners return the creators address
func (msg *MsgInvokeFeeGrant) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgInvokeFeeGrant) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input creation market
func (msg *MsgInvokeFeeGrant) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ticket param")
	}
	return nil
}

const typeMsgRevokeFeeGrant = "revoke_feegrant"

var _ sdk.Msg = &MsgRevokeFeeGrant{}

// NewMsgRevokeFeeGrant creates the new input for revoking a feegrant to blockchain.
func NewMsgRevokeFeeGrant(creator string, ticket string) *MsgRevokeFeeGrant {
	return &MsgRevokeFeeGrant{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (msg *MsgRevokeFeeGrant) Route() string {
	return RouterKey
}

// Type returns the msg add market type
func (msg *MsgRevokeFeeGrant) Type() string {
	return typeMsgRevokeFeeGrant
}

// GetSigners return the creators address
func (msg *MsgRevokeFeeGrant) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgRevokeFeeGrant) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input creation market
func (msg *MsgRevokeFeeGrant) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ticket param")
	}
	return nil
}
