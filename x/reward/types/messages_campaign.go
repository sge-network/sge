package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrortypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateCampaign = "create_campaign"
	TypeMsgUpdateCampaign = "update_campaign"
	TypeMsgDeleteCampaign = "delete_campaign"
)

var _ sdk.Msg = &MsgCreateCampaign{}

func NewMsgCreateCampaign(
	creator string,
	uid string,
	ticket string,
) *MsgCreateCampaign {
	return &MsgCreateCampaign{
		Creator: creator,
		Uid:     uid,
		Ticket:  ticket,
	}
}

func (msg *MsgCreateCampaign) Route() string {
	return RouterKey
}

func (msg *MsgCreateCampaign) Type() string {
	return TypeMsgCreateCampaign
}

func (msg *MsgCreateCampaign) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateCampaign) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateCampaign) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrortypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateCampaign{}

func NewMsgUpdateCampaign(
	creator string,
	uid string,
	ticket string,
) *MsgUpdateCampaign {
	return &MsgUpdateCampaign{
		Creator: creator,
		Uid:     uid,
		Ticket:  ticket,
	}
}

func (msg *MsgUpdateCampaign) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCampaign) Type() string {
	return TypeMsgUpdateCampaign
}

func (msg *MsgUpdateCampaign) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCampaign) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCampaign) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrortypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteCampaign{}

func NewMsgDeleteCampaign(
	creator string,
	index string,
) *MsgDeleteCampaign {
	return &MsgDeleteCampaign{
		Creator: creator,
		Uid:     index,
	}
}

func (msg *MsgDeleteCampaign) Route() string {
	return RouterKey
}

func (msg *MsgDeleteCampaign) Type() string {
	return TypeMsgDeleteCampaign
}

func (msg *MsgDeleteCampaign) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteCampaign) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteCampaign) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrortypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
