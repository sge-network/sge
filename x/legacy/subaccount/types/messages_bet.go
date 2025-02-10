package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/utils"
	bettypes "github.com/sge-network/sge/x/legacy/bet/types"
)

const (
	// typeMsgWager is type of message MsgWager
	typeMsgWager = "subacc_wager"
)

var _ sdk.Msg = &MsgWager{}

// NewMsgWager returns a MsgWager using given data
func NewMsgWager(creator, ticket string) *MsgWager {
	return &MsgWager{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route returns the module's message router key.
func (*MsgWager) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgWager) Type() string { return typeMsgWager }

func (msg *MsgWager) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgWager) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates basic constraints of original msgWager of bet module.
func (msg *MsgWager) ValidateBasic() error {
	if msg.Ticket == "" {
		return errors.ErrInvalidRequest.Wrap("ticket is nil")
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgWager) EmitEvent(ctx *sdk.Context, wagerMsg *bettypes.MsgWager, accOwnerAddr string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgWager, wagerMsg.Creator,
		sdk.NewAttribute(attributeKeyBetCreator, wagerMsg.Creator),
		sdk.NewAttribute(attributeKeyBetCreatorOwner, accOwnerAddr),
		sdk.NewAttribute(attributeKeyBetUID, wagerMsg.Props.UID),
	)
	emitter.Emit()
}
