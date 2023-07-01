package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/spf13/cast"
)

// typeMsgPubkeysChangeProposal is type of message MsgPubkeysChangeProposalRequest
const typeMsgPubkeysChangeProposal = "ovm_proposal_pubkeys_change"

var _ sdk.Msg = &MsgSubmitPubkeysChangeProposalRequest{}

// NewMsgPubkeysChangeProposalRequest returns a MsgSubmitPubkeysChangeProposalRequest using given data
func NewMsgPubkeysChangeProposalRequest(
	creator string,
	txs string,
) *MsgSubmitPubkeysChangeProposalRequest {
	return &MsgSubmitPubkeysChangeProposalRequest{
		Creator: creator,
		Ticket:  txs,
	}
}

// Route returns the module's message router key.
func (*MsgSubmitPubkeysChangeProposalRequest) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgSubmitPubkeysChangeProposalRequest) Type() string { return typeMsgPubkeysChangeProposal }

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

// EmitEvent emits the event for the message success.
func (msg *MsgSubmitPubkeysChangeProposalRequest) EmitEvent(ctx *sdk.Context, proposalID uint64, content string) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgPubkeysChangeProposal, msg.Creator,
		sdk.NewAttribute(attributeKeyPubkeysChangeProposalID, cast.ToString(proposalID)),
		sdk.NewAttribute(attributeKeyContent, content),
	)
	emitter.Emit()
}
