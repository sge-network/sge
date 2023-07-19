package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
)

const (
	// typeMsgPlace is type of message MsgPlace
	typeMsgPlace = "bet_place"
)

var _ sdk.Msg = &MsgPlace{}

// NewMsgPlace returns a MsgPlace using given data
func NewMsgPlace(
	creator string,
	props WagerProps,
) *MsgPlace {
	return &MsgPlace{
		Creator: creator,
		Props:   &props,
	}
}

// Route returns the module's message router key.
func (*MsgPlace) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgPlace) Type() string { return typeMsgPlace }

// GetSigners returns the signers of its message
func (msg *MsgPlace) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgPlace) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does some validate checks on its message
func (msg *MsgPlace) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil || msg.Creator == "" || strings.Contains(msg.Creator, " ") {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "%s", err)
	}

	return WagerValidation(msg.Props)
}

// EmitEvent emits the event for the message success.
func (msg *MsgPlace) EmitEvent(ctx *sdk.Context) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgPlace, msg.Creator,
		sdk.NewAttribute(attributeKeyBetCreator, msg.Creator),
		sdk.NewAttribute(attributeKeyBetUID, msg.Props.UID),
	)
	emitter.Emit()
}

// NewBet creates and returns a new bet from given message
func NewBet(creator string, props *WagerProps, oddsType OddsType, odds *BetOdds) *Bet {
	return &Bet{
		Creator:           creator,
		UID:               props.UID,
		MarketUID:         odds.MarketUID,
		OddsUID:           odds.UID,
		OddsValue:         odds.Value,
		Amount:            props.Amount,
		OddsType:          oddsType,
		MaxLossMultiplier: odds.MaxLossMultiplier,
	}
}
