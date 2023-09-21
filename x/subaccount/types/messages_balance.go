package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
)

const (
	// typeMsgTopUp is type of message MsgTopUp
	typeMsgTopUp = "subacc_topup"
	// typeMsgWithdrawUnlockedBalances is type of message MsgWithdrawUnlockedBalances
	typeMsgWithdrawUnlockedBalances = "subacc_withdraw_unlocked"
)

var (
	_ sdk.Msg = &MsgTopUp{}
	_ sdk.Msg = &MsgWithdrawUnlockedBalances{}
)

// Route returns the module's message router key.
func (*MsgTopUp) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgTopUp) Type() string { return typeMsgTopUp }

func (msg *MsgTopUp) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// EmitEvent emits the event for the message success.
func (msg *MsgTopUp) EmitEvent(ctx *sdk.Context, subAccAddr string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgTopUp, msg.Creator,
		sdk.NewAttribute(attributeKeySubAcc, subAccAddr),
	)
	emitter.Emit()
}

func (msg *MsgTopUp) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.ErrInvalidAddress
	}

	_, err = sdk.AccAddressFromBech32(msg.SubAccount)
	if err != nil {
		return errors.ErrInvalidAddress
	}

	for _, balanceUnlock := range msg.LockedBalances {
		if err = balanceUnlock.Validate(); err != nil {
			return sdkerrors.Wrapf(err, "invalid locked balance")
		}
	}

	return nil
}

// Route returns the module's message router key.
func (*MsgWithdrawUnlockedBalances) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgWithdrawUnlockedBalances) Type() string { return typeMsgWithdrawUnlockedBalances }

func (msg *MsgWithdrawUnlockedBalances) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (msg *MsgWithdrawUnlockedBalances) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errors.ErrInvalidAddress
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgWithdrawUnlockedBalances) EmitEvent(ctx *sdk.Context, subAccAddr string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgWithdrawUnlockedBalances, msg.Creator,
		sdk.NewAttribute(attributeKeySubAcc, subAccAddr),
	)
	emitter.Emit()
}
