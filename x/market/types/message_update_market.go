package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// typeMsgUpdateEvent is the market name of update market
const typeMsgUpdateEvent = "update_market"

var _ sdk.Msg = &MsgUpdateMarket{}

// NewMsgUpdateEvent accepts the params to create new update body
func NewMsgUpdateEvent(creator, ticket string) *MsgUpdateMarket {
	return &MsgUpdateMarket{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (msg *MsgUpdateMarket) Route() string {
	return RouterKey
}

// Type return the update market type
func (msg *MsgUpdateMarket) Type() string {
	return typeMsgUpdateEvent
}

// GetSigners return the creators address
func (msg *MsgUpdateMarket) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgUpdateMarket) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input update market
func (msg *MsgUpdateMarket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ticket param")
	}

	return nil
}
