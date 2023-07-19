package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sge-network/sge/utils"
	"github.com/spf13/cast"
)

// typeMsgVotePubkeysChange is type of message MsgPubkeysChangeProposalRequest
const typeMsgVotePubkeysChange = "ovm_vote_pubkeys_change"

var _ sdk.Msg = &MsgVotePubkeysChangeRequest{}

// MsgSubmitPubkeysChangeProposalRequest returns a MsgSubmitPubkeysChangeProposalRequest using given data
func NewMsgVotePubkeysChangeRequest(
	creator, ticket string,
	voterIndex uint32,
) *MsgVotePubkeysChangeRequest {
	return &MsgVotePubkeysChangeRequest{
		Creator:       creator,
		Ticket:        ticket,
		VoterKeyIndex: voterIndex,
	}
}

// Route returns the module's message router key.
func (*MsgVotePubkeysChangeRequest) Route() string { return RouterKey }

// Type returns type of its message
func (*MsgVotePubkeysChangeRequest) Type() string { return typeMsgVotePubkeysChange }

// GetSigners returns the signers of its message
func (msg *MsgVotePubkeysChangeRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns sortJson form of its message
func (msg *MsgVotePubkeysChangeRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validations on its message
func (msg *MsgVotePubkeysChangeRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}

// EmitEvent emits the event for the message success.
func (msg *MsgVotePubkeysChangeRequest) EmitEvent(ctx *sdk.Context, proposalID uint64, publicKey string, vote ProposalVote) {
	emitter := utils.NewEventEmitter(ctx, attributeValueCategory)
	emitter.AddMsg(typeMsgVotePubkeysChange, msg.Creator,
		sdk.NewAttribute(attributeKeyPubkeysChangeProposalID, cast.ToString(proposalID)),
		sdk.NewAttribute(attributeKeyVoterPubKey, publicKey),
		sdk.NewAttribute(attributeKeyVote, vote.String()),
	)
	emitter.Emit()
}
