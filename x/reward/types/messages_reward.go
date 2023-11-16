package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cast"

	"github.com/sge-network/sge/utils"
)

const (
	TypeMsgGrantReward = "apply_reward"
)

var _ sdk.Msg = &MsgGrantReward{}

func NewMsgGrantReward(
	creator string,
	uid string,
	campaignUID string,
	ticket string,
) *MsgGrantReward {
	return &MsgGrantReward{
		Creator:     creator,
		Uid:         uid,
		CampaignUid: campaignUID,
		Ticket:      ticket,
	}
}

func (msg *MsgGrantReward) Route() string {
	return RouterKey
}

func (msg *MsgGrantReward) Type() string {
	return TypeMsgGrantReward
}

func (msg *MsgGrantReward) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgGrantReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGrantReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !utils.IsValidUID(msg.CampaignUid) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid campaign uid")
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid ticket")
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgGrantReward) EmitEvent(ctx *sdk.Context, campaignUID string,
	rewardUID, promoterAddr string, rewardAmount RewardAmount, recevier Receiver,
) {
	mainAmount := sdkmath.ZeroInt()
	if !rewardAmount.MainAccountAmount.IsNil() {
		mainAmount = rewardAmount.MainAccountAmount
	}
	subAmount := sdkmath.ZeroInt()
	if !rewardAmount.SubaccountAmount.IsNil() {
		subAmount = rewardAmount.SubaccountAmount
	}

	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(TypeMsgGrantReward, msg.Creator,
		sdk.NewAttribute(attributeKeyCampaignUID, campaignUID),
		sdk.NewAttribute(attributeKeyRewardUID, rewardUID),
		sdk.NewAttribute(attributeKeyRewardPromoter, promoterAddr),
		sdk.NewAttribute(attributeKeyMainAmount, mainAmount.String()),
		sdk.NewAttribute(attributeKeySubAmount, subAmount.String()),
		sdk.NewAttribute(attributeKeyUnlockPeriod, cast.ToString(rewardAmount.UnlockPeriod)),
		sdk.NewAttribute(attributeKeyReceiver, recevier.String()),
	)
	emitter.Emit()
}
