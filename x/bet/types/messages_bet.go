package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
)

const (
	// typeMsgPlaceBet is type of message MsgPlaceBet
	typeMsgPlaceBet = "bet_place"
)

var _ sdk.Msg = &MsgPlaceBet{}

// NewMsgPlaceBet returns a MsgPlaceBet using given data
func NewMsgPlaceBet(
	creator string,
	bet PlaceBetFields,
) *MsgPlaceBet {
	return &MsgPlaceBet{
		Creator: creator,
		Bet:     &bet,
	}
}

// Route returns the module's message router key.
func (*MsgPlaceBet) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgPlaceBet) Type() string { return typeMsgPlaceBet }

// GetSigners returns the signers of its message
func (msg *MsgPlaceBet) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgPlaceBet) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does some validate checks on its message
func (msg *MsgPlaceBet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil || msg.Creator == "" || strings.Contains(msg.Creator, " ") {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "%s", err)
	}

	return BetFieldsValidation(msg.Bet)
}

// EmitEvent emits the event for the message success.
func (msg *MsgPlaceBet) EmitEvent(ctx *sdk.Context) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgPlaceBet, msg.Creator,
		sdk.NewAttribute(attributeKeyBetCreator, msg.Creator),
		sdk.NewAttribute(attributeKeyBetUID, msg.Bet.UID),
	)
	emitter.Emit()
}

// NewBet creates and returns a new bet from given message
func NewBet(creator string, bet *PlaceBetFields, oddsType OddsType, odds *BetOdds) *Bet {
	return &Bet{
		Creator:           creator,
		UID:               bet.UID,
		MarketUID:         odds.MarketUID,
		OddsUID:           odds.UID,
		OddsValue:         odds.Value,
		Amount:            bet.Amount,
		OddsType:          oddsType,
		MaxLossMultiplier: odds.MaxLossMultiplier,
	}
}
