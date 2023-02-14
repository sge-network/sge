package types

import (
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sge-network/sge/app/params"
)

// house message types
const (
	// TypeMsgDeposit is type of message MsgDeposit
	TypeMsgDeposit = "deposit"

	// TypeMsgWithdraw is type of message MsgWithdraw
	TypeMsgWithdraw = "withdraw"
)

// NewMsgDeposit creates a new MsgDeposit instance.
//nolint:interfacer
func NewMsgDeposit(depAddr sdk.AccAddress, sportEventUid string, amount sdk.Coin) *MsgDeposit {
	return &MsgDeposit{
		DepositorAddress: depAddr.String(),
		SportEventUid:    sportEventUid,
		Amount:           amount,
	}
}

// Route returns the module's message router key.
func (msg *MsgDeposit) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface.
func (msg *MsgDeposit) Type() string {
	return TypeMsgDeposit
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	depAddr, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{depAddr}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDeposit) ValidateBasic() error {
	if msg.DepositorAddress == "" {
		return ErrEmptyDepositorAddr
	}

	if !isValidUID(msg.SportEventUid) {
		return ErrInvalidSportEventUid
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid deposit amount",
		)
	}

	if msg.Amount.Denom != params.BaseCoinUnit {
		return sdkerrors.Wrapf(
			ErrInvalidDenom, ": got %s, expected %s", msg.Amount.Denom, params.BaseCoinUnit,
		)
	}

	return nil
}

// IsValidUID validates the uid
func isValidUID(uid string) bool {
	if len(uid) == 0 || uid == "" || strings.Contains(uid, " ") ||
		!isValidUUID(uid) {
		return false
	}

	return true
}

// isValidUUID validates the uid
func isValidUUID(uid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uid)
}

// NewMsgWithdraw creates a new MsgWithdraw instance.
//nolint:interfacer
func NewMsgWithdraw(depAddr sdk.AccAddress, sportEventUid string, amount sdk.Coin, pId uint64, mode WithdrawalMode) *MsgWithdraw {
	return &MsgWithdraw{
		DepositorAddress: depAddr.String(),
		SportEventUid:    sportEventUid,
		ParticipantId:    pId,
		Mode:             mode,
		Amount:           amount,
	}
}

// Route returns the module's message router key.
func (msg *MsgWithdraw) Route() string {
	return RouterKey
}

// Type implements the sdk.Msg interface.
func (msg *MsgWithdraw) Type() string {
	return TypeMsgWithdraw
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgWithdraw) GetSigners() []sdk.AccAddress {
	depAddr, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{depAddr}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgWithdraw) ValidateBasic() error {
	if msg.DepositorAddress == "" {
		return ErrEmptyDepositorAddr
	}

	if !isValidUID(msg.SportEventUid) {
		return ErrInvalidSportEventUid
	}

	if msg.ParticipantId < 1 {
		return ErrInvalidSportEventUid
	}

	if msg.Mode != WithdrawalMode_MODE_FULL && msg.Mode != WithdrawalMode_MODE_PARTIAL {
		return ErrInvalidSportEventUid
	}

	if msg.Mode == WithdrawalMode_MODE_PARTIAL {
		if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
			return sdkerrors.Wrap(
				sdkerrors.ErrInvalidRequest,
				"invalid withdrawal amount",
			)
		}

		if msg.Amount.Denom != params.BaseCoinUnit {
			return sdkerrors.Wrapf(
				ErrInvalidDenom, ": got %s, expected %s", msg.Amount.Denom, params.BaseCoinUnit,
			)
		}
	}

	return nil
}
