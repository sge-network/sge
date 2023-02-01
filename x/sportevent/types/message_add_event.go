package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const typeMsgAddEvent = "add_event"

var _ sdk.Msg = &MsgAddSportEventRequest{}

// NewMsgAddEvent creates the new input for adding an event to blockchain
func NewMsgAddEvent(creator string, ticket string) *MsgAddSportEventRequest {
	return &MsgAddSportEventRequest{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (msg *MsgAddSportEventRequest) Route() string {
	return RouterKey
}

// Type returns the msg add event type
func (msg *MsgAddSportEventRequest) Type() string {
	return typeMsgAddEvent
}

// GetSigners return the creators address
func (msg *MsgAddSportEventRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgAddSportEventRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input creation event
func (msg *MsgAddSportEventRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ticket param")
	}
	return nil
}
