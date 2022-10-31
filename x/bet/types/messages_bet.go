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
	// TypeMsgPlaceBetSlip is type of message MsgPlaceBetSlip
	TypeMsgPlaceBetSlip = "place_bet_slip"
	// TypeMsgSettleBet is type of message MsgSettleBet
	TypeMsgSettleBet = "settle_bet"
	// TypeMsgSettleBetBulk is type of message MsgSettleBetBulk
	TypeMsgSettleBetBulk = "settle_bet_bulk"

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

var _ sdk.Msg = &MsgPlaceBetSlip{}

// NewMsgPlaceBetSlip returns a MsgPlaceBetSlip using given data
func NewMsgPlaceBetSlip(creator string, bets []*BetPlaceFields) *MsgPlaceBetSlip {
	return &MsgPlaceBetSlip{
		Creator: creator,
		Bets:    bets,
	}
}

// Route returns the module's message router key.
func (msg *MsgPlaceBetSlip) Route() string {
	return RouterKey
}

// Type returns type of its message
func (msg *MsgPlaceBetSlip) Type() string {
	return TypeMsgPlaceBetSlip
}

// GetSigners returns the signers of its message
func (msg *MsgPlaceBetSlip) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgPlaceBetSlip) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does some validate checks on its message
func (msg *MsgPlaceBetSlip) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, ErrTextInvalidCreator, err)
	}
	// If the length of msgs.Bets is zero
	if len(msg.Bets) == 0 {
		return ErrEmptyBetListRequest
	}

	// If the length of msgs.Bets is more than the threshold
	if len(msg.Bets) > BetPlacementThreshold {
		return ErrTooManyBets
	}

	allNil := true
	for _, bet := range msg.Bets {
		if bet != nil {
			allNil = false
			break
		}
	}
	if allNil {
		return ErrEmptyBetListRequest
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

var _ sdk.Msg = &MsgSettleBetBulk{}

// NewMsgSettleBetBulk returns a MsgSettleBetBulk using given data
func NewMsgSettleBetBulk(creator string, betUIDs []string) *MsgSettleBetBulk {
	return &MsgSettleBetBulk{
		Creator: creator,
		BetUIDs: betUIDs,
	}
}

// Route returns the module's message router key.
func (msg *MsgSettleBetBulk) Route() string {
	return RouterKey
}

// Type returns type of its message
func (msg *MsgSettleBetBulk) Type() string {
	return TypeMsgSettleBetBulk
}

// GetSigners returns the signers of its message
func (msg *MsgSettleBetBulk) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgSettleBetBulk) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic does some validate checks on its message
func (msg *MsgSettleBetBulk) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, ErrTextInvalidCreator, err)
	}
	if len(msg.BetUIDs) == 0 {
		return ErrEmptyUidsList
	}
	if len(msg.BetUIDs) > SettlementUIDsThreshold {
		return ErrTooManyUids
	}

	allNil := true
	for _, betUID := range msg.BetUIDs {
		if betUID != "" {
			allNil = false
			break
		}
	}
	if allNil {
		return ErrEmptyUidsList
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
