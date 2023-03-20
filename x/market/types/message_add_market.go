package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const typeMsgAddEvent = "add_market"

var _ sdk.Msg = &MsgAddMarket{}

// NewMsgAddEvent creates the new input for adding an market to blockchain
func NewMsgAddEvent(creator string, ticket string) *MsgAddMarket {
	return &MsgAddMarket{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (msg *MsgAddMarket) Route() string {
	return RouterKey
}

// Type returns the msg add market type
func (msg *MsgAddMarket) Type() string {
	return typeMsgAddEvent
}

// GetSigners return the creators address
func (msg *MsgAddMarket) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgAddMarket) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input creation market
func (msg *MsgAddMarket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ticket param")
	}
	return nil
}
