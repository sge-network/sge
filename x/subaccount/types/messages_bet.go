package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgWager{}
)

func (m *MsgWager) ValidateBasic() error {
	if m.Msg == nil {
		return errors.ErrInvalidRequest.Wrap("msg is nil")
	}
	return m.Msg.ValidateBasic()
}

func (m *MsgWager) GetSigners() []sdk.AccAddress {
	return m.Msg.GetSigners()
}
