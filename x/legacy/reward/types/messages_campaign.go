package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
)

const (
	TypeMsgCreateCampaign = "create_campaign"
	TypeMsgUpdateCampaign = "update_campaign"
	TypeMsgWithdrawFunds  = "withdraw_funds"
)

var _ sdk.Msg = &MsgCreateCampaign{}

func NewMsgCreateCampaign(
	creator string,
	uid string,
	totalFunds sdkmath.Int,
	ticket string,
) *MsgCreateCampaign {
	return &MsgCreateCampaign{
		Creator:    creator,
		Uid:        uid,
		Ticket:     ticket,
		TotalFunds: totalFunds,
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
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateCampaign) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !utils.IsValidUID(msg.Uid) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid uid")
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid ticket")
	}

	if msg.TotalFunds.IsNil() || !msg.TotalFunds.GT(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "total funds should be positive")
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgCreateCampaign) EmitEvent(ctx *sdk.Context, uid string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(TypeMsgCreateCampaign, msg.Creator,
		sdk.NewAttribute(attributeKeyUID, uid),
	)
	emitter.Emit()
}

var _ sdk.Msg = &MsgUpdateCampaign{}

func NewMsgUpdateCampaign(
	creator string,
	uid string,
	topopFunds sdkmath.Int,
	ticket string,
) *MsgUpdateCampaign {
	return &MsgUpdateCampaign{
		Creator:    creator,
		Uid:        uid,
		TopupFunds: topopFunds,
		Ticket:     ticket,
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
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCampaign) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !utils.IsValidUID(msg.Uid) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid uid")
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid ticket")
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgUpdateCampaign) EmitEvent(ctx *sdk.Context, uid string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(TypeMsgUpdateCampaign, msg.Creator,
		sdk.NewAttribute(attributeKeyUID, uid),
	)
	emitter.Emit()
}

var _ sdk.Msg = &MsgWithdrawFunds{}

func NewMsgWithdrawFunds(
	creator string,
	uid string,
	amount sdkmath.Int,
	ticket string,
) *MsgWithdrawFunds {
	return &MsgWithdrawFunds{
		Creator: creator,
		Uid:     uid,
		Amount:  amount,
		Ticket:  ticket,
	}
}

func (msg *MsgWithdrawFunds) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawFunds) Type() string {
	return TypeMsgWithdrawFunds
}

func (msg *MsgWithdrawFunds) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawFunds) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawFunds) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !utils.IsValidUID(msg.Uid) {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid uid")
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid ticket")
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgWithdrawFunds) EmitEvent(ctx *sdk.Context, uid string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(TypeMsgWithdrawFunds, msg.Creator,
		sdk.NewAttribute(attributeKeyUID, uid),
	)
	emitter.Emit()
}
