package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/spf13/cast"
)

const typeMsgDeposit = "house_deposit"

var _ sdk.Msg = &MsgDeposit{}

// NewMsgDeposit creates the new input for adding deposit to blockchain
func NewMsgDeposit(creator, marketUID string, amount sdk.Int, ticket string) *MsgDeposit {
	return &MsgDeposit{
		Creator:   creator,
		MarketUID: marketUID,
		Amount:    amount,
		Ticket:    ticket,
	}
}

// Route return the message route for slashing
func (*MsgDeposit) Route() string { return RouterKey }

// Type returns the msg add market type
func (*MsgDeposit) Type() string { return typeMsgDeposit }

// GetSigners return the creators address
func (msg *MsgDeposit) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input creation market
func (msg *MsgDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !utils.IsValidUID(msg.MarketUID) {
		return ErrInvalidMarketUID
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid deposit amount",
		)
	}

	return nil
}

// ValidateSanity validates deposit acceptability
func (msg MsgDeposit) ValidateSanity(_ sdk.Context, p *Params) error {
	if msg.Amount.LT(p.MinDeposit) {
		return sdkerrors.Wrapf(
			ErrDepositTooSmall, ": got %s, expected greater or equal to %d",
			msg.Amount.String(), p.MinDeposit,
		)
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgDeposit) EmitEvent(ctx *sdk.Context, depositor string, participationIndex uint64) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgDeposit, msg.Creator,
		sdk.NewAttribute(attributeKeyCreator, msg.Creator),
		sdk.NewAttribute(attributeKeyDepositor, depositor),
		sdk.NewAttribute(attributeKeyDepositMarketUIDParticipantIndex,
			strings.Join([]string{msg.MarketUID, cast.ToString(participationIndex)}, "#"),
		),
	)
	emitter.Emit()
}
