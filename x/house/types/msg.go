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
		return ErrEmptyDelegatorAddr
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
