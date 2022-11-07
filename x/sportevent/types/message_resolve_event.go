package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const typeMsgResolveEvent = "resolve_event"

var _ sdk.Msg = &MsgResolveSportEvent{}

// NewMsgResolveEvent accepts the params to create new resolution body
func NewMsgResolveEvent(creator, ticket string) *MsgResolveSportEvent {
	return &MsgResolveSportEvent{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (msg *MsgResolveSportEvent) Route() string {
	return RouterKey
}

// Type return the resolve event type
func (msg *MsgResolveSportEvent) Type() string {
	return typeMsgResolveEvent
}

// GetSigners return the creators address
func (msg *MsgResolveSportEvent) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgResolveSportEvent) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input resolution event
func (msg *MsgResolveSportEvent) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ticket param")
	}

	return nil
}
