package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrortypes "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/utils"
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
		return sdkerrors.Wrapf(sdkerrortypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !utils.IsValidUID(msg.CampaignUid) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid campaign uid")
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid ticket")
	}

	return nil
}
