package types

import (
	"github.com/sge-network/sge/utils"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
)

const typeMsgAdd = "market_add"

var _ sdk.Msg = &MsgAdd{}

// NewMsgAdd creates the new input for adding a market to blockchain
func NewMsgAdd(creator, ticket string) *MsgAdd {
	return &MsgAdd{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (*MsgAdd) Route() string { return RouterKey }

// Type returns the msg add market type
func (*MsgAdd) Type() string { return typeMsgAdd }

// GetSigners return the creators address
func (msg *MsgAdd) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgAdd) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input creation market
func (msg *MsgAdd) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid ticket param")
	}
	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgAdd) EmitEvent(ctx *sdk.Context, marketUID, bookUID string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgAdd, msg.Creator,
		sdk.NewAttribute(attributeKeyMarketUID, marketUID),
		sdk.NewAttribute(attributeKeyMarketOrderBookUID, bookUID),
	)
	emitter.Emit()
}

// typeMsgUpdate is the market name of update market
const typeMsgUpdate = "market_update"

var _ sdk.Msg = &MsgUpdate{}

// NewMsgUpdate accepts the params to create new update body
func NewMsgUpdate(creator, ticket string) *MsgUpdate {
	return &MsgUpdate{
		Creator: creator,
		Ticket:  ticket,
	}
}

// Route return the message route for slashing
func (*MsgUpdate) Route() string { return RouterKey }

// Type return the update market type
func (*MsgUpdate) Type() string { return typeMsgUpdate }

// GetSigners return the creators address
func (msg *MsgUpdate) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgUpdate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input update market
func (msg *MsgUpdate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Ticket == "" {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidRequest, "invalid ticket param")
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgUpdate) EmitEvent(ctx *sdk.Context, marketUID string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgUpdate, msg.Creator,
		sdk.NewAttribute(attributeKeyMarketUID, marketUID),
	)
	emitter.Emit()
}
