package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
)

const typeMsgDeposit = "deposit"

var _ sdk.Msg = &MsgDeposit{}

// NewMsgDeposit creates the new input for adding deposit to blockchain
func NewMsgDeposit(creator, sportEventUID string, amount sdk.Int) *MsgDeposit {
	return &MsgDeposit{
		Creator:       creator,
		SportEventUID: sportEventUID,
		Amount:        amount,
	}
}

// Route return the message route for slashing
func (msg *MsgDeposit) Route() string {
	return RouterKey
}

// Type returns the msg add event type
func (msg *MsgDeposit) Type() string {
	return typeMsgDeposit
}

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

// ValidateBasic validates the input creation event
func (msg *MsgDeposit) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if !utils.IsValidUID(msg.SportEventUID) {
		return ErrInvalidSportEventUID
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			"invalid deposit amount",
		)
	}

	return nil
}
