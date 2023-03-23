package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const typeMsgResolveMarket = "resolve_market"

var _ sdk.Msg = &MsgResolveMarket{}

// NewMsgResolveMarket accepts the params to create new resolution body
func NewMsgResolveMarket(creator, ticket string) *MsgResolveMarket {
	return &MsgResolveMarket{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (msg *MsgResolveMarket) Route() string {
	return RouterKey
}

// Type return the resolve market type
func (msg *MsgResolveMarket) Type() string {
	return typeMsgResolveMarket
}

// GetSigners return the creators address
func (msg *MsgResolveMarket) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgResolveMarket) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input resolution market
func (msg *MsgResolveMarket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ticket param")
	}

	return nil
}
