package types

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/spf13/cast"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/utils"
)

const (
	// typeMsgHouseDeposit is type of message MsgHouseDeposit
	typeMsgHouseDeposit = "subacc_house_deposit"
	// typeHouseWithdraw is type of message MsgHouseWithdraw
	typeHouseWithdraw = "subacc_house_withdraw"
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
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgHouseDeposit) ValidateBasic() error {
	if msg.Msg == nil {
		return errors.ErrInvalidRequest.Wrap("msg is nil")
	}
	return msg.Msg.ValidateBasic()
}

// EmitEvent emits the event for the message success.
func (msg *MsgHouseDeposit) EmitEvent(ctx *sdk.Context, subAccAddr string, participationIndex uint64, feeAmount sdkmath.Int) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgHouseDeposit, msg.Msg.Creator,
		sdk.NewAttribute(attributeKeyDepositCreator, msg.Msg.Creator),
		sdk.NewAttribute(attributeKeySubaccountDepositor, subAccAddr),
		sdk.NewAttribute(attributeKeySubaccountDepositFee, feeAmount.String()),
		sdk.NewAttribute(attributeKeyDepositMarketUIDParticipantIndex,
			strings.Join([]string{msg.Msg.MarketUID, cast.ToString(participationIndex)}, "#"),
		),
	)
	emitter.Emit()
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
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgHouseWithdraw) ValidateBasic() error {
	if msg.Msg == nil {
		return errors.ErrInvalidRequest.Wrap("msg is nil")
	}
	return msg.Msg.ValidateBasic()
}

// EmitEvent emits the event for the message success.
func (msg *MsgHouseWithdraw) EmitEvent(ctx *sdk.Context, subAccAddr string, withdrawalID uint64) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeHouseWithdraw, msg.Msg.Creator,
		sdk.NewAttribute(attributeKeyDepositCreator, msg.Msg.Creator),
		sdk.NewAttribute(attributeKeySubaccountDepositor, subAccAddr),
		sdk.NewAttribute(attributeKeyWithdrawalID, cast.ToString(withdrawalID)),
		sdk.NewAttribute(attributeKeyWithdrawMarketUIDParticipantIndex,
			strings.Join([]string{msg.Msg.MarketUID, cast.ToString(msg.Msg.ParticipationIndex)}, "#"),
		),
	)
	emitter.Emit()
}
