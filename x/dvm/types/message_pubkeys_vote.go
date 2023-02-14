package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// typeMsgVotePubkeysChange is type of message MsgPubkeysChangeProposalRequest
const typeMsgVotePubkeysChange = "pubkeys_change_vote"

var _ sdk.Msg = &MsgVotePubkeysChangeRequest{}

// MsgSubmitPubkeysChangeProposalRequest returns a MsgSubmitPubkeysChangeProposalRequest using given data
func NewMsgVotePubkeysChangeRequest(creator, ticket, pubkey string) *MsgVotePubkeysChangeRequest {
	return &MsgVotePubkeysChangeRequest{
		Creator:   creator,
		Ticket:    ticket,
		PublicKey: pubkey,
	}
}

// Route returns the module's message router key.
func (msg *MsgVotePubkeysChangeRequest) Route() string {
	return RouterKey
}

// Type returns type of its message
func (msg *MsgVotePubkeysChangeRequest) Type() string {
	return typeMsgVotePubkeysChange
}

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

	if strings.TrimSpace(msg.PublicKey) == "" {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "public key is a mandatory propery")
	}

	return nil
}
