package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgApplyReward = "apply_reward"
)

var _ sdk.Msg = &MsgApplyReward{}

func NewMsgApplyReward(
	creator string,
	campaignUID string,
	ticket string,
) *MsgApplyReward {
	return &MsgApplyReward{
		Creator:     creator,
		CampaignUid: campaignUID,
		Ticket:      ticket,
	}
}

func (msg *MsgApplyReward) Route() string {
	return RouterKey
}

func (msg *MsgApplyReward) Type() string {
	return TypeMsgApplyReward
}

func (msg *MsgApplyReward) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgApplyReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgApplyReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
