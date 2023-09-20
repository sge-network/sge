package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// typeMsgWager is type of message MsgWager
	typeMsgWager = "subaccount_wager"
)

var _ sdk.Msg = &MsgWager{}

// Route returns the module's message router key.
func (*MsgWager) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgWager) Type() string { return typeMsgWager }

func (msg *MsgWager) GetSigners() []sdk.AccAddress {
	return msg.Msg.GetSigners()
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgWager) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWager) ValidateBasic() error {
	if msg.Msg == nil {
		return errors.ErrInvalidRequest.Wrap("msg is nil")
	}
	return msg.Msg.ValidateBasic()
}
