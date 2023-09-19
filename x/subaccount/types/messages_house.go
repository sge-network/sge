package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgHouseDeposit{}
	_ sdk.Msg = &MsgHouseWithdraw{}
)

func (m *MsgHouseDeposit) GetSigners() []sdk.AccAddress {
	return m.Msg.GetSigners()
}

func (m *MsgHouseDeposit) ValidateBasic() error {
	if m.Msg == nil {
		return errors.ErrInvalidRequest.Wrap("msg is nil")
	}
	return m.Msg.ValidateBasic()
}

func (m *MsgHouseWithdraw) GetSigners() []sdk.AccAddress {
	return m.Msg.GetSigners()
}

func (m *MsgHouseWithdraw) ValidateBasic() error {
	if m.Msg == nil {
		return errors.ErrInvalidRequest.Wrap("msg is nil")
	}
	return m.Msg.ValidateBasic()
}
