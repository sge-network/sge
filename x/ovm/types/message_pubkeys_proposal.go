package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// typeMsgPubkeysChangeProposal is type of message MsgPubkeysChangeProposalRequest
const typeMsgPubkeysChangeProposal = "pubkeys_change_proposal"

var _ sdk.Msg = &MsgSubmitPubkeysChangeProposalRequest{}

// MsgSubmitPubkeysChangeProposalRequest returns a MsgSubmitPubkeysChangeProposalRequest using given data
func NewMsgPubkeysChangeProposalRequest(creator string, txs string) *MsgSubmitPubkeysChangeProposalRequest {
	return &MsgSubmitPubkeysChangeProposalRequest{
		Creator: creator,
		Ticket:  txs,
	}
}

// Route returns the module's message router key.
func (msg *MsgSubmitPubkeysChangeProposalRequest) Route() string {
	return RouterKey
}

// Type returns type of its message
func (msg *MsgSubmitPubkeysChangeProposalRequest) Type() string {
	return typeMsgPubkeysChangeProposal
}

// GetSigners returns the signers of its message
func (msg *MsgSubmitPubkeysChangeProposalRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgSubmitPubkeysChangeProposalRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validations on its message
func (msg *MsgSubmitPubkeysChangeProposalRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
