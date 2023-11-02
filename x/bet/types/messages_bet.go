package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/mrz1836/go-sanitize"
	"github.com/sge-network/sge/utils"
)

const (
	// typeMsgWager is type of message MsgWager
	typeMsgWager = "bet_wager"
)

var _ sdk.Msg = &MsgWager{}

// NewMsgWager returns a MsgWager using given data
func NewMsgWager(
	creator string,
	props WagerProps,
) *MsgWager {
	return &MsgWager{
		Creator: creator,
		Props:   &props,
	}
}

// Route returns the module's message router key.
func (*MsgWager) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgWager) Type() string { return typeMsgWager }

// GetSigners returns the signers of its message
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

// ValidateBasic does some validate checks on its message
func (msg *MsgWager) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil || msg.Creator == "" || strings.Contains(msg.Creator, " ") {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "%s", err)
	}

	return WagerValidation(msg.Props)
}

// EmitEvent emits the event for the message success.
func (msg *MsgWager) EmitEvent(ctx *sdk.Context) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgWager, msg.Creator,
		sdk.NewAttribute(attributeKeyBetCreator, msg.Creator),
		sdk.NewAttribute(attributeKeyBetUID, msg.Props.UID),
	)
	emitter.Emit()
}

// NewBet creates and returns a new bet from given message
func NewBet(creator string, props *WagerProps, odds *BetOdds, meta MetaData) *Bet {
	meta.SelectedOddsValue = sanitize.XSS(meta.SelectedOddsValue)

	return &Bet{
		Creator:           creator,
		UID:               props.UID,
		MarketUID:         odds.MarketUID,
		OddsUID:           odds.UID,
		OddsValue:         odds.Value,
		Amount:            props.Amount,
		MaxLossMultiplier: odds.MaxLossMultiplier,
		Meta:              meta,
	}
}
