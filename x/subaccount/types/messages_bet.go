package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/utils"
)

const (
	// typeMsgWager is type of message MsgWager
	typeMsgWager = "subacc_wager"
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

// ValidateBasic validates basic constraints of original msgWager of bet module.
func (msg *MsgWager) ValidateBasic() error {
	if msg.Msg == nil {
		return errors.ErrInvalidRequest.Wrap("msg is nil")
	}
	return msg.Msg.ValidateBasic()
}

// EmitEvent emits the event for the message success.
func (msg *MsgWager) EmitEvent(ctx *sdk.Context, accOwnerAddr string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgWager, msg.Msg.Creator,
		sdk.NewAttribute(attributeKeyBetCreator, msg.Msg.Creator),
		sdk.NewAttribute(attributeKeyBetCreatorOwner, accOwnerAddr),
		sdk.NewAttribute(attributeKeyBetUID, msg.Msg.Props.UID),
	)
	emitter.Emit()
}
