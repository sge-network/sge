package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// typeMsgHouseDeposit is type of message MsgHouseDeposit
	typeMsgHouseDeposit = "subaccount_house_deposit"
	// typeHouseWithdraw is type of message MsgHouseWithdraw
	typeHouseWithdraw = "subaccount_house_withdraw"
)

var (
	_ sdk.Msg = &MsgHouseDeposit{}
	_ sdk.Msg = &MsgHouseWithdraw{}
)

// Route returns the module's message router key.
func (*MsgHouseDeposit) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgHouseDeposit) Type() string { return typeMsgHouseDeposit }

func (msg *MsgHouseDeposit) GetSigners() []sdk.AccAddress {
	return msg.Msg.GetSigners()
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgHouseDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (m *MsgHouseDeposit) ValidateBasic() error {
	if m.Msg == nil {
		return errors.ErrInvalidRequest.Wrap("msg is nil")
	}
	return m.Msg.ValidateBasic()
}

// Route returns the module's message router key.
func (*MsgHouseWithdraw) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgHouseWithdraw) Type() string { return typeHouseWithdraw }

func (msg *MsgHouseWithdraw) GetSigners() []sdk.AccAddress {
	return msg.Msg.GetSigners()
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgHouseWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (m *MsgHouseWithdraw) ValidateBasic() error {
	if m.Msg == nil {
		return errors.ErrInvalidRequest.Wrap("msg is nil")
	}
	return m.Msg.ValidateBasic()
}
