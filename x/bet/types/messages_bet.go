package types

import (
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// TypeMsgPlaceBet is type of message MsgPlaceBet
	TypeMsgPlaceBet = "place_bet"
	// TypeMsgSettleBet is type of message MsgSettleBet
	TypeMsgSettleBet = "settle_bet"

	// SettlementUIDsThreshold is the threshold for the number of UIDs in bulk settlement tx
	SettlementUIDsThreshold = 10
	// BetPlacementThreshold is the threshold for the number bets in bulk placement tx
	BetPlacementThreshold = 10
)

var _ sdk.Msg = &MsgPlaceBet{}

// NewMsgPlaceBet returns a MsgPlaceBet using given data
func NewMsgPlaceBet(
	creator string,
	bet BetPlaceFields,

) *MsgPlaceBet {
	return &MsgPlaceBet{
		Creator: creator,
		Bet:     &bet,
	}
}

// Route returns the module's message router key.
func (msg *MsgPlaceBet) Route() string {
	return RouterKey
}

// Type returns type of its message
func (msg *MsgPlaceBet) Type() string {
	return TypeMsgPlaceBet
}

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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, ErrTextInvalidCreator, err)
	}

	if err = BetFieldsValidation(msg.Bet); err != nil {
		return err
	}

	return nil
}

var _ sdk.Msg = &MsgSettleBet{}

// NewMsgSettleBet returns a MsgSettleBet using given data
func NewMsgSettleBet(creator string, betUID string) *MsgSettleBet {
	return &MsgSettleBet{
		Creator: creator,
		BetUID:  betUID,
	}
}

// Route returns the module's message router key.
func (msg *MsgSettleBet) Route() string {
	return RouterKey
}

// Type returns type of its message
func (msg *MsgSettleBet) Type() string {
	return TypeMsgSettleBet
}

// GetSigners returns the signers of its message
func (msg *MsgSettleBet) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgSettleBet) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does some validate checks on its message
func (msg *MsgSettleBet) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil || msg.Creator == "" || strings.Contains(msg.Creator, " ") {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, ErrTextInvalidCreator, err)
	}

	if !IsValidUID(msg.BetUID) {
		return ErrInvalidBetUID
	}
	return nil
}

// isValidUUID validates the uid
func isValidUUID(uid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uid)
	// _, err := uuid.Parse(uid)
	// return err == nil
}

// IsValidUID validates the uid
func IsValidUID(uid string) bool {
	if len(uid) == 0 || uid == "" || strings.Contains(uid, " ") ||
		!isValidUUID(uid) {
		return false
	}

	return true
}

// NewBet creates and returns a new bet from given message
func NewBet(creator string, bet *BetPlaceFields, odds *BetOdds) (*Bet, error) {
	ticketOddsValue, err := sdk.NewDecFromStr(odds.Value)
	if err != nil {
		return nil, sdkerrors.Wrapf(ErrInConvertingOddsToDec, "%s", err)
	}
	return &Bet{
		Creator:       creator,
		UID:           bet.UID,
		SportEventUID: odds.SportEventUID,
		OddsUID:       odds.UID,
		OddsValue:     ticketOddsValue,
		Amount:        bet.Amount,
		Ticket:        bet.Ticket,
		OddsType:      Bet_OddsType(bet.OddsType),
	}, nil
}

func ExtractSelectedOddsFromTicket(ticketData *BetPlacementTicketPayload) *BetOdds {
	for _, o := range ticketData.Odds {
		if o.UID == ticketData.OddsUID {
			return o
		}
	}
	return nil
}
