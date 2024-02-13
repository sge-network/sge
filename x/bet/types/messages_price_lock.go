package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrtypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
)

const (
	// typeMsgPriceLockPoolTopUp is type of message MsgPriceLockPoolTopUp
	typeMsgPriceLockPoolTopUp = "topup_price_lock"
)

var _ sdk.Msg = &MsgPriceLockPoolTopUp{}

// Route returns the module's message router key.
func (*MsgPriceLockPoolTopUp) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgPriceLockPoolTopUp) Type() string { return typeMsgPriceLockPoolTopUp }

// GetSigners returns the signers of its message
func (msg *MsgPriceLockPoolTopUp) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgPriceLockPoolTopUp) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does some validate checks on its message
func (msg *MsgPriceLockPoolTopUp) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil || msg.Creator == "" || strings.Contains(msg.Creator, " ") {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Funder)
	if err != nil || msg.Funder == "" || strings.Contains(msg.Funder, " ") {
		return sdkerrors.Wrapf(sdkerrtypes.ErrInvalidAddress, "%s", err)
	}

	if msg.Amount.IsNil() || msg.Amount.LTE(sdkmath.ZeroInt()) {
		return sdkerrors.Wrapf(ErrInvalidPriceLockTopUpAmount, "%s", msg.Amount)
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgPriceLockPoolTopUp) EmitEvent(ctx *sdk.Context) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgPriceLockPoolTopUp, msg.Creator,
		sdk.NewAttribute(attributeKeyTopUpPriceLock, msg.Funder),
		sdk.NewAttribute(attributeKeyTopUpPriceLockAmount, msg.Amount.String()),
	)
	emitter.Emit()
}

// NewMsgPriceLockPoolTopUp creates new message object for MsgPriceLockPoolTopUp
func NewMsgPriceLockPoolTopUp(creator, funder string, amount sdkmath.Int) *MsgPriceLockPoolTopUp {
	return &MsgPriceLockPoolTopUp{
		Creator: creator,
		Funder:  funder,
		Amount:  amount,
	}
}
