package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// typeMsgMutation is type of message MsgMutation
const typeMsgMutation = "mutation"

var _ sdk.Msg = &MsgMutation{}

// NewMsgMutation returns a MsgMutation using given data
func NewMsgMutation(creator string, txs string) *MsgMutation {
	return &MsgMutation{
		Creator: creator,
		Txs:     txs,
	}
}

// Route returns the module's message router key.
func (msg *MsgMutation) Route() string {
	return RouterKey
}

// Type returns type of its message
func (msg *MsgMutation) Type() string {
	return typeMsgMutation
}

// GetSigners returns the signers of its message
func (msg *MsgMutation) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgMutation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validations on its message
func (msg *MsgMutation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
