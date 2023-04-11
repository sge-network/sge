package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cast"
)

// Ticket is the Interface of ticket.
type Ticket interface {
	// Unmarshal unmarshals the information of the ticket to the v. v must be a pointer.
	Unmarshal(v interface{}) error

	// Verify verifies the ticket signature with the leader public key. if the ticket is verified,
	// then return nil else return invalid signature error
	Verify(pubKey string) error

	// VerifyAny verifies the ticket signature with any of the public keys. if the ticket is verified,
	// then return nil else return invalid signature error
	VerifyAny(pubKeys []string) error

	// ValidateExpiry verifies that the ticket is not expired yet.
	ValidateExpiry(ctx sdk.Context) error
}

func (payload *PubkeysChangeProposalPayload) Validate(leaderIndex uint32) error {
	keyVault := KeyVault{PublicKeys: payload.PublicKeys}
	if err := keyVault.validatePubKeys(); err != nil {
		return err
	}

	count := len(payload.PublicKeys)
	if leaderIndex >= cast.ToUint32(count) {
		return fmt.Errorf("leader index is out of range %d", leaderIndex)
	}

	return nil
}

func (payload *ProposalVotePayload) Validate() error {
	if !(payload.Vote == ProposalVote_PROPOSAL_VOTE_YES ||
		payload.Vote == ProposalVote_PROPOSAL_VOTE_NO) {
		return fmt.Errorf("vote should be no or yes %d", payload.Vote)
	}

	return nil
}
