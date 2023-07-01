package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/spf13/cast"
)

const typeMsgWithdraw = "house_withdraw"

var _ sdk.Msg = &MsgWithdraw{}

// NewMsgWithdraw creates the new input for withdrawal of a deposit
func NewMsgWithdraw(creator string, marketUID string, amount sdk.Int,
	participationIndex uint64, mode WithdrawalMode, ticket string,
) *MsgWithdraw {
	return &MsgWithdraw{
		Creator:            creator,
		MarketUID:          marketUID,
		ParticipationIndex: participationIndex,
		Mode:               mode,
		Amount:             amount,
		Ticket:             ticket,
	}
}

// Route return the message route for slashing
func (*MsgWithdraw) Route() string { return RouterKey }

// Type returns the msg add market type
func (*MsgWithdraw) Type() string { return typeMsgWithdraw }

// GetSigners return the creators address
func (msg *MsgWithdraw) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes return the marshalled bytes of the msg
func (msg *MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validates the input creation market
func (msg *MsgWithdraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Mode != WithdrawalMode_WITHDRAWAL_MODE_FULL &&
		msg.Mode != WithdrawalMode_WITHDRAWAL_MODE_PARTIAL {
		return ErrInvalidWithdrawMode
	}

	if !utils.IsValidUID(msg.MarketUID) {
		return ErrInvalidMarketUID
	}

	if msg.ParticipationIndex < 1 {
		return ErrInvalidIndex
	}

	if msg.Mode == WithdrawalMode_WITHDRAWAL_MODE_PARTIAL {
		if !msg.Amount.IsPositive() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid withdrawal amount")
		}
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgWithdraw) EmitEvent(ctx *sdk.Context, depositor string, withdrawalID uint64) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgWithdraw, msg.Creator,
		sdk.NewAttribute(attributeKeyCreator, msg.Creator),
		sdk.NewAttribute(attributeKeyDepositor, depositor),
		sdk.NewAttribute(attributeKeyWithdrawalID, cast.ToString(withdrawalID)),
		sdk.NewAttribute(attributeKeyWithdrawMarketUIDParticipantIndex,
			strings.Join([]string{msg.MarketUID, cast.ToString(msg.ParticipationIndex)}, "#"),
		),
	)
	emitter.Emit()
}
