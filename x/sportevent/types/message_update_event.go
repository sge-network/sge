package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// typeMsgUpdateEvent is the event name of update event
const typeMsgUpdateEvent = "update_event"

var _ sdk.Msg = &MsgUpdateEvent{}

// NewMsgUpdateEvent accepts the params to create new update body
func NewMsgUpdateEvent(creator, ticket string) *MsgUpdateEvent {
	return &MsgUpdateEvent{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (msg *MsgUpdateEvent) Route() string {
	return RouterKey
}

// Type return the update event type
func (msg *MsgUpdateEvent) Type() string {
	return typeMsgUpdateEvent
}

// GetSigners return the creators address
func (msg *MsgUpdateEvent) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgUpdateEvent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input update event
func (msg *MsgUpdateEvent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ticket param")
	}

	return nil
}
